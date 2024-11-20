package http

import (
	"net/http"
	"go.uber.org/zap"
	"github.com/gorilla/mux"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
)

type CategoryEndpoint struct {
	categoryUC usecase.CategoryUseCase
	logger     *zap.Logger
}

func NewCategoryEndpoint(categoryUC usecase.CategoryUseCase, logger *zap.Logger) *CategoryEndpoint {
	return &CategoryEndpoint{
		categoryUC: categoryUC,
		logger:     logger,
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
	categories, err := e.categoryUC.Get()
	if err != nil {	
		e.sendError(w, http.StatusInternalServerError, err, "error getting categories", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, categories)
}

func (e *CategoryEndpoint) sendError(w http.ResponseWriter, statusCode int, err error, context string, additionalInfo map[string]string) {
	e.logger.Error(err.Error(), zap.String("context", context), zap.Any("info", additionalInfo))
	utils.SendErrorResponse(w, statusCode, err.Error())
}
