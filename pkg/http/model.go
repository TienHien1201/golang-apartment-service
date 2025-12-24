package xhttp

type APIResponse struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"OK"`
	Data    interface{} `json:"data,omitempty"`
}

type APIResponse400Err struct {
	Status  int               `json:"status" example:"400"`
	Message string            `json:"message" example:"Bad Request"`
	Data    []ValidationError `json:"data,omitempty"`
}

type APIResponse401Err struct {
	Status  int    `json:"status" example:"401"`
	Message string `json:"message" example:"Unauthorized"`
	Data    string `json:"data,omitempty" example:"Token is invalid"`
}

type APIResponse500Err struct {
	Status  int    `json:"status" example:"500"`
	Message string `json:"message" example:"Internal Server Error"`
	Data    string `json:"data,omitempty" example:"Something went wrong"`
}

type OldAPIResponse struct {
	Message string      `json:"message" example:"OK"`
	Status  int         `json:"status" example:"1"`
	Code    int         `json:"code" example:"200"`
	Data    interface{} `json:"data,omitempty"`
}

type OldAPIResponse400Err struct {
	Message string `json:"message" example:"Bad Request"`
	Status  int    `json:"status" example:"0"`
	Code    int    `json:"code" example:"400"`
}

type OldAPIResponse500Err struct {
	Message string `json:"message" example:"Internal Server Error"`
	Status  int    `json:"status" example:"0"`
	Code    int    `json:"code" example:"500"`
	Data    string `json:"data,omitempty" example:"Something went wrong"`
}

type ValidationError struct {
	Code    string `json:"code,omitempty" example:"ERR_REQUIRED"`
	Field   string `json:"field,omitempty" example:"name"`
	Message string `json:"message,omitempty" example:"Name is required"`
}

type ListDataResponse struct {
	Rows  interface{} `json:"rows"`
	Total int64       `json:"total,omitempty"`
}

type PaginationOptions struct {
	Page         int  `query:"page" default:"1" validate:"omitempty,min=1"`
	Limit        int  `query:"limit" default:"10" validate:"omitempty,min=1,max=100"`
	ExcludeTotal bool `query:"exclude_total" validate:"omitempty"`
}

type DateRangeOptions struct {
	FromDate string `query:"from_date" validate:"omitempty,datetime=2006-01-02"`
	ToDate   string `query:"to_date" validate:"omitempty,datetime=2006-01-02"`
	RangeBy  string `query:"range_by" default:"created_at" validate:"omitempty,oneof=created_at updated_at"`
}

type SortOptions struct {
	SortBy  string `query:"sort_by" default:"created_at" validate:"omitempty,oneof=created_at updated_at"`
	OrderBy string `query:"order_by" default:"desc" validate:"omitempty,oneof=asc desc"`
}

type PaginationResponse struct {
	Page      int         `json:"page"`
	PageSize  int         `json:"pageSize"`
	TotalItem int64       `json:"totalItem"`
	TotalPage int64       `json:"totalPage"`
	Items     interface{} `json:"items"`
}
