package rsm

import (
	"reflect"
	"testing"
)

func assertEqual[T any](t testing.TB, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but didn't want one:\n%s", got)
	}
}
