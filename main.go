package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shubhangcs/agromart-server/internal/app"
	"github.com/shubhangcs/agromart-server/internal/env"
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
	port := env.GetInt("PORT", 8080)

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		application.Logger.Info("server listening", "port", port)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err = <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			application.Logger.Error("server error", "error", err)
			os.Exit(1)
		}
	case sig := <-quit:
		application.Logger.Info("shutting down", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		defer cancel()

		if err = server.Shutdown(ctx); err != nil {
			application.Logger.Error("graceful shutdown failed", "error", err)
		} else {
			application.Logger.Info("server stopped cleanly")
		}
	}
}
