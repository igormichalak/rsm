package rsm

import (
	"time"
)

type session struct {
	token  string
	values map[string]any
	expiry time.Time
}
