package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/apiserver/config"
)

type Config struct {
	SecureService   SecureInfo
	InsecureService InsecureInfo
	Middlewares     []string
	Timeout         time.Duration
	Mode            string
	Healthz         bool
	Version         bool
	EnableMetrics   bool
	EnableProfiling bool
}

type SecureInfo struct {
	Address string
	CertKey TLSInfo
}

type InsecureInfo struct {
	Address string
}

type TLSInfo struct {
	CertFile string
	KeyFile  string
}

func NewConfig() *Config {
	return &Config{
		SecureService:   SecureInfo{Address: "127.0.0.1:80", CertKey: TLSInfo{}},
		InsecureService: InsecureInfo{Address: "127.0.0.1:443"},
		Middlewares:     []string{"recovery", "logger"},
		Timeout:         10 * time.Second,
		Healthz:         true,
		Version:         true,
		EnableMetrics:   true,
		EnableProfiling: true,
	}
}

func (c *Config) Apply(cfg *config.Config) *Config {
	c.SecureService.Address = cfg.SecureOp.Address
	c.SecureService.CertKey.CertFile = cfg.SecureOp.TLSConfig.CertKey.CertFile
	c.SecureService.CertKey.KeyFile = cfg.SecureOp.TLSConfig.CertKey.CertKey
	c.InsecureService.Address = cfg.InsecureOp.Address
	c.Middlewares = cfg.RunOp.Middlewares
	c.Timeout = time.Duration(cfg.RunOp.TimeOut)
	c.Mode = cfg.RunOp.Mode
	c.Healthz = cfg.RunOp.Healthz
	c.Version = cfg.RunOp.Version
	c.EnableMetrics = cfg.Feature.EnableMetrics
	c.EnableProfiling = cfg.Feature.EnableProfiling

	return c
}

type CompletedConfig struct {
	*Config
}

func (c *Config) Complete() *CompletedConfig {
	return &CompletedConfig{c}
}

func (cc *CompletedConfig) New() *GenericServer {
	gin.SetMode(cc.Mode)

	gs := &GenericServer{
		SecureService:   cc.SecureService,
		InsecureService: cc.InsecureService,
		middlewares:     cc.Middlewares,
		timeout:         cc.Timeout,
		healthz:         cc.Healthz,
		version:         cc.Version,
		enableMetrics:   cc.EnableMetrics,
		enableProfiling: cc.EnableProfiling,
		Engine:          gin.New(),
		secureServer:    new(http.Server),
		insecureServer:  new(http.Server),
	}

	InitGenericServer(gs)

	return gs
}
