package rsm

import (
	"testing"
	"time"
)

func TestCodec(t *testing.T) {
	token := "token_id"
	values := make(map[string]any)
	values["test_key"] = "test_value"
	expiry := time.Now().UTC()

	want := session{token, values, expiry}
	got := session{token: token}

	data, err := want.encodeData()
	assertNoError(t, err)

	err = got.decodeData(data)
	assertNoError(t, err)

	assertEqual(t, got, want)
}
