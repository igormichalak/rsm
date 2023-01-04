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

func (sm *SessionManager) RetrieveValues(token string) (map[string]any, error) {
	if len(token) == 0 {
		return nil, ErrEmptyToken
	}
	data, err := sm.store.Retrieve(token)
	if err != nil {
		return nil, err
	}
	s := session{token: token}
	err = s.decodeData(data)
	if err != nil {
		return nil, err
	}
	return s.values, nil
}
