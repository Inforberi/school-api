package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"restapi/internal/app"
	"restapi/internal/config"
	applog "restapi/internal/log"
)

func main() {
	// init context
	ctx := context.Background()

	// init config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// init logger
	logger, err := applog.NewLogger(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	applog.SetDefault(logger)
	defer applog.Sync()

	// init signal
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// init app
	a := app.NewApp(cfg)

	// run app
	if err := a.RunUntil(ctx, 10*time.Second); err != nil {
		applog.Error("app failed", "err", err)
		os.Exit(1)
	}
}
