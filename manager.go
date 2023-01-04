package rsm

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type SessionManager struct {
	store    Store
	Lifetime time.Duration
}

func New(pool *redis.Pool) *SessionManager {
	return &SessionManager{
		store: &redisStore{pool},
	}
}

func (sm *SessionManager) InitSession() (string, error) {
	token, err := generateRandomToken(32)
	if err != nil {
		return "", err
	}
	return token, nil
}
