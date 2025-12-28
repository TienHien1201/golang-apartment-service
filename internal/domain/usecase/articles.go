package usecase

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model"
)

type ArticlesUsecase interface {
	ListArticles(ctx context.Context, req *model.ListArticleRequest, filter *model.ArticlesFilters) ([]*model.Articles, int64, error)
}
