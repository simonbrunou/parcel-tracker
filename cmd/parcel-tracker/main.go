package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/simonbrunou/parcel-tracker/internal/auth"
	"github.com/simonbrunou/parcel-tracker/internal/config"
	"github.com/simonbrunou/parcel-tracker/internal/handler"
	"github.com/simonbrunou/parcel-tracker/internal/notifier"
	"github.com/simonbrunou/parcel-tracker/internal/server"
	"github.com/simonbrunou/parcel-tracker/internal/store"
	"github.com/simonbrunou/parcel-tracker/internal/tracker"
	"github.com/simonbrunou/parcel-tracker/web"
)

var version = "dev"

func main() {
	cfg := config.Load()

	logLevel := slog.LevelInfo
	if cfg.Dev {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	logger.Info("starting parcel-tracker", "version", version, "port", cfg.Port)

	// Database
	db, err := store.NewSQLiteStore(cfg.DatabasePath)
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Auth
	a := auth.New(db)
	if cfg.Password != "" && !a.IsConfigured(context.Background()) {
		logger.Info("setting up initial password from environment")
		if err := a.Setup(context.Background(), cfg.Password); err != nil {
			logger.Error("failed to setup password", "error", err)
			os.Exit(1)
		}
	}

	// Tracker registry
	registry := tracker.NewRegistry()

	// Notifier
	n := &notifier.Notifier{
		Store:  db,
		Logger: logger,
	}
	if err := n.EnsureVAPIDKeys(context.Background()); err != nil {
		logger.Error("failed to initialize VAPID keys", "error", err)
		os.Exit(1)
	}

	// Handlers
	h := &handler.Handler{
		Store:    db,
		Auth:     a,
		Tracker:  registry,
		Logger:   logger,
		Notifier: n,
	}

	// Embedded frontend
	distFS, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		logger.Error("failed to load embedded frontend", "error", err)
		os.Exit(1)
	}

	// HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      server.New(h, a, distFS, logger),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Background tracking worker
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	if cfg.RefreshInterval > 0 {
		w := &tracker.Worker{
			Store:    db,
			Registry: registry,
			Interval: cfg.RefreshInterval,
			Logger:   logger,
			Notifier: n,
		}
		go w.Run(workerCtx)
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		logger.Info("shutting down...")
		workerCancel()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	logger.Info("listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
