package http

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type CategoryEndpoint struct {
	categoryUC usecase.CategoryUseCase
}

func NewCategoryEndpoint(categoryUC usecase.CategoryUseCase) *CategoryEndpoint {
	return &CategoryEndpoint{
		categoryUC: categoryUC,
	}
}

func (e *CategoryEndpoint) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/categories", e.Get).Methods("GET")
}

// Get godoc
// @Summary Get all categories
// @Description Retrieve a list of all categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} entity.Category
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/categories [get]
func (e *CategoryEndpoint) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(context.Background())
	logger.Info("get categories request")
	categories, err := e.categoryUC.Get()
	if err != nil {
		e.sendError(w, http.StatusInternalServerError, err, "error getting categories", nil)
		return
	}
	logger.Info("get categories response", zap.Any("categories", categories))
	utils.SendJSONResponse(w, http.StatusOK, categories)
}

func (e *CategoryEndpoint) sendError(w http.ResponseWriter, statusCode int, err error, contextInfo string, additionalInfo map[string]string) {
	logger := middleware.GetLogger(context.Background())
	logger.Error(err.Error(), zap.String("context", contextInfo), zap.Any("info", additionalInfo))
	utils.SendErrorResponse(w, statusCode, err.Error())
}
