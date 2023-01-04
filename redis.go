package rsm

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type redisStore struct {
	pool *redis.Pool
}

func (rs *redisStore) Retrieve(token string) (data []byte, err error) {
	conn := rs.pool.Get()
	defer conn.Close()

	data, err = redis.Bytes(conn.Do("GET", token))
	if err == redis.ErrNil {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return data, nil
}

func (rs *redisStore) Insert(token string, data []byte, expiry time.Time) error {
	conn := rs.pool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("SET", token, data)
	if err != nil {
		return err
	}
	err = conn.Send("PEXPIREAT", token, expiry.UnixMilli())
	if err != nil {
		return err
	}
	_, err = conn.Do("EXEC")
	return err
}

func (rs *redisStore) Delete(token string) error {
	conn := rs.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", token)
	return err
}
