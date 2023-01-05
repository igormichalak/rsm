package rsm

import (
	"fmt"
)

type managerError string

func (e managerError) Error() string {
	return fmt.Sprintf("rsm: %s", string(e))
}

const (
	ErrSessionNotFound  = managerError("session not found")
	ErrPropertyNotFound = managerError("property not found")
	ErrEmptyToken       = managerError("the provided token is empty")
	ErrEmptyValueKey    = managerError("the provided value key is empty")
)
