package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/middleware"
	"github.com/neee333ko/component-base/pkg/core"
	"github.com/neee333ko/component-base/pkg/version"
	"github.com/neee333ko/log"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/sync/errgroup"
)

type GenericServer struct {
	SecureService   SecureInfo
	InsecureService InsecureInfo
	middlewares     []string
	timeout         time.Duration
	healthz         bool
	version         bool
	enableMetrics   bool
	enableProfiling bool
	*gin.Engine
	secureServer   *http.Server
	insecureServer *http.Server
}

func InitGenericServer(s *GenericServer) {
	s.InstallMiddlewares()
	s.InstallAPI()
}

func (s *GenericServer) InstallMiddlewares() {
	for i, m := range s.middlewares {
		if fn, ok := middleware.Middlewares[m]; !ok {
			log.Warnf("middleware: %s not found\n", s.middlewares[i])
		} else {
			log.Infof("install %s middleware\n", s.middlewares[i])
			s.Use(fn)
		}
	}
}

func (s *GenericServer) InstallAPI() {
	if s.healthz {
		s.GET("/healthz", func(ctx *gin.Context) {
			core.WriteResponse(ctx, nil, map[string]string{"status": "ok"})
		})
	}

	if s.enableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}

	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	if s.version {
		s.GET("/version", func(ctx *gin.Context) {
			core.WriteResponse(ctx, nil, version.Get())
		})
	}
}

func (s *GenericServer) Run() {
	s.insecureServer = &http.Server{
		Addr:    s.InsecureService.Address,
		Handler: s.Engine,
	}

	s.secureServer = &http.Server{
		Addr:    s.SecureService.Address,
		Handler: s.Engine,
	}

	var eg errgroup.Group

	eg.Go(func() error {
		log.Infof("http server starting, listening on %s\n", s.insecureServer.Addr)

		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		log.Infof("https server starting, listening on %s\n", s.secureServer.Addr)

		if err := s.secureServer.ListenAndServeTLS(s.SecureService.CertKey.CertFile, s.SecureService.CertKey.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if s.healthz {
		go s.ping(ctx)
	}

	err := eg.Wait()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (s *GenericServer) ping(ctx context.Context) {
	var addr string
	if strings.Contains(s.InsecureService.Address, "0.0.0.0") {
		addr = fmt.Sprintf("127.0.0.1:%s/healthz", strings.Split(s.InsecureService.Address, ":")[1])
	} else {
		addr = s.InsecureService.Address + "/healthz"
	}

	for {
		resp, err := http.Get(addr)

		if err == nil && resp.Status == "200 OK" {
			resp.Body.Close()
			log.Info("ping succeed!")
			return
		}

		time.Sleep(time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("ping failed")
		default:
			log.Info("retry ping...")
		}
	}
}

func (s *GenericServer) Close() error {
	var eg errgroup.Group

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	eg.Go(func() error {
		if err := s.insecureServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		if err := s.secureServer.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	})

	return eg.Wait()
}
