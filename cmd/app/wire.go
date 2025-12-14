package main

import (
	"context"
	"fmt"
	"sync"

	"thomas.vn/apartment_service/internal/config"
	"thomas.vn/apartment_service/internal/di"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xserver "thomas.vn/apartment_service/pkg/server"
)

type App struct {
	logger  *xlogger.Logger
	servers []xserver.Server
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config) (*App, func(), error) {
	// Initialize logger
	logger, err := initLogger(cfg)
	if err != nil {
		return nil, nil, err
	}

	// Initialize dependencies
	container, cleanup, err := di.NewAppContainer(cfg, logger)
	if err != nil {
		return nil, nil, err
	}

	// Initialize HTTP server
	httpServer := xhttp.NewHTTPServer(logger, cfg.Server.HTTP.Host, cfg.Server.HTTP.Port, container.HTTPHandler)

	return &App{
		logger:  logger,
		servers: []xserver.Server{httpServer},
	}, cleanup, nil

}

func (a *App) Start() error {
	for _, srv := range a.servers {
		if err := srv.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(a.servers))

	for _, srv := range a.servers {
		wg.Add(1)
		go func(s xserver.Server) {
			defer wg.Done()
			if err := s.Stop(ctx); err != nil {
				errChan <- err
			}
		}(srv)
	}

	// Wait for all servers to stop
	wg.Wait()
	close(errChan)

	// Collect errors
	errs := make([]error, 0, len(a.servers))
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors stopping servers: %v", errs)
	}
	return nil
}

func initLogger(cfg *config.Config) (*xlogger.Logger, error) {
	logCfg := &xlogger.Config{
		Level:      cfg.Logger.Level,
		Format:     cfg.Logger.Format,
		Output:     cfg.Logger.Output,
		TimeFormat: cfg.Logger.TimeFormat,
	}

	logger, err := xlogger.New(logCfg)
	if err != nil {
		return nil, err
	}

	return logger, nil
}
