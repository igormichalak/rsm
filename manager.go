package rsm

import (
	"github.com/gomodule/redigo/redis"
)

type SessionManager struct {
	Store Store
}

func New(pool *redis.Pool) *SessionManager {
	return &SessionManager{
		Store: &redisStore{pool},
	}
}
