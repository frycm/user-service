package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardanlabs/conf/v3"
	"github.com/frycm/user-service/cmd/user-service/server"
	"go.uber.org/zap"
)

func main() {
	err := run()
	if err != nil {
		zap.S().Errorw("server run failed", "cause", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := zap.NewDevelopment() // or production
	if err != nil {
		return fmt.Errorf("logger initization failed: %w", err)
	}
	zap.ReplaceGlobals(logger)
	defer func() {
		_ = logger.Sync()
	}()

	const prefix = "USER_SERVICE"
	cfg := server.Conf{}
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	zap.S().Infow("user service", "conf", cfg)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdown
		zap.S().Info("shutdown signal received")
		cancel()
	}()

	return server.Serve(ctx, cfg)
}
