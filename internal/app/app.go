package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/faasf/functions-api/internal/config"
	v1 "github.com/faasf/functions-api/internal/controllers/http/v1"
	functionsRepo "github.com/faasf/functions-api/internal/repositories/functions"
	functionsService "github.com/faasf/functions-api/internal/services/functions"
	"github.com/faasf/functions-api/pkg/httpserver"
	"github.com/faasf/functions-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	var err error

	functionsService := functionsService.New(functionsRepo.New(cfg))

	handler := gin.New()
	v1.NewRouter(handler, l, cfg, functionsService)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
