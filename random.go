package rsm

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func generateRandomToken(length int) (string, error) {
	buf := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
