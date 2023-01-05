package rsm

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

func assertEqual[T any](t testing.TB, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if errors.Is(got, want) {
		return
	} else if got != nil {
		t.Fatalf("got a different error than expected:\n%s", got)
	} else {
		t.Fatal("didn't get the error but expected one")
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but didn't want one:\n%s", got)
	}
}

type fakeStore struct {
	records []*struct {
		token  string
		data   []byte
		expiry time.Time
	}
	mu sync.Mutex
}

func (fs *fakeStore) Retrieve(token string) (data []byte, err error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	for _, rec := range fs.records {
		if rec.token == token {
			return rec.data, nil
		}
	}
	return nil, ErrSessionNotFound
}

func (fs *fakeStore) Insert(token string, data []byte, expiry time.Time) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	for _, rec := range fs.records {
		if rec.token == token {
			rec.data = data
			rec.expiry = expiry
			return nil
		}
	}
	fs.records = append(fs.records, &struct {
		token  string
		data   []byte
		expiry time.Time
	}{token, data, expiry})
	return nil
}

func (fs *fakeStore) Delete(token string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	for i, rec := range fs.records {
		if rec.token == token {
			lastIndex := len(fs.records) - 1
			fs.records[i] = fs.records[lastIndex]
			fs.records = fs.records[:lastIndex]
		}
	}
	return nil
}

func TestFakeStore(t *testing.T) {
	fs := &fakeStore{}
	err := fs.Insert("token_1", []byte("token_1_data"), time.Now().UTC())
	assertNoError(t, err)

	data, err := fs.Retrieve("token_1")
	assertNoError(t, err)
	assertEqual(t, data, []byte("token_1_data"))

	err = fs.Delete("token_1")
	assertNoError(t, err)

	_, err = fs.Retrieve("token_1")
	assertError(t, err, ErrSessionNotFound)
}
