package authzserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/neee333ko/IAM/internal/authzserver/analytics"
	"github.com/neee333ko/IAM/internal/authzserver/config"
	"github.com/neee333ko/IAM/internal/authzserver/load"
	"github.com/neee333ko/IAM/internal/authzserver/load/cache"
	"github.com/neee333ko/IAM/internal/authzserver/store"
	grpcclient "github.com/neee333ko/IAM/internal/authzserver/store/grpc"
	"github.com/neee333ko/IAM/internal/pkg/option"
	"github.com/neee333ko/IAM/internal/pkg/server"
	"github.com/neee333ko/IAM/pkg/shutdown"
	"github.com/neee333ko/IAM/pkg/shutdown/posixsignal"
	"github.com/neee333ko/IAM/pkg/storage"
	pb "github.com/neee333ko/api/proto/v1"
	"github.com/neee333ko/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const RedisKeyPrefix = "analytics-"

type Server struct {
	shutdown      *shutdown.GracefulShutdown
	genericServer *server.GenericServer
	redisOp       *option.RedisOption
	analytics     *analytics.Analytics
	load          *load.Reloader
}

func CreateServerFromConfig(cfg *config.Config) *Server {
	gfshutdown := shutdown.NewGS()
	gfshutdown.AddShutdownManager(posixsignal.NewPSManager())

	genericServer := NewGenericConfig(cfg).ApplyTo().Complete().New()
	NewExtraConfig(cfg).Complete().New()

	ctx, cancel := context.WithCancel(context.Background())
	gfshutdown.AddShutdownCallback(shutdown.CallBackFn(func(s string) error {
		log.Infof("manager: %s, loader shutdown callback called...\n", s)
		cancel()
		return nil
	}))
	loader := load.NewReloader(ctx, cache.GetCacheInsOr())

	analytics := analytics.GetAnalytics()
	gfshutdown.AddShutdownCallback(shutdown.CallBackFn(func(s string) error {
		log.Infof("manager: %s, analytics shutdown callback called...\n", s)
		analytics.Stop()
		return nil
	}))

	return &Server{
		shutdown:      gfshutdown,
		genericServer: genericServer,
		redisOp:       cfg.RedisOp,
		analytics:     analytics,
		load:          loader,
	}
}

type PreparedServer struct {
	*Server
}

func (server *Server) PrepareRun() *PreparedServer {
	InitRoute(server.genericServer.Engine)
	server.InitRedisClient()

	server.shutdown.AddShutdownCallback(shutdown.CallBackFn(func(s string) error {
		log.Infof("manager: %s, genericServer shutdown callback called...\n", s)
		return server.genericServer.Close()
	}))

	server.shutdown.AddShutdownCallback(shutdown.CallBackFn(func(s string) error {
		log.Infof("manager: %s, analytics shutdown callback called...\n", s)
		server.analytics.Stop()
		return nil
	}))

	return &PreparedServer{
		Server: server,
	}
}

func (pserver *PreparedServer) Run() error {
	pserver.load.Run()
	pserver.analytics.Run()
	pserver.shutdown.Start()

	return pserver.genericServer.Run()
}

type GenericConfig struct {
	InsecureOp *option.InsecureOption
	SecureOp   *option.SecureOption
	RunOp      *option.RunOption
	Feature    *option.Feature
}

func NewGenericConfig(cfg *config.Config) *GenericConfig {
	return &GenericConfig{
		InsecureOp: cfg.InsecureOp,
		SecureOp:   cfg.SecureOp,
		RunOp:      cfg.RunOp,
		Feature:    cfg.Feature,
	}
}

func (genericCfg *GenericConfig) ApplyTo() (serverCfg *server.Config) {
	genericCfg.InsecureOp.ApplyTo(serverCfg)
	genericCfg.SecureOp.ApplyTo(serverCfg)
	genericCfg.RunOp.ApplyTo(serverCfg)
	genericCfg.Feature.ApplyTo(serverCfg)

	return
}

type ExtraConfig struct {
	GrpcClientOp *option.GrpcClientOption
	AnalyticsOp  *analytics.AnalyticOption
}

func NewExtraConfig(cfg *config.Config) *ExtraConfig {
	return &ExtraConfig{
		GrpcClientOp: cfg.GrpcClientOp,
		AnalyticsOp:  cfg.AnalyticsOp,
	}
}

type CompletedExtraConfig struct {
	*ExtraConfig
}

func (exCfg *ExtraConfig) Complete() *CompletedExtraConfig {
	return &CompletedExtraConfig{
		ExtraConfig: exCfg,
	}
}

func (completeExCfg *CompletedExtraConfig) New() {
	client := SetGrpcClient(completeExCfg.GrpcClientOp.Address, completeExCfg.GrpcClientOp.Cert)
	factory := grpcclient.GetGrpcClientInsOrDie(client)

	store.SetClient(factory)
	cache.GetCacheInsOr()

	analytics.NewAnalytics(completeExCfg.AnalyticsOp, RedisKeyPrefix)
}

func (server *Server) InitRedisClient() {
	ctx, cancel := context.WithCancel(context.Background())

	server.shutdown.AddShutdownCallback(shutdown.CallBackFn(func(s string) error {
		log.Infof("manager: %s, redis shutdown callback called...\n", s)
		cancel()
		return nil
	}))

	config := &storage.Config{
		Host:                  server.redisOp.Host,
		Port:                  server.redisOp.Port,
		Addrs:                 server.redisOp.Addrs,
		MasterName:            server.redisOp.MasterName,
		Username:              server.redisOp.Username,
		Password:              server.redisOp.Password,
		Database:              server.redisOp.Database,
		MaxIdle:               server.redisOp.MaxIdle,
		MaxActive:             server.redisOp.MaxActive,
		Timeout:               server.redisOp.Timeout,
		EnableCluster:         server.redisOp.EnableCluster,
		UseSSL:                server.redisOp.UseSSL,
		SSLInsecureSkipVerify: server.redisOp.SSLInsecureSkipVerify,
	}

	go storage.ConnectToRedis(ctx, config)
}

func SetGrpcClient(addr string, cert string) pb.CacheClient {
	caCert, err := os.ReadFile(cert)
	if err != nil {
		log.Fatalf("failed to read ca cert: %s.", err.Error())
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		log.Fatal("failed to append ca cert")
	}

	creds := credentials.NewTLS(&tls.Config{
		RootCAs:    certPool,
		ServerName: "iam-apiserver",
	})

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to create grpc conn: %s.", err.Error())
	}

	return pb.NewCacheClient(conn)
}
