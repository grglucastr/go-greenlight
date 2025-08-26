package main

import (
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

		//exit the application with a 0 (success) status code
		os.Exit(0)
	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	return srv.ListenAndServe()
}
