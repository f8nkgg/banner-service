package app

import (
	"banner/config"
	v1 "banner/internal/controller/http/v1"
	"banner/internal/repository"
	"banner/internal/service"
	"banner/pkg/cache"
	"banner/pkg/db/postgres"
	"banner/pkg/httpserver"
	"banner/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	pgURL, _ := os.LookupEnv("PG_URL")
	pg, err := postgres.New(pgURL, cfg.PG.PoolMax, cfg.PG.ConnAttempts, cfg.PG.ConnTimeout)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %v", err))
	}
	defer pg.Close()

	memCache := cache.NewMemoryCache(1000, 20)

	bannerController := v1.NewBannerController(
		service.NewBannerService(
			repository.NewBannerRepository(pg),
			memCache,
			5*time.Minute,
		),
		l,
	)

	handler := gin.New()
	v1.RegisterRoutes(handler, bannerController)
	httpServer := httpserver.New(handler, cfg.HTTPServer.ReadTimeout, cfg.HTTPServer.WriteTimeout, cfg.HTTPServer.Host, cfg.HTTPServer.Port, cfg.HTTPServer.MaxHeaderBytes, cfg.HTTPServer.ShutdownTimeout)
	l.Info("Server is starting on " + cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify: %v", err)
	}
	l.Info("Server shutting down...")
	dropTables(pgURL, l)
	err = httpServer.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown: %v", err)
	}

}
