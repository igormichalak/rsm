package rsm

import (
	"fmt"
	"time"
)

type Store interface {
	Retrieve(token string) (data []byte, err error)
	Insert(token string, data []byte, expiry time.Time) error
	Delete(token string) error
}

type storeError string

func (se storeError) Error() string {
	return fmt.Sprintf("rsm: %s", se)
}

const (
	ErrNotFound = storeError("entry not found")
)
