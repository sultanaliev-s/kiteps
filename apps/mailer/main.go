package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sultanaliev-s/kiteps/apps/mailer/config"
	"github.com/sultanaliev-s/kiteps/apps/mailer/domain"
	"github.com/sultanaliev-s/kiteps/apps/mailer/transport/httprest"
	"github.com/sultanaliev-s/kiteps/pkg/logging"
	"github.com/sultanaliev-s/kiteps/pkg/validation"
)

var isDev = flag.Bool("dev", false, "development mode")

func main() {
	flag.Parse()

	cfg, err := config.NewConfig(*isDev)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("logger created")

	validator := validation.NewValidator()

	service := domain.NewService(cfg.Addr, cfg.Password, cfg.From, cfg.Port, validator, logger)
	logger.Info("service created")

	server := httprest.NewServer(service, logger, validator, cfg.ServerAddress)
	logger.Info("server created")

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server start error: ", logging.Error("err", err))
		}
	}()

	logger.Info("server started", logging.String("address", cfg.ServerAddress))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit // await quit signal

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Fatal("server shutdown error: ", logging.Error("err", err))
	}

	logger.Warn("shutting down server...")
}
