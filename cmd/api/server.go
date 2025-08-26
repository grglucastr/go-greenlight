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

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	// Shutdown channel. Use to receive any errors returned
	// by the graceful Shutdown() function.
	shutdownError := make(chan error)

	go func() {
		// creates a quit buffered channel with size 1,
		// which carries os.Signal values
		quit := make(chan os.Signal, 1)

		// listen for incoming SIGINT and SIGTERM signals
		// and relay them to the quit channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// read the signal from the quit channel.
		// this code will block until a signal is received
		s := <-quit

		app.logger.Info("caught signal", "signal", s.String())

		// give any in-flight request a period of 30 seconds to complete before
		// the application termination
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// We wait to receive a return value from Shutdown() on the shutdownError channel
	// If the return value is an error, we know that there was a problem with the
	// graceful shutdown and we return the error.
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", srv.Addr)

	return nil
}
