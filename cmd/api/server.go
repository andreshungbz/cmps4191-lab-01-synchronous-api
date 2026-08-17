package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serve starts the HTTP server and monitors for shutdown signals so that the server
// and background tasks can gracefully terminate.
func (app *application) serve() error {
	// HTTP Server Configuration

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	// Goroutine for gracefully shutting down HTTP Server on SIGINT (Ctrl + C) and SIGTERM (pkill)
	// and finishing any background tasks

	shutdownError := make(chan error)
	go func() {
		// create a single buffered channel that blocks until a signal is received
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("Shutting down server", "signal", s.String())

		// allow HTTP server to close any remaining connections with the configured timeout period
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.config.timeout)*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("Completing background tasks", "addr", srv.Addr)

		// block until all goroutines are finished
		app.wg.Wait()
		shutdownError <- nil
	}()

	// Run HTTP Server

	app.logger.Info("Starting server", "addr", srv.Addr, "env", app.config.env)

	// we should expect the ErrServerClosed error since Shutdown was called.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// block until shutdown signal is received
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("Stopped server", "addr", srv.Addr)

	return nil
}
