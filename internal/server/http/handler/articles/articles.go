package articles

import (
	"encoding/json"
	"net/url"

	"github.com/labstack/echo/v4"
	"thomas.vn/apartment_service/internal/domain/model"
	"thomas.vn/apartment_service/internal/domain/usecase"
	xhttp "thomas.vn/apartment_service/pkg/http"
	xlogger "thomas.vn/apartment_service/pkg/logger"
)

type ArticleHandler struct {
	logger     *xlogger.Logger
	articlesUC usecase.ArticlesUsecase
}

func NewArticleHandler(logger *xlogger.Logger, articlesUC usecase.ArticlesUsecase) *ArticleHandler {
	return &ArticleHandler{
		logger:     logger,
		articlesUC: articlesUC,
	}
}

// List godoc
// @Summary List articles
// @Description List articles with pagination and filters
// @Tags articles
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param sort_by query string false "Sort by field"
// @Param order_by query string false "Order by asc/desc"
// @Param filters query string false "JSON encoded filters"
// @Success 200 {object} xhttp.APIResponse{data=[]model.Articles}
// @Failure 400 {object} xhttp.APIResponse400Err{}
// @Failure 500 {object} xhttp.APIResponse500Err{}
// @Router /api/articles [get]
func (h *ArticleHandler) List(c echo.Context) error {
	var req model.ListArticleRequest
	if err := xhttp.ReadAndValidateRequest(c, &req); err != nil {
		return xhttp.BadRequestResponse(c, err)
	}
	var filters model.ArticlesFilters

	if req.Filters != "" && req.Filters != "undefined" {
		decoded, err := url.QueryUnescape(req.Filters)
		if err != nil {
			return xhttp.BadRequestResponse(c, err)
		}
		if err := json.Unmarshal([]byte(decoded), &filters); err != nil {
			return xhttp.BadRequestResponse(c, err)
		}
	}

	res, total, err := h.articlesUC.ListArticles(
		c.Request().Context(),
		&req,
		&filters)
	if err != nil {
		return xhttp.AppErrorResponse(c, err)
	}
	return xhttp.PaginationListResponse(c, &req.PaginationOptions, res, total)
}
