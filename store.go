package rsm

import (
	"time"
)

type Store interface {
	Retrieve(token string) (data []byte, err error)
	Insert(token string, data []byte, expiry time.Time) error
	Delete(token string) error
}
