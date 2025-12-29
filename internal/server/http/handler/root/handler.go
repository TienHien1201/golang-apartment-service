package root

import (
	"github.com/labstack/echo/v4"
	handler2 "thomas.vn/apartment_service/internal/server/http/handler/ai"
	"thomas.vn/apartment_service/internal/server/http/handler/articles"
	xAuth "thomas.vn/apartment_service/internal/server/http/handler/auth"
	"thomas.vn/apartment_service/internal/server/http/handler/chatgroup"
	"thomas.vn/apartment_service/internal/server/http/handler/chatmessage"
	xtotp "thomas.vn/apartment_service/internal/server/http/handler/totp"
	xmiddleware "thomas.vn/apartment_service/pkg/http/middleware"
	ws "thomas.vn/apartment_service/pkg/websocket"

	xuser "thomas.vn/apartment_service/internal/server/http/handler/user"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type handler struct {
	logger               *xlogger.Logger
	user                 *xuser.Handler
	auth                 *xAuth.Handler
	ai                   *handler2.Handler
	authMiddleware       *xmiddleware.AuthMiddleware
	permissionMiddleware *xmiddleware.PermissionMiddleware
	totp                 *xtotp.Handler
	chatMessage          *chatmessage.Handler
	chatGroup            *chatgroup.Handler
	wsHandler            *ws.Handler
	article              *articles.Handler
}

func NewHTTPHandler(
	logger *xlogger.Logger,
	user *xuser.Handler,
	auth *xAuth.Handler,
	ai *handler2.Handler,
	authMiddleware *xmiddleware.AuthMiddleware,
	permissionMiddleware *xmiddleware.PermissionMiddleware,
	totp *xtotp.Handler,
	chatMessage *chatmessage.Handler,
	chatGroup *chatgroup.Handler,
	wsHandler *ws.Handler,
	article *articles.Handler,
) xhttp.Handler {
	return &handler{
		logger:               logger,
		user:                 user,
		chatMessage:          chatMessage,
		auth:                 auth,
		ai:                   ai,
		authMiddleware:       authMiddleware,
		permissionMiddleware: permissionMiddleware,
		totp:                 totp,
		chatGroup:            chatGroup,
		wsHandler:            wsHandler,
		article:              article,
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

	//	Totp routes
	h.registerTotpRoutes(api)

	//	Chat message routes
	h.registerChatMessageRoutes(api)

	//	Chat group routes
	h.registerChatGroupRoutes(api)

	// Articles routes
	h.registerArticleRoutes(api)

	// WebSocket
	e.GET("/ws", h.wsHandler.Handle())

	//Public
	e.Static("/attachments", "attachments")

}

func (h *handler) registerUserRoutes(e *echo.Group) {
	users := e.Group("/users")
	{
		users.POST("", h.user.User().Create)
		users.GET("/:id", h.user.User().Get)
		users.PUT("/:id", h.user.User().Update)
		users.DELETE("/:id", h.user.User().Delete)
		users.GET("", h.user.User().List, h.authMiddleware.Protect)
		users.POST("/avatar-local", h.user.User().UploadLocal, h.authMiddleware.Protect)
		users.POST("/avatar-cloud", h.user.User().UploadCloud, h.authMiddleware.Protect)
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

func (h *handler) registerTotpRoutes(e *echo.Group) {
	tOtp := e.Group("/totp")
	{
		tOtp.POST("/generate", h.totp.Totp().Generate, h.authMiddleware.Protect)
		tOtp.POST("/verify", h.totp.Totp().Verify, h.authMiddleware.Protect)
		tOtp.POST("/generate", h.totp.Totp().Generate, h.authMiddleware.Protect)
		tOtp.POST("/save", h.totp.Totp().Save, h.authMiddleware.Protect)
	}
}

func (h *handler) registerChatMessageRoutes(e *echo.Group) {
	chat := e.Group("/chat-message")
	{
		chat.GET("", h.chatMessage.ChatMessage().List, h.authMiddleware.Protect, h.permissionMiddleware.Check)
	}
}

func (h *handler) registerChatGroupRoutes(e *echo.Group) {
	chat := e.Group("/chat-group")
	{
		chat.GET("", h.chatGroup.ChatGroup().List, h.authMiddleware.Protect, h.permissionMiddleware.Check)
	}
}

func (h *handler) registerArticleRoutes(e *echo.Group) {
	chat := e.Group("/article")
	{
		chat.GET("/all", h.article.Articles().List, h.authMiddleware.Protect)
		chat.GET("", h.article.Articles().List, h.authMiddleware.Protect, h.permissionMiddleware.Check)
	}
}
