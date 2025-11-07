package apiserver

import (
	"context"
	"net"

	"github.com/neee333ko/IAM/internal/apiserver/config"
	storage "github.com/neee333ko/IAM/internal/apiserver/store"
	"github.com/neee333ko/IAM/internal/apiserver/store/mysql"
	"github.com/neee333ko/IAM/internal/pkg/server"
	"github.com/neee333ko/IAM/pkg/shutdown"
	"github.com/neee333ko/IAM/pkg/shutdown/posixsignal"
	redis "github.com/neee333ko/IAM/pkg/storage"
	pb "github.com/neee333ko/api/proto/v1"
	"github.com/neee333ko/log"
	"google.golang.org/grpc"
)

type Server struct {
	genericserver *server.GenericServer
	grpcserver    *GrpcServer
	shutdown      *shutdown.GracefulShutdown
}

func CreateServerFromConfig(config *config.Config) *Server {
	s := new(Server)
	s.shutdown = shutdown.NewGS()
	s.shutdown.AddShutdownManager(posixsignal.NewPSManager())

	genericserver := server.NewConfig().Apply(config).Complete().New()
	completedExtraCfg := NewExtraConfig().Apply(config).Complete()

	s.InitRedisClient(completedExtraCfg)

	s.genericserver = genericserver
	s.grpcserver = completedExtraCfg.New()

	InitRoute(s.genericserver.Engine)

	return s
}

type SServer struct {
	*Server
}

func (s *Server) PreparedRun() *SServer {
	ss := &SServer{Server: s}

	ss.shutdown.AddShutdownCallback(shutdown.CallBackFn(func(s string) error {
		log.Infof("shutdown manager: %s, grpc callback fn called...\n", s)
		ss.grpcserver.Close()

		return nil
	}))

	ss.shutdown.AddShutdownCallback(shutdown.CallBackFn(func(s string) error {
		log.Infof("shutdown manager: %s, http&&https callback fn called...\n", s)
		ss.genericserver.Close()

		return nil
	}))

	return ss
}

func (ss *SServer) Run(basename string) {
	log.Info("iam-apiserver running...")

	go ss.grpcserver.Run()

	ss.shutdown.Start()

	ss.genericserver.Run()
}

type ExtraConfig struct {
	*config.Config
}

func (c *ExtraConfig) Apply(config *config.Config) *ExtraConfig {
	return &ExtraConfig{Config: config}
}

func NewExtraConfig() *ExtraConfig {
	return &ExtraConfig{}
}

type CompletedExConfig struct {
	*ExtraConfig
}

func (c *ExtraConfig) Complete() *CompletedExConfig {
	return &CompletedExConfig{c}
}

func (c *CompletedExConfig) New() *GrpcServer {
	listener, err := net.Listen("tcp", c.GRPCOp.Address)
	if err != nil {
		log.Fatal("error when creating listener for grpc!")
	}

	store := mysql.GetMysqlInsOr(c.MysqlOp)
	storage.SetClient(store)

	server := grpc.NewServer()
	grpcserver := &GrpcServer{Store: store, server: server, listener: listener}
	pb.RegisterCacheServer(server, grpcserver)

	return grpcserver
}

func (s *Server) InitRedisClient(c *CompletedExConfig) {
	ctx, cancel := context.WithCancel(context.Background())
	s.shutdown.AddShutdownCallback(shutdown.CallBackFn(func(s string) error {
		log.Infof("shutdown manager: %s, redis callback fn called...\n", s)
		cancel()
		return nil
	}))

	cfg := &redis.Config{
		Host:                  c.RedisOp.Host,
		Port:                  c.RedisOp.Port,
		Addrs:                 c.RedisOp.Addrs,
		MasterName:            c.RedisOp.MasterName,
		Username:              c.RedisOp.Username,
		Password:              c.RedisOp.Password,
		Database:              c.RedisOp.Database,
		MaxIdle:               c.RedisOp.MaxIdle,
		MaxActive:             c.RedisOp.MaxActive,
		Timeout:               c.RedisOp.Timeout,
		EnableCluster:         c.RedisOp.EnableCluster,
		UseSSL:                c.RedisOp.UseSSL,
		SSLInsecureSkipVerify: c.RedisOp.SSLInsecureSkipVerify,
	}

	go redis.ConnectToRedis(ctx, cfg)
}
