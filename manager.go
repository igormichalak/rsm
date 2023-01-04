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
	s := &session{
		token:  token,
		values: make(map[string]any),
		expiry: time.Now().Add(sm.Lifetime).UTC(),
	}
	data, err := s.encodeData()
	if err != nil {
		return "", err
	}
	if err = sm.store.Insert(s.token, data, s.expiry); err != nil {
		return "", err
	}
	return s.token, nil
}

func (sm *SessionManager) RetrieveValues(token string) (map[string]any, error) {
	if len(token) == 0 {
		return nil, ErrEmptyToken
	}
	data, err := sm.store.Retrieve(token)
	if err != nil {
		return nil, err
	}
	s := &session{token: token}
	if err = s.decodeData(data); err != nil {
		return nil, err
	}
	return s.values, nil
}
