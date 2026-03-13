package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shubhangcs/agromart-server/internal/app"
	"github.com/shubhangcs/agromart-server/internal/routes"
)

// @title           Agromart API
// @version         1.0
// @description     Agromart B2B agricultural marketplace backend server.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	application, err := app.NewApplication()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialise application: %v\n", err)
		os.Exit(1)
	}
	defer application.DB.Close()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           routes.SetupRoutes(application),
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	// Channel to receive OS shutdown signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine so we can listen for signals concurrently.
	serverErr := make(chan error, 1)
	go func() {
		application.Logger.Info("server listening", "port", port)
		serverErr <- server.ListenAndServe()
	}()

	// Block until a signal or a fatal server error arrives.
	select {
	case err = <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			application.Logger.Error("server error", "error", err)
			os.Exit(1)
		}
	case sig := <-quit:
		application.Logger.Info("received signal — shutting down gracefully", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err = server.Shutdown(ctx); err != nil {
			application.Logger.Error("graceful shutdown failed", "error", err)
		} else {
			application.Logger.Info("server stopped cleanly")
		}
	}
}
