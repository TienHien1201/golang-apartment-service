package di

import (
	"thomas.vn/apartment_service/internal/config"
	"thomas.vn/apartment_service/internal/repository"
	"thomas.vn/apartment_service/internal/server/http/handler"
	xAuth "thomas.vn/apartment_service/internal/server/http/handler/auth"
	xuser "thomas.vn/apartment_service/internal/server/http/handler/user"
	queuejobs "thomas.vn/apartment_service/internal/server/queue/jobs"
	"thomas.vn/apartment_service/internal/usecase"
	auth2 "thomas.vn/apartment_service/internal/usecase/auth"
	"thomas.vn/apartment_service/internal/usecase/user"
	"thomas.vn/apartment_service/pkg/auth"
	xcloudinary "thomas.vn/apartment_service/pkg/cloudinary"
	xfile "thomas.vn/apartment_service/pkg/file"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xmiddleware "thomas.vn/apartment_service/pkg/http/middleware"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	mail "thomas.vn/apartment_service/pkg/mailer"
	xgoogle "thomas.vn/apartment_service/pkg/oauth/google"
	xqueue "thomas.vn/apartment_service/pkg/queue"
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

	tokenCfg := auth.Config{
		AccessSecret:  cfg.JWT.AccessSecret,
		AccessExpire:  cfg.JWT.AccessExpire,
		RefreshSecret: cfg.JWT.RefreshSecret,
		RefreshExpire: cfg.JWT.RefreshExpire,
	}

	mailer := mail.NewMailer(
		mail.SMTPConfig{
			Host: "smtp.gmail.com",
			Port: "587",
			User: "phamtienhien08072018@gmail.com",
			Pass: "gqqrdaxprykasskf",
		},
		"Hien CNTT <tienhien.cntt@gmail.com>",
	)

	googleOAuth := xgoogle.New(
		cfg.Auth.Google.ClientID,
		cfg.Auth.Google.ClientSecret,
		cfg.Auth.Google.CallbackURL,
	)

	cld, err := xcloudinary.NewCloudinary(cfg.Cloudinary)
	if err != nil {
		logger.Error("failed to init cloudinary", xlogger.Error(err))
		return nil, nil, err
	}

	// === REPOSITORIES ===
	userRepo := repository.NewUserRepository(logger, mysqlClient.DB)
	aiRepo := repository.NewAiRepository(logger, httpClient, fileSvc, aiURLConfig.URL)
	permissionRepo := repository.NewPermissionRepository(logger, mysqlClient.DB)

	tokenSvc := auth.NewToken(tokenCfg)

	// === USECASES ===
	userUC := user.NewUserUsecase(logger, userRepo, redisCache, fileSvc, inMemoryQueue)
	authUC := auth2.NewAuthUsecase(logger, userRepo, tokenSvc, inMemoryQueue)
	aiUC := usecase.NewAiUsecase(logger, *aiRepo, aiURLConfig.DownloadURL, inMemoryQueue)
	permissionUC := usecase.NewPermissionUsecase(permissionRepo)
	mailUC := usecase.NewMailUsecase(mailer)

	// === HANDLERS ===
	userHandler := xuser.NewHandler(logger, xuser.WithUserUsecase(userUC))
	authHandler := xAuth.NewHandler(logger, xAuth.WithGoogleOAuth(googleOAuth), xAuth.WithAuthUsecase(authUC))
	aiHandler := handler.NewAiHandler(logger, aiUC)
	authMiddlewareHandler := xmiddleware.NewAuthMiddleware(logger, tokenSvc, userRepo)
	permissionMiddlewareHandler := xmiddleware.NewPermissionMiddleware(permissionUC)

	// === HTTP ROOT HANDLER ===
	httpHandler := handler.NewHTTPHandler(
		logger,
		userHandler,
		authHandler,
		aiHandler,
		authMiddlewareHandler,
		permissionMiddlewareHandler,
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
