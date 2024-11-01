package http

import (
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
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
	router.HandleFunc("/api/v1/files/{fileId}", h.GetStaticById).Methods("GET")
}

// GetStaticById godoc
// @Summary Get file by ID
// @Description Get a file by its ID
// @Tags static
// @Produce json
// @Param fileId path string true "File ID"
// @Success 200 {string} string "URL of the static file"
// @Failure 400 {object} utils.ErrResponse "Invalid static ID"
// @Failure 404 {object} utils.ErrResponse "Static file not found"
// @Failure 500 {object} utils.ErrResponse "Failed to get static file"
// @Router /api/v1/files/{fileId} [get]
func (h *StaticEndpoints) GetStaticById(writer http.ResponseWriter, r *http.Request) {
	staticIdStr := mux.Vars(r)["fileId"]
	staticId, err := uuid.Parse(staticIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid static ID", nil)
		return
	}

	staticURL, err := h.StaticUseCase.GetStaticURL(staticId)
	if err != nil {
		if errors.Is(err, ErrStaticFileNotFound) {
			h.sendError(writer, http.StatusNotFound, err, "static file not found", nil)
		} else {
			h.sendError(writer, http.StatusInternalServerError, err, "failed to get static file", nil)
		}
		return
	}

	utils.SendJSONResponse(writer, http.StatusOK, staticURL)
}

func (e *StaticEndpoints) sendError(w http.ResponseWriter, statusCode int, err error, context string, additionalInfo map[string]string) {
	e.logger.Error(err.Error(), zap.String("context", context), zap.Any("info", additionalInfo))
	utils.SendErrorResponse(w, statusCode, err.Error())
}