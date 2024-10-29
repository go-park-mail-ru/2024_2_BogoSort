package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const userSessionPlaceholder = "user_sessions:"

type SessionDB struct {
	rdb              *redis.Client
	sessionAliveTime int
	ctx              context.Context
}

func NewSessionRepository(rdb *redis.Client, sessionAliveTime int) *SessionDB {
	return &SessionDB{
		rdb:              rdb,
		sessionAliveTime: sessionAliveTime,
		ctx:              context.Background(),
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

	err := sdb.rdb.Set(sdb.ctx, sessionID, userID, time.Duration(sdb.sessionAliveTime)*time.Second).Err()
	if err != nil {
		return "", entity.RedisWrap(repository.ErrSessionCreationFailed, err)
	}

	err = sdb.rdb.SAdd(sdb.ctx, userSessionPlaceholder+userID.String(), sessionID).Err()
	if err != nil {
		return "", entity.RedisWrap(repository.ErrSessionCreationFailed, err)
	}

	return sessionID, nil
}

func (sdb *SessionDB) GetSession(sessionID string) (uuid.UUID, error) {
	userID, err := sdb.rdb.Get(sdb.ctx, sessionID).Result()
	if errors.Is(err, redis.Nil) {
		return uuid.Nil, entity.RedisWrap(repository.ErrSessionNotFound, err)
	}
	if err != nil {
		return uuid.Nil, entity.RedisWrap(repository.ErrSessionCheckFailed, err)
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, entity.RedisWrap(repository.ErrIncorrectID, err)
	}
	return id, nil
}

func (sdb *SessionDB) DeleteSession(sessionID string) error {
	userID, err := sdb.rdb.Get(sdb.ctx, sessionID).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return entity.RedisWrap(repository.ErrSessionDeleteFailed, err)
	}

	err = sdb.rdb.Del(sdb.ctx, sessionID).Err()
	if err != nil {
		return entity.RedisWrap(repository.ErrSessionDeleteFailed, err)
	}

	err = sdb.rdb.SRem(sdb.ctx, userSessionPlaceholder+userID).Err()
	if err != nil {
		return entity.RedisWrap(repository.ErrSessionDeleteFailed, err)
	}

	return nil
}
