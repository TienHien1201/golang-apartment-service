package repository

import (
	"context"

	"thomas.vn/apartment_service/internal/domain/model"
)

type ArticlesRepository interface {
	ListArticles(ctx context.Context, req *model.ListArticleRequest, filters *model.ArticlesFilters) ([]*model.Articles, int64, error)
}
