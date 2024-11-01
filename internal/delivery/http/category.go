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

func (e *CategoryEndpoints) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := e.CategoryUseCase.GetCategories()
	if err != nil {
		e.logger.Error("Error getting categories", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, categories)
}
