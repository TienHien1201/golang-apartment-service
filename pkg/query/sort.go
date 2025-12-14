package query

type SortOptions struct {
	SortBy  string `query:"sort_by"`
	OrderBy string `query:"order_by"`
}
