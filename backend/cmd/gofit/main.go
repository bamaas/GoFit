package main

import (
	"log/slog"
	"net/http"
	"time"
	"database/sql"

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

func setupDB(logger *slog.Logger) (*sql.DB, error) {
	logger.Debug("Initializing database...")

	db, err := sql.Open("sqlite", ":memory:")
	// db, err := sql.Open("sqlite", "gofit.db")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	createCheckinsTableQuery := `
	CREATE TABLE IF NOT EXISTS checkins (
	uuid STRING NOT NULL PRIMARY KEY,
	datetime INTEGER NOT NULL,
	weight FLOAT NOT NULL,
	notes STRING
	);`
	_, err = db.Exec(createCheckinsTableQuery)
	if err != nil {
		return nil, err
	}

	createUsersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
		email CITEXT UNIQUE NOT NULL,
		password_hash BYTEA NOT NULL,
		activated BOOL NOT NULL,
		version INTEGER NOT NULL DEFAULT 1
	);
	`
	_, err = db.Exec(createUsersTableQuery)
	if err != nil {
		return nil, err
	}

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

	// Database
	db, err := setupDB(logger)
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
	err = app.models.CheckIns.InjectSampleData()
	if err != nil {
		panic(err)
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
