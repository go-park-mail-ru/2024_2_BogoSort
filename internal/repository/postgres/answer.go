package postgres

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"go.uber.org/zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
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
		ID:          dbAnswer.ID,
		Value:       dbAnswer.Value,
		QuestionID:  dbAnswer.QuestionID,
		UserID:      dbAnswer.UserID,
	}, nil
}