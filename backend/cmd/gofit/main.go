package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/bamaas/gofit/internal/assets"
	"github.com/bamaas/gofit/internal/config"
	"github.com/bamaas/gofit/internal/data"
	"github.com/bamaas/gofit/internal/logger"
)

const version = "0.0.2"

type application struct {
	config   config.Config
	logger   *slog.Logger
	models   data.Models
}

func setupDB(logger *slog.Logger, developmentMode bool) (*sql.DB, error) {
	logger.Debug("Initializing database...")

	datasourceName := "/data/gofit.db"
	if (developmentMode) {
		datasourceName = ":memory:"
	}
	db, err := sql.Open("sqlite", datasourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Apply migrations
	fsDriver, err := iofs.New(assets.MigrationsFs, assets.MigrationsPath)
    if err != nil {
        return nil, err
    }
	dbDriver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return nil, err
	}
	migrator, err := migrate.NewWithInstance("iofs", fsDriver, "sqlite", dbDriver)
	if err != nil {
		return nil, err
	}
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	logger.Info("database migrations applied")

	return db, nil
}

func main() {

	// Retrieve config
	cfg, err := config.Retrieve()
	if err != nil {
		panic(err)
	}

	// Setup logger
	logger, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	logger.Info(cfg.Users[0].Email)
	logger.Info(cfg.Users[0].Password)
	logger.Info(cfg.Users[1].Email)
	logger.Info(cfg.Users[1].Password)

	// Database
	db, err := setupDB(logger, cfg.DevelopmentMode)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Init application
	app := &application{
		config:   *cfg,
		logger:   logger,
		models:   data.NewModels(db, logger),
	}

	// Bootstrap
	err = app.Bootstrap()
	if err != nil {
		panic(err)
	}

	// Start server
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
