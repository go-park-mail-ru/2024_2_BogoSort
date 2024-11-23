package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/survey"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	ErrInvalidPage = errors.New("invalid page")
)

type SurveyEndpoint struct {
	surveyClient   survey.SurveyClient
	sessionManager *utils.SessionManager
	logger         *zap.Logger
}

func NewSurveyEndpoint(surveyClient survey.SurveyClient, sessionManager *utils.SessionManager, logger *zap.Logger) *SurveyEndpoint {
	return &SurveyEndpoint{
		surveyClient:   surveyClient,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

func (s *SurveyEndpoint) ConfigureProtectedRoutes(router *mux.Router) {
	protected := router.PathPrefix("/api/v1").Subrouter()
	sessionMiddleware := middleware.NewAuthMiddleware(s.sessionManager)
	protected.Use(sessionMiddleware.SessionMiddleware)

	protected.HandleFunc("/questions/{page}", s.GetQuestions).Methods(http.MethodGet)
	protected.HandleFunc("/answer", s.PostAnswer).Methods(http.MethodPost)
	protected.HandleFunc("/stats", s.GetStats).Methods(http.MethodGet)
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

	questions, err := s.surveyClient.GetQuestions(ctx, &dto.GetQuestionsRequest{Page: page})
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err, "failed to get questions", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, questions)
}

func (s *SurveyEndpoint) PostAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	if page == "" {
		s.sendError(w, http.StatusBadRequest, ErrInvalidPage, "page is required", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var answer entity.Answer
	err := json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		s.sendError(w, http.StatusBadRequest, err, "failed to decode answers", nil)
		return
	}

	message, err := s.surveyClient.AddAnswer(ctx, &dto.PostAnswersRequest{Answer: answer})
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err, "failed to add answer", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, message)
}

func (s *SurveyEndpoint) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats, err := s.surveyClient.GetStats(ctx)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err, "failed to get stats", nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, stats)
}
