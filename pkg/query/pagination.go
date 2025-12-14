package query

type PaginationOptions struct {
	Page         int  `query:"page" validate:"omitempty,gt=0"`
	Limit        int  `query:"limit" validate:"omitempty,gt=0"`
	ExcludeTotal bool `query:"exclude_total"`
}
