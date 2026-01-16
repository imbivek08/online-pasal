package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imbivek08/hamropasal/internal/config"
	"github.com/imbivek08/hamropasal/internal/database"
	"github.com/imbivek08/hamropasal/internal/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo   *echo.Echo
	db     *database.Database
	config *config.Config
}

func New(cfg *config.Config, db *database.Database) *Server {
	return &Server{
		echo:   echo.New(),
		db:     db,
		config: cfg,
	}
}

func (s *Server) Start() error {
	// Configure Echo
	s.echo.HideBanner = true
	s.echo.HidePort = false

	// Apply global middleware
	s.setupMiddleware()

	// Setup routes
	router.SetupRoutes(s.echo, s.db, s.config)

	// Start server with graceful shutdown
	return s.startWithGracefulShutdown()
}

func (s *Server) setupMiddleware() {
	// Logger middleware
	s.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${method} ${uri} (${latency_human})\n",
	}))

	// Recover middleware
	s.echo.Use(middleware.Recover())

	// CORS middleware
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Request ID middleware
	s.echo.Use(middleware.RequestID())

	// Timeout middleware
	s.echo.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))
}

func (s *Server) startWithGracefulShutdown() error {
	// Channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		address := fmt.Sprintf(":%s", s.config.ServerPort)
		s.echo.Logger.Info(fmt.Sprintf("Starting server on %s", address))
		serverErrors <- s.echo.Start(address)
	}()

	// Channel to listen for interrupt or terminate signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		s.echo.Logger.Info(fmt.Sprintf("Received signal: %v. Starting graceful shutdown...", sig))

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := s.echo.Shutdown(ctx); err != nil {
			s.echo.Logger.Error(fmt.Sprintf("Graceful shutdown failed: %v", err))
			return s.echo.Close()
		}

		s.echo.Logger.Info("Server stopped gracefully")
		return nil
	}
}

func (s *Server) Close() error {
	s.db.Close()
	return s.echo.Close()
}
