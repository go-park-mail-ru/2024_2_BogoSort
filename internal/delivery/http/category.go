package http

import (
	"net/http"
	"go.uber.org/zap"
	"github.com/gorilla/mux"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
)

type CategoryEndpoints struct {
	CategoryUseCase usecase.CategoryUseCase
	logger        *zap.Logger
}

func NewCategoryEndpoints(categoryUseCase usecase.CategoryUseCase, logger *zap.Logger) *CategoryEndpoints {
	return &CategoryEndpoints{
		CategoryUseCase: categoryUseCase,
		logger:        logger,
	}
}

func (e *CategoryEndpoints) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/categories", e.GetCategories).Methods("GET")
}

// GetCategories godoc
// @Summary Get all categories
// @Description Retrieve a list of all categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} entity.Category
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/categories [get]
func (e *CategoryEndpoints) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := e.CategoryUseCase.GetCategories()
	if err != nil {
		e.sendError(w, http.StatusInternalServerError, err, "error getting categories", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, categories)
}

func (e *CategoryEndpoints) sendError(w http.ResponseWriter, statusCode int, err error, context string, additionalInfo map[string]string) {
	e.logger.Error(err.Error(), zap.String("context", context), zap.Any("info", additionalInfo))
	utils.SendErrorResponse(w, statusCode, err.Error())
}
