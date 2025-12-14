package di

import (
	"thomas.vn/hr_recruitment/internal/config"
	"thomas.vn/hr_recruitment/internal/repository"
	"thomas.vn/hr_recruitment/internal/server/http/handler"
	xAuth "thomas.vn/hr_recruitment/internal/server/http/handler/auth"
	xuser "thomas.vn/hr_recruitment/internal/server/http/handler/user"
	"thomas.vn/hr_recruitment/internal/usecase"
	"thomas.vn/hr_recruitment/pkg/auth"
	xfile "thomas.vn/hr_recruitment/pkg/file"
	xhttp "thomas.vn/hr_recruitment/pkg/http"
	xlogger "thomas.vn/hr_recruitment/pkg/logger"
	xqueue "thomas.vn/hr_recruitment/pkg/queue"
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

	// === TOKEN SERVICE (JWT) ===
	tokenSvc := auth.NewToken(
		cfg.JWT.AccessSecret,
		cfg.JWT.AccessExpire,
		cfg.JWT.RefreshSecret,
		cfg.JWT.RefreshExpire,
	)

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

	// === HTTP ROOT HANDLER ===
	httpHandler := handler.NewHTTPHandler(
		logger,
		userHandler,
		authHandler,
		aiHandler,
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
