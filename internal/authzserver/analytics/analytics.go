package analytics

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/neee333ko/IAM/pkg/storage"
	"github.com/neee333ko/log"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	KeyAnalytics     = "iam.cluster.analyics"
	ErrAnalyticsDown = errors.New("analytics is down")
)

type AnalyticRecord struct {
	Timestamp  string    `json:"timestamp"`
	Username   string    `json:"username"`
	Request    string    `json:"request"`
	Pool       string    `json:"pool"`
	Deciders   string    `json:"deciders"`
	Conclusion string    `json:"conclusion"`
	Effect     string    `json:"effect"`
	ExpireAt   time.Time `json:"expireAt" bson:"expireAt"`
}

func (record *AnalyticRecord) SetExpire(duration int64) {
	d := time.Duration(duration) * time.Second

	if d == 0 {
		d = 24 * 365 * 100 * time.Hour
	}

	time := time.Now()
	expireAt := time.Add(d)
	record.ExpireAt = expireAt
}

type Analytics struct {
	analyticHandler    storage.AnalyticsHandler
	poolSize           int64
	workerBufferSize   int64
	status             *atomic.Int64
	forcedPushInterval time.Duration
	recordChan         chan AnalyticRecord
	wg                 sync.WaitGroup
}

func NewAnalytics(option *AnalyticOption, prefix string) {
	analytics = Analytics{
		analyticHandler:    &storage.RedisCluster{KeyPrefix: prefix},
		poolSize:           option.PoolSize,
		workerBufferSize:   option.BufferSize / option.PoolSize,
		status:             &atomic.Int64{},
		forcedPushInterval: time.Duration(option.ForcedPushInterval),
		recordChan:         make(chan AnalyticRecord),
		wg:                 sync.WaitGroup{},
	}
}

var analytics Analytics

func GetAnalytics() *Analytics {
	return &analytics
}

func (a *Analytics) Run() {
	log.Info("analytics start running...")

	a.status.Store(0)

	for range a.poolSize {
		a.wg.Add(1)

		go a.Listen()
	}
}

func (a *Analytics) SendRecord(record AnalyticRecord) error {
	status := a.status.Load()

	if status != 0 {
		return ErrAnalyticsDown
	}

	a.recordChan <- record

	return nil
}

func (a *Analytics) Stop() {
	a.status.Store(1)

	close(a.recordChan)

	a.wg.Wait()
}

func (a *Analytics) Listen() {
	defer a.wg.Done()

	buffer := make([][]byte, 0, a.workerBufferSize)
	shouldSend := false

	for {
		shouldSend = false
		select {
		case record, ok := <-a.recordChan:
			if !ok {
				// chan closed, need to push rest buffer
				a.analyticHandler.AppendToSetPipelined(KeyAnalytics, buffer)
				return
			}

			msg, err := msgpack.Marshal(record)
			if err != nil {
				log.Errorf("marshal record failed: %s\n", err.Error())
				continue
			}

			buffer = append(buffer, msg)

			if len(buffer) == int(a.workerBufferSize) {
				shouldSend = true
			}
		case <-time.After(a.forcedPushInterval * time.Millisecond):
			shouldSend = true
		}

		if len(buffer) != 0 && shouldSend {
			a.analyticHandler.AppendToSetPipelined(KeyAnalytics, buffer)
		}
	}
}
