package query

type DateRangeOptions struct {
	FromDate string `query:"from_date"`
	ToDate   string `query:"to_date"`
	RangeBy  string `query:"range_by"`
}
