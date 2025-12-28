package model

import (
	"time"

	"thomas.vn/apartment_service/pkg/query"
)

type Articles struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	ImageURL  string     `json:"image_url"`
	Views     int        `json:"views"`
	UserID    int        `json:"user_id"`
	DeletedBy int        `json:"deleted_by"`
	IsDeleted int        `json:"is_deleted"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type ListArticleRequest struct {
	query.PaginationOptions
	query.DateRangeOptions
	query.SortOptions
	IsDeleted int    `query:"is_deleted"`
	Filters   string `query:"filters"`
}

type ArticlesFilters struct {
	ID      int    `query:"id"`
	Content string `query:"content"`
	Views   int    `query:"views"`
}
