package postgres

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	createQuestionQuery         = `INSERT INTO question (title, description, page, trigger_value, lower_description, upper_description, parent_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	getQuestionsByPageTypeQuery = `SELECT * FROM question WHERE page = $1`
	getAllQuestionsQuery        = `SELECT * FROM question`
)

type QuestionDB struct {
	db      DBExecutor
	logger  *zap.Logger
	ctx     context.Context
	timeout time.Duration
}

func NewQuestionRepository(db *pgxpool.Pool, logger *zap.Logger, ctx context.Context, timeout time.Duration) (repository.QuestionRepository, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &QuestionDB{db: db, logger: logger, ctx: ctx, timeout: timeout}, nil
}

func (q *QuestionDB) Create(question *entity.Question) error {
	var questionDB entity.Question

	ctx, cancel := context.WithTimeout(q.ctx, q.timeout)
	defer cancel()

	err := q.db.QueryRow(ctx, createQuestionQuery, question.Title, question.Description, question.Page, question.TriggerValue, question.LowerDescription, question.UpperDescription, question.ParentID).Scan(&questionDB.ID)
	if err != nil {
		q.logger.Error("failed to create question", zap.Error(err))
		return entity.PSQLWrap(err, err)
	}

	return nil
}

func (q *QuestionDB) GetByPageType(pageType entity.PageType) ([]entity.Question, error) {
	ctx, cancel := context.WithTimeout(q.ctx, q.timeout)
	defer cancel()

	rows, err := q.db.Query(ctx, getQuestionsByPageTypeQuery, pageType)
	if err != nil {
		q.logger.Error("failed to get questions by page type", zap.Error(err))
		return nil, entity.PSQLWrap(err, err)
	}
	defer rows.Close()

	var questions []entity.Question
	for rows.Next() {
		var question entity.Question
		if err := rows.Scan(&question.ID, &question.Title, &question.Description, &question.Page, &question.TriggerValue, &question.LowerDescription, &question.UpperDescription, &question.ParentID); err != nil {
			q.logger.Error("failed to scan question", zap.Error(err))
			return nil, entity.PSQLWrap(err, err)
		}
		questions = append(questions, question)
	}

	return questions, nil
}

func (q *QuestionDB) GetAll() ([]entity.Question, error) {
	ctx, cancel := context.WithTimeout(q.ctx, q.timeout)
	defer cancel()

	rows, err := q.db.Query(ctx, getAllQuestionsQuery)
	if err != nil {
		q.logger.Error("failed to get all questions", zap.Error(err))
		return nil, entity.PSQLWrap(err, err)
	}
	defer rows.Close()

	var questions []entity.Question
	for rows.Next() {
		var question entity.Question
		if err := rows.Scan(&question.ID, &question.Title, &question.Description, &question.Page, &question.TriggerValue, &question.LowerDescription, &question.UpperDescription, &question.ParentID); err != nil {
			q.logger.Error("failed to scan question", zap.Error(err))
			return nil, entity.PSQLWrap(err, err)
		}
		questions = append(questions, question)
	}

	return questions, nil
}
