package apiserver

import (
	"context"

	"github.com/neee333ko/IAM/internal/apiserver/config"
	cachev1 "github.com/neee333ko/IAM/internal/apiserver/controller/v1/cache"
	storage "github.com/neee333ko/IAM/internal/apiserver/store"
	"github.com/neee333ko/IAM/internal/apiserver/store/mysql"
	"github.com/neee333ko/IAM/internal/pkg/option"
	"github.com/neee333ko/IAM/internal/pkg/server"
	"github.com/neee333ko/IAM/pkg/shutdown"
	"github.com/neee333ko/IAM/pkg/shutdown/posixsignal"
	redis "github.com/neee333ko/IAM/pkg/storage"
	pb "github.com/neee333ko/api/proto/v1"
	"github.com/neee333ko/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	redisOp       *option.RedisOption
	genericserver *server.GenericServer
	grpcserver    *GRPCServer
	shutdown      *shutdown.GracefulShutdown
}

func CreateServerFromConfig(config *config.Config) (*Server, error) {
	s := new(Server)
	s.shutdown = shutdown.NewGS()
	s.shutdown.AddShutdownManager(posixsignal.NewPSManager())
	s.redisOp = config.RedisOp

	genericCfg, err := buildGenericConfig(config)
	if err != nil {
		return nil, err
	}
	genericserver := genericCfg.Complete().New()

	grpcserver, err := buildExtraConfig(config).Complete().New()
	if err != nil {
		return nil, err
	}

	s.genericserver = genericserver
	s.grpcserver = grpcserver

	return s, nil
}

type PreparedServer struct {
	*Server
}

func (s *Server) PreparedRun() *PreparedServer {
	InitRoute(s.genericserver.Engine)
	s.InitRedisClient(s.redisOp)

	s.shutdown.AddShutdownCallback(shutdown.CallBackFn(func(managerName string) error {
		log.Infof("shutdown manager: %s, grpc callback fn called...\n", managerName)
		s.grpcserver.Close()

		return nil
	}))

	s.shutdown.AddShutdownCallback(shutdown.CallBackFn(func(managerName string) error {
		log.Infof("shutdown manager: %s, http && https callback fn called...\n", managerName)
		s.genericserver.Close()

		return nil
	}))

	s.shutdown.AddShutdownCallback(shutdown.CallBackFn(func(managerName string) error {
		store := mysql.GetMysqlInsOr(nil)
		if store != nil {
			log.Infof("shutdown manager: %s, mysql callback fn called...\n", managerName)
			return store.Close()
		}

		return nil
	}))

	return &PreparedServer{Server: s}
}

func (s *PreparedServer) Run() error {
	log.Info("iam-apiserver running...")

	if err := s.grpcserver.Run(); err != nil {
		return err
	}

	s.shutdown.Start()

	return s.genericserver.Run()
}

func buildGenericConfig(config *config.Config) (cfg *server.Config, err error) {
	err = config.InsecureOp.ApplyTo(cfg)
	if err != nil {
		return nil, err
	}

	err = config.SecureOp.ApplyTo(cfg)
	if err != nil {
		return nil, err
	}

	err = config.Feature.ApplyTo(cfg)
	if err != nil {
		return nil, err
	}

	err = config.RunOp.ApplyTo(cfg)
	if err != nil {
		return nil, err
	}

	return
}

type ExtraConfig struct {
	GRPCOp  *option.GRPCOption
	MysqlOp *option.MysqlOption
	RedisOp *option.RedisOption
	CertKey *option.CertKey
}

func buildExtraConfig(config *config.Config) *ExtraConfig {
	return &ExtraConfig{
		GRPCOp:  config.GRPCOp,
		MysqlOp: config.MysqlOp,
		RedisOp: config.RedisOp,
		CertKey: &config.SecureOp.TLSConfig.CertKey,
	}
}

type CompletedExConfig struct {
	*ExtraConfig
}

func (c *ExtraConfig) Complete() *CompletedExConfig {
	return &CompletedExConfig{c}
}

func (c *CompletedExConfig) New() (*GRPCServer, error) {
	tls, err := credentials.NewServerTLSFromFile(c.CertKey.CertFile, c.CertKey.CertKey)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer(grpc.Creds(tls), grpc.MaxRecvMsgSize(c.GRPCOp.MaxRecvMsgSize))

	store := mysql.GetMysqlInsOr(c.MysqlOp)
	storage.SetClient(store)

	cacheServer := cachev1.GetCacheInsOr(store)
	pb.RegisterCacheServer(server, cacheServer)
	reflection.Register(server)

	return &GRPCServer{Address: c.GRPCOp.Address, Server: server}, nil
}

func (s *Server) InitRedisClient(redisOp *option.RedisOption) {
	ctx, cancel := context.WithCancel(context.Background())
	s.shutdown.AddShutdownCallback(shutdown.CallBackFn(func(managerName string) error {
		log.Infof("shutdown manager: %s, redis callback fn called...\n", managerName)
		cancel()
		return nil
	}))

	cfg := &redis.Config{
		Host:                  redisOp.Host,
		Port:                  redisOp.Port,
		Addrs:                 redisOp.Addrs,
		MasterName:            redisOp.MasterName,
		Username:              redisOp.Username,
		Password:              redisOp.Password,
		Database:              redisOp.Database,
		MaxIdle:               redisOp.MaxIdle,
		MaxActive:             redisOp.MaxActive,
		Timeout:               redisOp.Timeout,
		EnableCluster:         redisOp.EnableCluster,
		UseSSL:                redisOp.UseSSL,
		SSLInsecureSkipVerify: redisOp.SSLInsecureSkipVerify,
	}

	go redis.ConnectToRedis(ctx, cfg)
}
