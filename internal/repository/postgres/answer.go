package postgres

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AnswerDB struct {
	DB      DBExecutor
	logger  *zap.Logger
	ctx     context.Context
	timeout time.Duration
}

func NewAnswerRepository(db *pgxpool.Pool,
	logger *zap.Logger,
	ctx context.Context,
	timeout time.Duration) (repository.AnswerRepository, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &AnswerDB{
		DB:      db,
		logger:  logger,
		ctx:     ctx,
		timeout: timeout,
	}, nil
}

const insertAnswerQuery = `
	INSERT INTO answer (value, question_id, user_id)
	VALUES ($1, $2, $3)
	RETURNING id, value, question_id, user_id
`

const getAnswersByQuestionIDQuery = `
	SELECT * FROM answer WHERE question_id = $1
`

func (r *AnswerDB) Add(answer *entity.Answer) (*entity.Answer, error) {
	var dbAnswer entity.Answer

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	r.logger.Info("adding answer", zap.Any("answer", answer))

	err := r.DB.QueryRow(ctx, insertAnswerQuery,
		answer.Value,
		answer.QuestionID,
		answer.UserID).Scan(
		&dbAnswer.ID,
		&dbAnswer.Value,
		&dbAnswer.QuestionID,
		&dbAnswer.UserID,
	)

	if err != nil {
		r.logger.Error("error adding answer", zap.Error(err))
		return nil, entity.PSQLWrap(err)
	}

	r.logger.Info("answer added", zap.Any("answer", dbAnswer))

	return &entity.Answer{
		ID:         dbAnswer.ID,
		Value:      dbAnswer.Value,
		QuestionID: dbAnswer.QuestionID,
		UserID:     dbAnswer.UserID,
	}, nil
}

func (r *AnswerDB) GetByQuestionID(questionID string) ([]entity.Answer, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	rows, err := r.DB.Query(ctx, getAnswersByQuestionIDQuery, questionID)
	if err != nil {
		r.logger.Error("failed to get answers by question ID", zap.Error(err))
		return nil, entity.PSQLWrap(err, err)
	}
	defer rows.Close()

	var answers []entity.Answer
	var createdAt time.Time
	for rows.Next() {
		var answer entity.Answer
		if err := rows.Scan(&answer.ID, &answer.QuestionID, &answer.UserID, &answer.Value, &createdAt); err != nil {
			r.logger.Error("failed to scan answer", zap.Error(err))
			return nil, entity.PSQLWrap(err, err)
		}
		answers = append(answers, answer)
	}

	return answers, nil
}
