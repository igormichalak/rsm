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

func (sm *SessionManager) DestroySession(token string) error {
	if len(token) == 0 {
		return ErrEmptyToken
	}
	if err := sm.store.Delete(token); err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) RetrieveData(token string) (map[string]any, time.Time, error) {
	if len(token) == 0 {
		return nil, time.Time{}, ErrEmptyToken
	}
	data, err := sm.store.Retrieve(token)
	if err != nil {
		return nil, time.Time{}, err
	}
	s := &session{token: token}
	if err = s.decodeData(data); err != nil {
		return nil, time.Time{}, err
	}
	return s.values, s.expiry, nil
}

func (sm *SessionManager) setValues(token string, values map[string]any) error {
	if len(token) == 0 {
		return ErrEmptyToken
	}
	s := &session{
		token:  token,
		values: values,
		expiry: time.Now().Add(sm.Lifetime).UTC(),
	}
	data, err := s.encodeData()
	if err != nil {
		return err
	}
	if err = sm.store.Insert(s.token, data, s.expiry); err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) RenewSession(token string) error {
	values, _, err := sm.RetrieveData(token)
	if err != nil {
		return err
	}
	if err = sm.setValues(token, values); err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) GetValue(token, key string) (any, error) {
	if len(key) == 0 {
		return nil, ErrEmptyValueKey
	}
	values, _, err := sm.RetrieveData(token)
	if err != nil {
		return nil, err
	}
	return values[key], nil
}
