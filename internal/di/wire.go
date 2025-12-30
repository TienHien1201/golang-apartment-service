package di

import (
	"thomas.vn/apartment_service/internal/config"
	"thomas.vn/apartment_service/internal/repository"
	"thomas.vn/apartment_service/internal/server/http/handler/ai"
	"thomas.vn/apartment_service/internal/server/http/handler/articles"
	xAuth "thomas.vn/apartment_service/internal/server/http/handler/auth"
	"thomas.vn/apartment_service/internal/server/http/handler/chatgroup"
	"thomas.vn/apartment_service/internal/server/http/handler/chatmessage"
	"thomas.vn/apartment_service/internal/server/http/handler/permission"
	"thomas.vn/apartment_service/internal/server/http/handler/root"
	xtotp "thomas.vn/apartment_service/internal/server/http/handler/totp"
	xuser "thomas.vn/apartment_service/internal/server/http/handler/user"
	queuejobs "thomas.vn/apartment_service/internal/server/queue/jobs"
	"thomas.vn/apartment_service/internal/usecase"
	auth2 "thomas.vn/apartment_service/internal/usecase/auth"
	"thomas.vn/apartment_service/internal/usecase/totp"
	"thomas.vn/apartment_service/internal/usecase/user"
	xcloudinary "thomas.vn/apartment_service/pkg/cloudinary"
	xfile "thomas.vn/apartment_service/pkg/file"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	mail "thomas.vn/apartment_service/pkg/mailer"
	xgoogle "thomas.vn/apartment_service/pkg/oauth/google"
	xqueue "thomas.vn/apartment_service/pkg/queue"
	ws "thomas.vn/apartment_service/pkg/websocket"
)

type AppContainer struct {
	HTTPHandler   xhttp.Handler
	InMemoryQueue *xqueue.InMemoryQueue
}

