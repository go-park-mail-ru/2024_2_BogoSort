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

func NewSessionRepository(rdb *redis.Client, sessionAliveTime int, logger *zap.Logger) *SessionDB {
	return &SessionDB{
		rdb:              rdb,
		sessionAliveTime: sessionAliveTime,
		ctx:              context.Background(),
		logger:           zap.L(),
	}
}

func (sdb *SessionDB) CreateSession(userID uuid.UUID) (string, error) {
	sessionID := uuid.NewString()

	for {
		_, err := sdb.rdb.Get(sdb.ctx, sessionID).Result()
		if errors.Is(err, redis.Nil) {
			break
		}
		sessionID = uuid.NewString()
	}

	err := sdb.rdb.Set(sdb.ctx, sessionID, userID.String(), time.Duration(sdb.sessionAliveTime)*time.Second).Err()
	if err != nil {
		sdb.logger.Error("error creating session", zap.String("sessionID", sessionID), zap.String("userID", userID.String()), zap.Error(err))
		return "", entity.RedisWrap(repository.ErrSessionCreationFailed, err)
	}

	err = sdb.rdb.SAdd(sdb.ctx, userSessionPlaceholder+userID.String(), sessionID).Err()
	if err != nil {
		sdb.logger.Error("error adding session to user", zap.String("sessionID", sessionID), zap.String("userID", userID.String()), zap.Error(err))
		return "", entity.RedisWrap(repository.ErrSessionCreationFailed, err)
	}
	sdb.logger.Info("session created", zap.String("sessionID", sessionID), zap.String("userID", userID.String()))
	return sessionID, nil
}

func (sdb *SessionDB) GetSession(sessionID string) (uuid.UUID, error) {
	userID, err := sdb.rdb.Get(sdb.ctx, sessionID).Result()
	if errors.Is(err, redis.Nil) {
		sdb.logger.Error("session not found", zap.String("sessionID", sessionID))
		return uuid.Nil, entity.RedisWrap(repository.ErrSessionNotFound, err)
	}
	if err != nil {
		sdb.logger.Error("error getting session", zap.String("sessionID", sessionID), zap.Error(err))
		return uuid.Nil, entity.RedisWrap(repository.ErrSessionCheckFailed, err)
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		sdb.logger.Error("error parsing userID", zap.String("userID", userID), zap.Error(err))
		return uuid.Nil, entity.RedisWrap(repository.ErrIncorrectID, err)
	}
	sdb.logger.Info("session found", zap.String("sessionID", sessionID), zap.String("userID", id.String()))
	return id, nil
}

func (sdb *SessionDB) DeleteSession(sessionID string) error {
	userID, err := sdb.rdb.Get(sdb.ctx, sessionID).Result()
	if errors.Is(err, redis.Nil) {
		sdb.logger.Info("session not found", zap.String("sessionID", sessionID))
		return nil
	}
	if err != nil {
		sdb.logger.Error("error getting session", zap.String("sessionID", sessionID), zap.Error(err))
		return entity.RedisWrap(repository.ErrSessionDeleteFailed, err)
	}

	err = sdb.rdb.Del(sdb.ctx, sessionID).Err()
	if err != nil {
		sdb.logger.Error("error deleting session", zap.String("sessionID", sessionID), zap.Error(err))
		return entity.RedisWrap(repository.ErrSessionDeleteFailed, err)
	}

	if userID == "" {
		sdb.logger.Error("userID is empty", zap.String("sessionID", sessionID))
		return entity.RedisWrap(repository.ErrSessionDeleteFailed, errors.New("userID is empty"))
	}

	err = sdb.rdb.SRem(sdb.ctx, userSessionPlaceholder+userID, sessionID).Err()
	if err != nil {
		sdb.logger.Error("error deleting session from user", zap.String("sessionID", sessionID), zap.String("userID", userID), zap.Error(err))
		return entity.RedisWrap(repository.ErrSessionDeleteFailed, err)
	}
	sdb.logger.Info("session deleted", zap.String("sessionID", sessionID), zap.String("userID", userID))
	return nil
}
