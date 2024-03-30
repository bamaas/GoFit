package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/bamaas/gofit/internal/config"
	"github.com/bamaas/gofit/internal/database"
	"github.com/bamaas/gofit/internal/logger"
)

const version = "0.0.2"

type application struct {
	config   config.Config
	logger   *slog.Logger
	database *database.Database
}

func main() {

	// Retrieve config
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	// Setup logger
	logger, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	// Database
	db, err := database.New(logger)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.InjectSampleData()
	if err != nil {
		panic(err)
	}

	// Init application
	app := &application{
		config:   *cfg,
		logger:   logger,
		database: db,
	}

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting server")
	if err = srv.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
