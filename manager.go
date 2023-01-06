package rsm

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type SessionManager struct {
	Store       Store
	Lifetime    time.Duration
	TokenLength uint
}

func New(pool *redis.Pool) *SessionManager {
	return &SessionManager{
		Store:       &redisStore{pool},
		Lifetime:    20 * time.Minute,
		TokenLength: 32,
	}
}

func (sm *SessionManager) putSession(token string, values map[string]any) error {
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
	if err = sm.Store.Insert(s.token, data, s.expiry); err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) InitSession() (string, error) {
	token, err := generateRandomToken(sm.TokenLength)
	if err != nil {
		return "", err
	}
	values := make(map[string]any)
	if err = sm.putSession(token, values); err != nil {
		return "", err
	}
	return token, nil
}

func (sm *SessionManager) RetrieveSession(token string) (map[string]any, time.Time, error) {
	if len(token) == 0 {
		return nil, time.Time{}, ErrEmptyToken
	}
	data, err := sm.Store.Retrieve(token)
	if err != nil {
		return nil, time.Time{}, err
	}
	s := &session{token: token}
	if err = s.decodeData(data); err != nil {
		return nil, time.Time{}, err
	}
	return s.values, s.expiry, nil
}

func (sm *SessionManager) RenewToken(token string) (newToken string, err error) {
	values, _, err := sm.RetrieveSession(token)
	if err != nil {
		return "", err
	}
	newToken, err = generateRandomToken(sm.TokenLength)
	if err != nil {
		return "", err
	}
	if err = sm.putSession(newToken, values); err != nil {
		return "", err
	}
	return newToken, nil
}

func (sm *SessionManager) DestroySession(token string) error {
	if len(token) == 0 {
		return ErrEmptyToken
	}
	if err := sm.Store.Delete(token); err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) GetValue(token, key string) (any, error) {
	if len(key) == 0 {
		return nil, ErrEmptyValueKey
	}
	values, _, err := sm.RetrieveSession(token)
	if err != nil {
		return nil, err
	}
	value, ok := values[key]
	if !ok {
		return nil, ErrPropertyNotFound
	}
	return value, nil
}

func (sm *SessionManager) SetValue(token, key string, v any) error {
	if len(key) == 0 {
		return ErrEmptyValueKey
	}
	values, _, err := sm.RetrieveSession(token)
	if err != nil {
		return err
	}
	values[key] = v
	if err = sm.putSession(token, values); err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) DeleteValue(token, key string) error {
	if len(key) == 0 {
		return ErrEmptyValueKey
	}
	values, _, err := sm.RetrieveSession(token)
	if err != nil {
		return err
	}
	delete(values, key)
	if err = sm.putSession(token, values); err != nil {
		return err
	}
	return nil
}
