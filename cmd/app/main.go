package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"thomas.vn/apartment_service/internal/config"
)

var (
	// Name is the name of the application.
	Name = "app"
	// Version is the version of the application.
	Version = "1.0.0"
	// Env is the environment of the application.
	Env = "dev"
)

func init() {
	flag.StringVar(&Name, "name", Name, "Name")
	flag.StringVar(&Version, "version", Version, "Version")
	flag.StringVar(&Env, "env", Env, "Environment")
}

// @title Apartment Business Service API
// @version 1.0
// @description This is an Apartment business system API
// @BasePath /
// @host localhost:1424
// @schemes http

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
// @description Type "Basic" followed by a space and base64 encoded username and password

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// Load configuration for the specified environment
	cfg, err := config.LoadConfig(config.Environment(Env), "./config")
	if err != nil {
		log.Panicf("Failed to load configuration: %v", err)
	}

	// Override application name, version and environment
	cfg.App.Name = Name
	cfg.App.Version = Version
	cfg.App.Env = Env
	log.Printf("\033[1;36mStarting\033[0m \033[1;33m%s\033[0m \033[1;32mv%s\033[0m (\033[1;35m%s\033[0m)", Name, Version, Env)

	// Initialize application with config
	app, cleanup, err := NewApp(cfg)
	if err != nil {
		log.Panicf("Failed to initialize application: %v", err)
	}
	defer cleanup()

	// Start the application
	if err := app.Start(); err != nil {
		log.Panicf("Application error: %v", err)
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop the application
	if err := app.Stop(ctx); err != nil {
		log.Panicf("Failed to stop application: %v", err)
	}
}
