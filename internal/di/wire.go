package di

import (
	"thomas.vn/apartment_service/internal/config"
	"thomas.vn/apartment_service/internal/repository"
	"thomas.vn/apartment_service/internal/server/http/handler"
	xAuth "thomas.vn/apartment_service/internal/server/http/handler/auth"
	xuser "thomas.vn/apartment_service/internal/server/http/handler/user"
	"thomas.vn/apartment_service/internal/usecase"
	"thomas.vn/apartment_service/pkg/auth"
	xfile "thomas.vn/apartment_service/pkg/file"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xmiddleware "thomas.vn/apartment_service/pkg/http/middleware"
	xlogger "thomas.vn/apartment_service/pkg/logger"
	xqueue "thomas.vn/apartment_service/pkg/queue"
)

type AppContainer struct {
	HTTPHandler xhttp.Handler
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

	// === REPOSITORIES ===
	userRepo := repository.NewUserRepository(logger, mysqlClient.DB)
	aiRepo := repository.NewAiRepository(logger, httpClient, fileSvc, aiURLConfig.URL)

	tokenCfg := auth.Config{
		AccessSecret:  cfg.JWT.AccessSecret,
		AccessExpire:  cfg.JWT.AccessExpire,
		RefreshSecret: cfg.JWT.RefreshSecret,
		RefreshExpire: cfg.JWT.RefreshExpire,
	}

	tokenSvc := auth.NewToken(tokenCfg)

	// === USECASES ===
	userUC := usecase.NewUserUsecase(logger, userRepo, redisCache)
	authUC := usecase.NewAuthUsecase(logger, userRepo, tokenSvc)
	aiUC := usecase.NewAiUsecase(logger, *aiRepo, aiURLConfig.DownloadURL, inMemoryQueue)

	// === HANDLERS ===
	userHandler := xuser.NewHandler(logger, xuser.WithUserUsecase(userUC))

	authHandler := xAuth.NewHandler(
		logger,
		xAuth.WithAuthUsecase(authUC),
	)

	aiHandler := handler.NewAiHandler(logger, aiUC)

	authMiddleware := xmiddleware.NewAuthMiddleware(
		logger,
		tokenSvc,
		userRepo,
	)

	// === HTTP ROOT HANDLER ===
	httpHandler := handler.NewHTTPHandler(
		logger,
		userHandler,
		authHandler,
		aiHandler,
		authMiddleware,
	)

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
		HTTPHandler: httpHandler,
	}, cleanup, nil
}
