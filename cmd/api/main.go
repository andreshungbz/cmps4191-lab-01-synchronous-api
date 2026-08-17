package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/andreshungbz/cmps4191-lab-01-synchronous-api/internal/data"
	"github.com/andreshungbz/cmps4191-lab-01-synchronous-api/internal/vcs"
	_ "github.com/lib/pq"
)

var (
	version = vcs.Version()
)

// config stores the API server configuration.
type config struct {
	port    int    // API server port
	env     string // (development|staging|production)
	timeout int    // grace period in seconds for HTTP server shutdown
	db      struct {
		dsn string // data source name
	}
}

// application holds the dependencies for the HTTP handlers, helpers, middleware,
// etc. so that they are all accessible through dependency injection.
type application struct {
	config config
	logger *slog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// FLAGS

	// server flags
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.IntVar(&cfg.timeout, "timeout", 30, "Grace period in seconds for HTTP server shutdown")

	// database flags
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")

	// version flag
	displayVersion := flag.Bool("version", false, "Display program version")

	flag.Parse()

	// display program version and exit if the version flag was passed
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	// DATABASE

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Database connection pool established")

	// APPLICATION

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	// start the API server
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// openDB connects to the PostgreSQL database using the provided DSN and
// and returns a pointer to a handler to that database.
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// test the connection with a ping
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
