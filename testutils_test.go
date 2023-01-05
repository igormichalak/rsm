package rsm

import (
	"errors"
	"reflect"
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

type FakeStore struct {
	records []*struct {
		token  string
		data   []byte
		expiry time.Time
	}
}

func (fs *FakeStore) Retrieve(token string) (data []byte, err error) {
	for _, rec := range fs.records {
		if rec.token == token {
			return rec.data, nil
		}
	}
	return nil, ErrSessionNotFound
}

func (fs *FakeStore) Insert(token string, data []byte, expiry time.Time) error {
	fs.records = append(fs.records, &struct {
		token  string
		data   []byte
		expiry time.Time
	}{token, data, expiry})
	return nil
}

func (fs *FakeStore) Delete(token string) error {
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
	fakeStore := &FakeStore{}
	err := fakeStore.Insert("token_1", []byte("token_1_data"), time.Now().UTC())
	assertNoError(t, err)

	data, err := fakeStore.Retrieve("token_1")
	assertNoError(t, err)
	assertEqual(t, data, []byte("token_1_data"))

	err = fakeStore.Delete("token_1")
	assertNoError(t, err)

	_, err = fakeStore.Retrieve("token_1")
	assertError(t, err, ErrSessionNotFound)
}
