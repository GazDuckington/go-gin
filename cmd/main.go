package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	database "github.com/GazDuckington/go-gin/db"
	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/routes"
)

func main() {
	cfg := config.LoadConfig()

	// connect DB (safe to skip in dev if you want; error out otherwise)
	if err := database.Connect(cfg); err != nil {
		cfg.Logger.Warnf("database connection failed, continuing without DB: %v", err)
	} else {
		cfg.Logger.Info("database connected successfully")
	}

	// NOTE: we manage schema with migrate CLI; DO NOT call AutoMigrate here in prod.
	// If you want to auto-migrate for quick dev, you can call it explicitly.

	r := routes.SetupRouter(cfg)
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	cfg.Logger.Infof("starting server on %s", addr)

	// graceful shutdown basic pattern
	go func() {
		if err := r.Run(addr); err != nil {
			cfg.Logger.Fatalf("server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cfg.Logger.Info("shutting down")
}
