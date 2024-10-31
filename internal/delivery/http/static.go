package http

import (
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var (
	ErrInvalidStaticID    = errors.New("invalid static file ID")
	ErrStaticFileNotFound = errors.New("static file not found")
	ErrFailedToGetStatic  = errors.New("failed to get static file")
)

type StaticEndpoints struct {
	StaticUseCase usecase.StaticUseCase
	logger        *zap.Logger
}

func NewStaticEndpoints(staticUseCase usecase.StaticUseCase, logger *zap.Logger) *StaticEndpoints {
	return &StaticEndpoints{
		StaticUseCase: staticUseCase,
		logger:        logger,
	}
}

func (h *StaticEndpoints) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/static/{staticId}", h.GetStaticById).Methods("GET")
}

// GetStaticById godoc
// @Summary Get static file by ID
// @Description Get a static file by its ID
// @Tags static
// @Produce json
// @Param staticId path string true "Static file ID"
// @Success 200 {string} string "URL of the static file"
// @Failure 400 {object} utils.ErrResponse "Invalid static ID"
// @Failure 404 {object} utils.ErrResponse "Static file not found"
// @Failure 500 {object} utils.ErrResponse "Failed to get static file"
// @Router /api/v1/static/{staticId} [get]
func (h *StaticEndpoints) GetStaticById(writer http.ResponseWriter, r *http.Request) {
	staticIdStr := mux.Vars(r)["staticId"]
	staticId, err := strconv.Atoi(staticIdStr)
	if err != nil || staticId <= 0 {
		h.logger.Error("invalid static ID", zap.Error(err))
		utils.SendErrorResponse(writer, http.StatusBadRequest, ErrInvalidStaticID.Error())
		return
	}

	staticURL, err := h.StaticUseCase.GetStaticURL(uuid.MustParse(staticIdStr))
	if err != nil {
		h.logger.Error("failed to get static file", zap.Error(err))
		if errors.Is(err, ErrStaticFileNotFound) {
			h.logger.Error("static file not found", zap.Error(err))
			utils.SendErrorResponse(writer, http.StatusNotFound, ErrStaticFileNotFound.Error())
		} else {
			h.logger.Error("failed to get static file", zap.Error(err))
			utils.SendErrorResponse(writer, http.StatusInternalServerError, ErrFailedToGetStatic.Error())
		}
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, staticURL)
}
