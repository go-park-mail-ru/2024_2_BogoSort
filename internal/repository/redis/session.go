package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const userSessionPlaceholder = "user_sessions:"

type SessionDB struct {
	rdb              *redis.Client
	sessionAliveTime int
	ctx              context.Context
	logger           *zap.Logger
}

func NewSessionRepository(rdb *redis.Client, sessionAliveTime int, ctx context.Context, logger *zap.Logger) (*SessionDB, error) {
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &SessionDB{
		rdb:              rdb,
		sessionAliveTime: sessionAliveTime,
		ctx:              ctx,
		logger:           logger,
	}, nil
}

func (s *SessionDB) Create(userID uuid.UUID) (string, error) {
	sessionID := uuid.NewString()

	for {
		_, err := s.rdb.Get(s.ctx, sessionID).Result()
		if errors.Is(err, redis.Nil) {
			break
		}
		sessionID = uuid.NewString()
	}

	err := s.rdb.Set(s.ctx, sessionID, userID.String(), time.Duration(s.sessionAliveTime)*time.Second).Err()
	if err != nil {
		s.logger.Error("error creating session", zap.String("sessionID", sessionID), zap.String("userID", userID.String()), zap.Error(err))
		return "", entity.RedisWrap(repository.ErrSessionCreationFailed, err)
	}

	err = s.rdb.SAdd(s.ctx, userSessionPlaceholder+userID.String(), sessionID).Err()
	if err != nil {
		s.logger.Error("error adding session to user", zap.String("sessionID", sessionID), zap.String("userID", userID.String()), zap.Error(err))
		return "", entity.RedisWrap(repository.ErrSessionCreationFailed, err)
	}
	return sessionID, nil
}

func (s *SessionDB) Get(sessionID string) (uuid.UUID, error) {
	userID, err := s.rdb.Get(s.ctx, sessionID).Result()
	if errors.Is(err, redis.Nil) {
		s.logger.Error("session not found", zap.String("sessionID", sessionID))
		return uuid.Nil, entity.RedisWrap(repository.ErrSessionNotFound, err)
	}
	if err != nil {
		s.logger.Error("error getting session", zap.String("sessionID", sessionID), zap.Error(err))
		return uuid.Nil, entity.RedisWrap(repository.ErrSessionCheckFailed, err)
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Error("error parsing userID", zap.String("userID", userID), zap.Error(err))
		return uuid.Nil, entity.RedisWrap(repository.ErrIncorrectID, err)
	}
	return id, nil
}

func (s *SessionDB) Delete(sessionID string) error {
	userID, err := s.rdb.Get(s.ctx, sessionID).Result()
	if errors.Is(err, redis.Nil) {
		s.logger.Info("session not found", zap.String("sessionID", sessionID))
		return nil
	}
	if err != nil {
		s.logger.Error("error getting session", zap.String("sessionID", sessionID), zap.Error(err))
		return entity.RedisWrap(repository.ErrSessionDeleteFailed, err)
	}

	err = s.rdb.Del(s.ctx, sessionID).Err()
	if err != nil {
		s.logger.Error("error deleting session", zap.String("sessionID", sessionID), zap.Error(err))
		return entity.RedisWrap(repository.ErrSessionDeleteFailed, err)
	}

	if userID == "" {
		s.logger.Error("userID is empty", zap.String("sessionID", sessionID))
		return entity.RedisWrap(repository.ErrSessionDeleteFailed, errors.New("userID is empty"))
	}

	err = s.rdb.SRem(s.ctx, userSessionPlaceholder+userID, sessionID).Err()
	if err != nil {
		s.logger.Error("error deleting session from user", zap.String("sessionID", sessionID), zap.String("userID", userID), zap.Error(err))
		return entity.RedisWrap(repository.ErrSessionDeleteFailed, err)
	}
	s.logger.Info("session deleted", zap.String("sessionID", sessionID), zap.String("userID", userID))
	return nil
}
