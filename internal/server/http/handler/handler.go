package handler

import (
	"github.com/labstack/echo/v4"
	xAuth "thomas.vn/apartment_service/internal/server/http/handler/auth"
	xmiddleware "thomas.vn/apartment_service/pkg/http/middleware"

	xuser "thomas.vn/apartment_service/internal/server/http/handler/user"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type handler struct {
	logger               *xlogger.Logger
	user                 *xuser.Handler
	auth                 *xAuth.Handler
	ai                   *AiHandler
	authMiddleware       *xmiddleware.AuthMiddleware
	permissionMiddleware *xmiddleware.PermissionMiddleware
}

func NewHTTPHandler(logger *xlogger.Logger, user *xuser.Handler, auth *xAuth.Handler, ai *AiHandler, authMiddleware *xmiddleware.AuthMiddleware, permissionMiddleware *xmiddleware.PermissionMiddleware) xhttp.Handler {
	return &handler{
		logger:               logger,
		user:                 user,
		auth:                 auth,
		ai:                   ai,
		authMiddleware:       authMiddleware,
		permissionMiddleware: permissionMiddleware,
	}
}

func (h *handler) HealthCheck(c echo.Context) error {
	return c.String(200, "OK")
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api")

	// Base routes
	api.GET("/health", h.HealthCheck)

	// User routes
	h.registerUserRoutes(api)

	// AI routes
	h.registerAiRoutes(api)

	//	Auth routes
	h.registerAuthRoutes(api)

}

func (h *handler) registerUserRoutes(e *echo.Group) {
	users := e.Group("/users")
	{
		users.POST("", h.user.User().Create)
		users.GET("/:id", h.user.User().Get)
		users.PUT("/:id", h.user.User().Update)
		users.DELETE("/:id", h.user.User().Delete)
		users.GET("", h.user.User().List, h.authMiddleware.Protect)
	}
}

func (h *handler) registerAiRoutes(e *echo.Group) {
	ai := e.Group("/ai")
	{
		ai.POST("/scan-cv", h.ai.VerifyCV)
	}
}

func (h *handler) registerAuthRoutes(e *echo.Group) {
	auth := e.Group("/auth")
	{
		auth.POST("/register", h.auth.Auth().Register)
		auth.POST("/login", h.auth.Auth().Login)
		auth.POST("/refresh", h.auth.Auth().Refresh)
		auth.POST("/logout", h.auth.Auth().Logout)
		auth.GET("/get-info", h.auth.Auth().GetInfo, h.authMiddleware.Protect, h.permissionMiddleware.Check)
		auth.POST("/refresh-token", h.auth.Auth().Refresh)
		auth.GET("/google", h.auth.Auth().GoogleLogin)
		auth.GET("/google/callback", h.auth.Auth().GoogleCallback)
	}
}
