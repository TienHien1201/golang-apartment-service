package xutils

import (
	"fmt"

	"gorm.io/gorm"
)

// BuildlikeQuery constructs a GORM with multiple OR conditions for LIKE searches.
// It applies the search pattern to each specified fied, e.g., name LIKE ? OR mail LIKE ....
func BuildLikeQuery(query *gorm.DB, searchPattern string, fields ...string) *gorm.DB {
	if len(fields) == 0 {
		return query
	}

	q := query
	for i, field := range fields {
		if i == 0 {
			q = q.Where(fmt.Sprintf("%s LIKE ?", field), searchPattern)
		} else {
			q = q.Or(fmt.Sprintf("%s LIKE ?", field), searchPattern)
		}
	}

	return q
}

// ApplyPagination adds OFFSET and LIMIT clauses to GORM query for pagination.
func ApplyPagination(query *gorm.DB, page, limit int) *gorm.DB {
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		return query.Offset(offset).Limit(limit)
	}

	return query
}

// Apply Sorting adds Orderby clause to query
func ApplySorting(query *gorm.DB, sortBy, orderBy string) *gorm.DB {
	if sortBy != "" && orderBy != "" {
		return query.Order(sortBy + " " + orderBy)
	}

	return query
}

// ApplyFilters applies filters to query
func ApplyInFilter(query *gorm.DB, field string, values interface{}) *gorm.DB {
	return query.Where(fmt.Sprintf("%s IN ?", field), values)
}

// ApplyEqualFilter adds WHERE = clause for single value
func ApplyEqualFilter(query *gorm.DB, field string, value interface{}) *gorm.DB {
	return query.Where(fmt.Sprintf("%s = ?", field), value)
}
