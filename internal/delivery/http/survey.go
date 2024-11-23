package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	ErrInvalidPage = errors.New("invalid page")
)

type SurveyEndpoint struct {
	surveyGrpcClient survey.SurveyGrpcClient
	sessionManager   *utils.SessionManager
	logger           *zap.Logger
}

func NewSurveyEndpoint(surveyGrpcClient survey.SurveyGrpcClient, sessionManager *utils.SessionManager, logger *zap.Logger) *SurveyEndpoint {
	return &SurveyEndpoint{
		surveyGrpcClient: surveyGrpcClient,
		sessionManager:   sessionManager,
		logger:           logger,
	}
}

func (s *SurveyEndpoint) ConfigureProtectedRoutes(router *mux.Router) {
	protected := router.PathPrefix("/api/v1").Subrouter()
	sessionMiddleware := middleware.NewAuthMiddleware(s.sessionManager)
	protected.Use(sessionMiddleware.SessionMiddleware)

	protected.HandleFunc("/questions/{page}", s.GetQuestions).Methods(http.MethodGet)
	protected.HandleFunc("/answers/{page}", s.PostAnswers).Methods(http.MethodPost)
}

func (s *SurveyEndpoint) sendError(w http.ResponseWriter, statusCode int, err error, message string, details map[string]string) {
	s.logger.Error(message, zap.Error(err), zap.Any("details", details))
	utils.SendErrorResponse(w, statusCode, message)
}

func (s *SurveyEndpoint) GetQuestions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	if page == "" {
		s.sendError(w, http.StatusBadRequest, ErrInvalidPage, "page is required", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	questions, err := s.surveyGrpcClient.GetQuestions(ctx, &dto.GetQuestionsRequest{Page: page})
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err, "failed to get questions", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, questions)
}

func (s *SurveyEndpoint) PostAnswers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	if page == "" {
		s.sendError(w, http.StatusBadRequest, ErrInvalidPage, "page is required", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var answers []dto.Answer
	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		s.sendError(w, http.StatusBadRequest, err, "failed to decode answers", nil)
		return
	}

}
