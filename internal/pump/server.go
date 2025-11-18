package pump

import (
	"context"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/neee333ko/IAM/internal/pump/analytics"
	"github.com/neee333ko/IAM/internal/pump/config"
	"github.com/neee333ko/IAM/internal/pump/option"
	"github.com/neee333ko/IAM/internal/pump/pump"
	"github.com/neee333ko/IAM/internal/pump/store"
	redisstore "github.com/neee333ko/IAM/internal/pump/store/redis"
	"github.com/neee333ko/log"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

type Server struct {
	Client           store.AnalyticsStorage
	PurgeDelay       int64
	OmitDetailedPump bool
	PumpOp           map[string]*option.PumpOption
	mutex            *redsync.Mutex
}

var pmps map[string]pump.Pump

const KeyAnalytics = "iam.system.analyics"

func CreateServerFromConfig(cfg *config.Config) *Server {
	client := CreateRedisPool(cfg)
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	return &Server{
		Client:           &redisstore.RedisCluster{},
		PurgeDelay:       cfg.PurgeDelay,
		OmitDetailedPump: cfg.OmitDetailedPump,
		PumpOp:           cfg.PumpOp,
		mutex:            rs.NewMutex("iam-pump", redsync.WithExpiry(10*time.Minute)),
	}
}

type PreparedServer struct {
	*Server
}

func (s *Server) PrepareRun() *PreparedServer {
	s.initialize()

	return &PreparedServer{Server: s}
}

func (ps *PreparedServer) Run(c <-chan int) error {
	t := time.NewTicker(time.Duration(ps.PurgeDelay))
	defer t.Stop()

	for {
		select {
		case <-t.C:
			ps.pump()
		case <-c:
			log.Info("pump stop")
			return nil
		}
	}
}

func (ps *PreparedServer) pump() {
	log.Info("pumping...")

	if err := ps.mutex.Lock(); err != nil {
		log.Errorf("error when acquiring the lock: %s.\n", err.Error())
		return
	}

	defer func() {
		if _, err := ps.mutex.Unlock(); err != nil {
			log.Errorf("error when unlocking: %s.\n", err.Error())
		}
	}()

	data := ps.Client.GetAndDeleteSet(KeyAnalytics)
	if len(data) == 0 {
		return
	}

	records := make([]interface{}, 0, len(data))

	for _, value := range data {
		var record *analytics.AnalyticsRecord = new(analytics.AnalyticsRecord)
		s, _ := value.(string)
		err := msgpack.Unmarshal([]byte(s), record)
		if err != nil {
			continue
		}

		if ps.OmitDetailedPump {
			record.Pool = ""
			record.Deciders = ""
		}

		records = append(records, record)
	}

	var wg *sync.WaitGroup = new(sync.WaitGroup)
	for _, pmp := range pmps {
		wg.Add(1)
		go executePump(pmp, wg, ps.PurgeDelay, records)
	}

	wg.Wait()

	log.Info("pump finished.")
}

func filter(pump pump.Pump, data []interface{}) []interface{} {
	filtered_data := make([]interface{}, 0)

	f := pump.GetFilters()
	var should_filter bool
	if f != nil && f.HasFilter() {
		should_filter = true
	}

	for _, item := range data {
		r, _ := item.(*analytics.AnalyticsRecord)

		if should_filter {
			if f.ShouldFilter(r) {
				continue
			}
		}

		if pump.GetOmitDetailedRecording() {
			r.Pool = ""
			r.Deciders = ""
		}
	}

	return filtered_data
}

func executePump(pump pump.Pump, wg *sync.WaitGroup, purgeDelay int64, data []interface{}) {
	timer := time.AfterFunc(time.Duration(purgeDelay), func() {
		if pump.GetTimeout() == 0 {
			log.Warnf("[%s]: pumping has taken too much time, should set timeout of pump.\n", pump.GetName())
		} else if pump.GetTimeout() > int(purgeDelay) {
			log.Warnf("[%s]: pumping has taken too much time than purgeDelay, should set lower timeout.\n", pump.GetName())
		}
	})
	defer timer.Stop()
	defer wg.Done()

	filtered_data := filter(pump, data)

	var ctx context.Context
	var cancel context.CancelFunc
	if pump.GetTimeout() != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(pump.GetTimeout()))
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	errchan := make(chan error)

	go func() {
		errchan <- pump.WriteData(ctx, filtered_data)
	}()

	select {
	case err := <-errchan:
		log.Errorf("[%s]: error when pumping data:%s.\n", pump.GetName(), err.Error())
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			log.Errorf("[%s]: canceled because of timeout.\n", pump.GetName())
		case context.Canceled:
			log.Errorf("[%s]: canceled because of cancel called.\n", pump.GetName())
		}
	}
}

func (s *Server) initialize() {
	for key, value := range s.PumpOp {
		var Type string
		if value.Type != "" {
			Type = value.Type
		} else {
			Type = key
		}

		pmp, err := pump.GetPumpByName(Type)
		if err != nil {
			log.Errorf("failed to get pump: %s.\n", Type)
			continue
		}

		pmpins := pmp.New()

		err = pmpins.Init(value.Meta)
		if err != nil {
			log.Errorf("failed to initialize pump: %s, err: %s.\n", pmpins.GetName(), err.Error())
			continue
		}

		pmpins.SetFilters(value.Filters)
		pmpins.SetTimeout(int(value.Timeout))
		pmpins.SetOmitDetailedRecording(value.OmitDetailedPump)

		pmps[pmpins.GetName()] = pmpins
	}
}

func CreateRedisPool(cfg *config.Config) redis.UniversalClient {
	config := &redisstore.Config{
		Host:                  cfg.RedisOp.Host,
		Port:                  cfg.RedisOp.Port,
		Addrs:                 cfg.RedisOp.Addrs,
		MasterName:            cfg.RedisOp.MasterName,
		Username:              cfg.RedisOp.Username,
		Password:              cfg.RedisOp.Password,
		Database:              cfg.RedisOp.Database,
		MaxIdle:               cfg.RedisOp.MaxIdle,
		MaxActive:             cfg.RedisOp.MaxIdle,
		Timeout:               cfg.RedisOp.Timeout,
		EnableCluster:         cfg.RedisOp.EnableCluster,
		UseSSL:                cfg.RedisOp.UseSSL,
		SSLInsecureSkipVerify: cfg.RedisOp.SSLInsecureSkipVerify,
	}

	return redisstore.NewRedisClusterPool(config)
}
