package http

import (
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/static"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"io"
)

var (
	ErrInvalidStaticID    = errors.New("invalid static file ID")
	ErrStaticFileNotFound = errors.New("static file not found")
	ErrFailedToGetStatic  = errors.New("failed to get static file")
)

type StaticEndpoint struct {
	staticGrpcClient  static.StaticGrpcClient
	logger        *zap.Logger
}

func NewStaticEndpoint(staticGrpcClient static.StaticGrpcClient, 
	logger *zap.Logger) *StaticEndpoint {
	return &StaticEndpoint{
		staticGrpcClient: staticGrpcClient,
		logger:        logger,
	}
}

func (h *StaticEndpoint) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/files/{fileId}", h.GetById).Methods("GET")
	router.HandleFunc("/api/v1/files/stream/{fileId}", h.GetFileStream).Methods("GET")
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
func (h *StaticEndpoint) GetById(writer http.ResponseWriter, r *http.Request) {
	staticIdStr := mux.Vars(r)["fileId"]
	staticId, err := uuid.Parse(staticIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid static ID", nil)
		return
	}

	staticURL, err := h.staticGrpcClient.GetStatic(staticId)
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

// GetFileStream godoc
// @Summary Get static file stream by ID
// @Description Get a static file as a byte stream by its ID
// @Tags static
// @Produce octet-stream
// @Param fileId path string true "File ID"
// @Success 200 {binary} []byte "Static file content"
// @Failure 400 {object} utils.ErrResponse "Invalid file ID"
// @Failure 404 {object} utils.ErrResponse "Static file not found"
// @Failure 500 {object} utils.ErrResponse "Failed to get static file"
// @Router /api/v1/files/stream/{fileId} [get]
func (h *StaticEndpoint) GetFileStream(writer http.ResponseWriter, r *http.Request) {
	fileIdStr := mux.Vars(r)["fileId"]
	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		h.sendError(writer, http.StatusBadRequest, err, "invalid file ID", nil)
		return
	}

	filePath, err := h.staticGrpcClient.GetStatic(fileId)
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to get static file path", nil)
		return
	}

	fileStream, err := h.staticGrpcClient.GetStaticFile(filePath)
	if err != nil {
		if errors.Is(err, ErrStaticFileNotFound) {
			h.sendError(writer, http.StatusNotFound, err, "static file not found", nil)
		} else {
			h.sendError(writer, http.StatusInternalServerError, err, "failed to get static file", nil)
		}
		return
	}

	writer.Header().Set("Content-Type", "image/webp")
	writer.WriteHeader(http.StatusOK)

	_, err = io.Copy(writer, fileStream) 
	if err != nil {
		h.sendError(writer, http.StatusInternalServerError, err, "failed to write file stream", nil)
	}
}

func (e *StaticEndpoint) sendError(w http.ResponseWriter, statusCode int, err error, context string, additionalInfo map[string]string) {
	e.logger.Error(err.Error(), zap.String("context", context), zap.Any("info", additionalInfo))
	utils.SendErrorResponse(w, statusCode, err.Error())
}