func NewAppContainer(cfg *config.Config, logger *xlogger.Logger) (*AppContainer, func(), error) {

	// === INIT DATABASE & SERVICES ===
	mysqlClient, err := cfg.InitMySQLDB()
	if err != nil {
		return nil, nil, err
	}

	redisCache, err := cfg.InitRedisCache()
	if err != nil {
		return nil, nil, err
	}

	esClient, err := cfg.InitElasticSearch()
	if err != nil {
		return nil, nil, err
	}

	httpClient := xhttp.NewHTTPClient()
	inMemoryQueue := xqueue.NewInMemoryQueue(logger, nil)

	fileSvc := xfile.NewHTTPFile(httpClient.HTTPClient())
	aiURLConfig := cfg.Ai

	tokenCfg := config.TokenConfig{AccessSecret: cfg.JWT.AccessSecret, AccessExpire: cfg.JWT.AccessExpire, RefreshSecret: cfg.JWT.RefreshSecret, RefreshExpire: cfg.JWT.RefreshExpire}

	mailer := mail.NewMailer(mail.SMTPConfig{Host: "smtp.gmail.com", Port: "587", User: "phamtienhien08072018@gmail.com", Pass: "gqqrdaxprykasskf"}, "Hien CNTT <tienhien.cntt@gmail.com>")
	googleOAuth := xgoogle.New(cfg.Auth.Google.ClientID, cfg.Auth.Google.ClientSecret, cfg.Auth.Google.CallbackURL)
	cld, _ := xcloudinary.NewCloudinary(cfg.Cloudinary)

	// === REPOSITORIES ===
	userRepo := repository.NewUserRepository(logger, mysqlClient.DB)
	aiRepo := repository.NewAiRepository(logger, httpClient, fileSvc, aiURLConfig.URL)
	permissionRepo := repository.NewPermissionRepository(logger, mysqlClient.DB)
	chatMessageRepo := repository.NewChatMessageRepository(logger, mysqlClient.DB)
	chatGroupRepo := repository.NewChatGroupRepository(logger, mysqlClient.DB)
	tokenSvc := usecase.NewToken(tokenCfg)
	articlesRepo := repository.NewArticlesRepository(logger, mysqlClient.DB)

	// === USECASES ===
	userUC := user.NewUserUsecase(logger, userRepo, redisCache, fileSvc, inMemoryQueue)
	chatMessageUC := usecase.NewChatMessageUsecase(logger, chatMessageRepo)
	chatGroupUc := usecase.NewChatGroupUsecase(logger, chatGroupRepo)
	authUC := auth2.NewAuthUsecase(logger, userRepo, tokenSvc, inMemoryQueue)
	aiUC := usecase.NewAiUsecase(logger, *aiRepo, aiURLConfig.DownloadURL, inMemoryQueue)
	permissionUC := usecase.NewPermissionUsecase(logger, permissionRepo)
	mailUC := usecase.NewMailUsecase(mailer)
	totpUc := totp.NewTotpUsecase(logger, userRepo)
	chatWsUC := usecase.NewChatUcase(logger, chatGroupUc, chatMessageUC)
	articleUc := usecase.NewArticlesUsecase(logger, articlesRepo)

	// === HANDLERS ===
	userHandler := xuser.NewHandler(logger, xuser.WithUserUsecase(userUC))
	chatMessageHandler := chatmessage.NewHandler(logger, chatmessage.WithChatMessageUsecase(chatMessageUC))
	chatGroupHandler := chatgroup.NewHandler(logger, chatgroup.WithChatGroupUsecase(chatGroupUc))
	authHandler := xAuth.NewHandler(logger, xAuth.WithGoogleOAuth(googleOAuth), xAuth.WithAuthUsecase(authUC))
	aiHandler := ai.NewAiHandler(logger, aiUC)
	authMiddlewareHandler := xAuth.NewAuthMiddleware(logger, tokenSvc, userRepo)
	permissionMiddlewareHandler := permission.NewPermissionMiddleware(logger, permissionUC)
	tOtpHandler := xtotp.NewHandler(logger, xtotp.WithTotpUsecase(totpUc))
	articleHandler := articles.NewHandler(logger, articles.WithArticleUsecase(articleUc))
	permissionHandler := permission.NewHandler(logger, permission.WithPermissionUsecase(permissionUC))
	hub := ws.NewHub()
	wsServer := &ws.Server{Hub: hub, ChatUC: chatWsUC, Token: tokenSvc}
	wsHandler := ws.NewHandler(wsServer)

	// === HTTP ROOT HANDLER ===
	httpHandler := root.NewHTTPHandler(
		logger,
		userHandler,
		authHandler,
		aiHandler,
		authMiddlewareHandler,
		permissionMiddlewareHandler,
		tOtpHandler,
		chatMessageHandler,
		chatGroupHandler,
		wsHandler,
		articleHandler,
		permissionHandler,
	)

	//========= Create job ==============
	mailJob := queuejobs.NewMailJob(logger, mailUC)
	uploadLocalAvatarJob := queuejobs.NewUploadUserAvatarJob(logger, fileSvc, userUC)
	uploadCloudAvatarJob := queuejobs.NewUploadAvatarCloudJob(logger, cld, userUC)
	deleteCloudAssetJob := queuejobs.NewDeleteCloudinaryAssetJob(logger, cld)
	inMemoryQueue.RegisterJobs([]xqueue.Job{mailJob, uploadLocalAvatarJob, uploadCloudAvatarJob, deleteCloudAssetJob})

	if err := inMemoryQueue.Start(); err != nil {
		return nil, nil, err
	}
	// === CLEANUP FUNCTION ===
	cleanup := func() {
		if err := mysqlClient.Close(); err != nil {
			logger.Error("Close MySQL client failed", xlogger.Error(err))
		}
		if err := redisCache.Close(); err != nil {
			logger.Error("Close Redis cache failed", xlogger.Error(err))
		}
		if err := esClient.Close(); err != nil {
			logger.Error("Close ElasticSearch client failed", xlogger.Error(err))
		}
	}

	return &AppContainer{
		HTTPHandler:   httpHandler,
		InMemoryQueue: inMemoryQueue,
	}, cleanup, nil
}
