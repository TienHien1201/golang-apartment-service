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

	var (
		articles []*model.Articles
		total    int64
	)

	query := r.articlesTable.WithContext(ctx)

	query = r.applyFilters(query, req, filters)

	if !req.ExcludeTotal {
		if err := r.count(query, &total); err != nil {
			return nil, 0, err
		}
	}

	query = r.applyPaging(query, req)
	query = r.applySorting(query, req)

	if err := query.Find(&articles).Error; err != nil {
		r.logger.Error("List articles failed", xlogger.Error(err))
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *articlesRepository) applyFilters(query *gorm.DB, req *model.ListArticleRequest, filters *model.ArticlesFilters) *gorm.DB {

	if filters == nil {
		return query
	}

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

	query = r.applyDateRange(query, req)

	return query
}

func (r *articlesRepository) applyDateRange(query *gorm.DB, req *model.ListArticleRequest) *gorm.DB {

	if req.FromDate != "" {
		query = query.Where(req.RangeBy+" >= ?", req.FromDate+" 00:00:00")
	}

	if req.ToDate != "" {
		query = query.Where(req.RangeBy+" <= ?", req.ToDate+" 23:59:59")
	}

	return query
}

func (r *articlesRepository) applyPaging(query *gorm.DB, req *model.ListArticleRequest) *gorm.DB {

	if req.Page > 0 && req.Limit > 0 {
		return query.Offset((req.Page - 1) * req.Limit).Limit(req.Limit)
	}

	return query
}

func (r *articlesRepository) applySorting(query *gorm.DB, req *model.ListArticleRequest) *gorm.DB {

	if req.SortBy != "" && req.OrderBy != "" {
		return query.Order(req.SortBy + " " + req.OrderBy)
	}

	return query
}

func (r *articlesRepository) count(query *gorm.DB, total *int64) error {
	if err := query.Count(total).Error; err != nil {
		r.logger.Error("Count articles failed", xlogger.Error(err))
		return err
	}
	return nil
}
