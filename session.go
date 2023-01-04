package rsm

import (
	"bytes"
	"encoding/gob"
	"time"
)

type session struct {
	token  string
	values map[string]any
	expiry time.Time
}

func (s *session) encodeData() ([]byte, error) {
	x := &struct {
		values map[string]any
		expiry time.Time
	}{s.values, s.expiry}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(x); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *session) decodeData(b []byte) error {
	x := &struct {
		values map[string]any
		expiry time.Time
	}{}

	r := bytes.NewReader(b)
	if err := gob.NewDecoder(r).Decode(x); err != nil {
		return err
	}
	s.values = x.values
	s.expiry = x.expiry
	return nil
}
