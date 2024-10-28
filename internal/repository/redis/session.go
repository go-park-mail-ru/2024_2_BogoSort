package redis

import "github.com/redis/go-redis/v9"

type SessionDB struct {
	rdb *redis.Client
}

func newSessionDB(rdb *redis.Client) *SessionDB {
	return &SessionDB{rdb: rdb}
}
