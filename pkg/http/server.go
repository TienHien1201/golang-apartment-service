package xhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	xmiddleware "thomas.vn/apartment_service/pkg/http/middleware"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type Server struct {
	e       *echo.Echo
	logger  *xlogger.Logger
	host    string
	port    int
	handler Handler
}

func NewHTTPServer(logger *xlogger.Logger, host string, port int, h Handler) *Server {
	e := echo.New()

	// Setup basic middleware
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.OPTIONS, echo.HEAD, echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderAccessControlAllowOrigin},
	}))

	// Setup custom middleware
	e.Use(xmiddleware.RequestLogging(logger))

	// Customize Echo server
	e.HideBanner = true
	e.HidePort = true

	// Setup routes
	h.RegisterRoutes(e)

	return &Server{
		e:       e,
		logger:  logger,
		host:    host,
		port:    port,
		handler: h,
	}
}

func (s *Server) Start() error {
	go func() {
		addr := fmt.Sprintf("%s:%d", s.host, s.port)
		if err := s.e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("http server error", xlogger.Error(err))
		}
	}()

	s.logger.Info("http server started", xlogger.String("host", s.host), xlogger.Int("port", s.port))
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.e.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down http server: %w", err)
	}

	s.logger.Info("http server stopped gracefully")
	return nil
}
