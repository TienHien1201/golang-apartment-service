package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model"
	"thomas.vn/apartment_service/internal/domain/repository"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type articlesUsecase struct {
	logger            *xlogger.Logger
	articleRepository repository.ArticlesRepository
}

func NewArticlesUsecase(logger *xlogger.Logger, articleRepository repository.ArticlesRepository) usecase.ArticlesUsecase {
	return &articlesUsecase{
		logger:            logger,
		articleRepository: articleRepository,
	}
}

func (u *articlesUsecase) ListArticles(ctx context.Context, req *model.ListArticleRequest, filters *model.ArticlesFilters) ([]*model.Articles, int64, error) {
	return u.articleRepository.ListArticles(ctx, req, filters)
}
