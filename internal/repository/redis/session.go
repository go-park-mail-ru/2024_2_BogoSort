package redis

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

const userSessionPlaceholder = "user_sessions:"

type SessionDB struct {
	rdb              *redis.Client
	sessionAliveTime int
	ctx              context.Context
}

func newSessionRepository(rdb *redis.Client, sessionAliveTime int) *SessionDB {
	return &SessionDB{
		rdb:              rdb,
		sessionAliveTime: sessionAliveTime,
		ctx:              context.Background(),
	}
}
func (sdb *SessionDB) CreateSession(userID string) (string, error) {
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

	err = sdb.rdb.SAdd(sdb.ctx, userSessionPlaceholder+userID, sessionID).Err()
	if err != nil {
		return "", entity.RedisWrap(repository.ErrSessionCreationFailed, err)
	}

	return sessionID, nil
}

func (sdb *SessionDB) GetSession(sessionID string) (string, error) {
	userID, err := sdb.rdb.Get(sdb.ctx, sessionID).Result()
	if errors.Is(err, redis.Nil) {
		return "", entity.RedisWrap(repository.ErrSessionNotFound, err)
	}
	if err != nil {
		return "", entity.RedisWrap(repository.ErrSessionCheckFailed, err)
	}

	return userID, nil
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
