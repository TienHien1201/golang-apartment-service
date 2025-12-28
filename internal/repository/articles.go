package repository

import (
	"context"

	"gorm.io/gorm"
	"thomas.vn/apartment_service/internal/domain/model"
	"thomas.vn/apartment_service/internal/domain/repository"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type articlesRepository struct {
	logger        *xlogger.Logger
	articlesTable *gorm.DB
}

func NewArticlesRepository(logger *xlogger.Logger, db *gorm.DB) repository.ArticlesRepository {
	return &articlesRepository{
		logger:        logger,
		articlesTable: db.Table("articles"),
	}
}

func (r *articlesRepository) ListArticles(ctx context.Context, req *model.ListArticleRequest, filters *model.ArticlesFilters) ([]*model.Articles, int64, error) {
	var articles []*model.Articles
	var total int64

	query := r.articlesTable.WithContext(ctx)
	if filters.Views != 0 {
		query = query.Where("views = ?", filters.Views)
	}
	if filters.Content != "" {
		query = query.Where("content LIKE ?", "%"+filters.Content+"%")
	}
	if filters.ID != 0 {
		query = query.Where("id = ?", filters.ID)
	}
	if req.IsDeleted != 0 {
		query = query.Where("is_deleted = ?", req.IsDeleted)
	}
	if req.FromDate != "" {
		query = query.Where(req.RangeBy+" >= ?", req.FromDate+" 00:00:00")
	}
	if req.ToDate != "" {
		query = query.Where(req.RangeBy+" <= ?", req.ToDate+" 23:59:59")
	}
	// Get total count if not exclude
	if !req.ExcludeTotal {
		if err := query.Count(&total).Error; err != nil {
			r.logger.Error("Count articles failed", xlogger.Error(err))
			return nil, 0, err
		}
	}

	// Apply pagination
	if req.Page > 0 && req.Limit > 0 {
		query = query.Offset((req.Page - 1) * req.Limit).Limit(req.Limit)
	}

	// Apply sorting
	if req.SortBy != "" && req.OrderBy != "" {
		query = query.Order(req.SortBy + " " + req.OrderBy)
	}

	// Execute query
	if err := query.Find(&articles).Error; err != nil {
		r.logger.Error("List articles failed", xlogger.Error(err))
		return nil, 0, err
	}
	return articles, total, nil
}
