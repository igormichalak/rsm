package rsm

import (
	"testing"
	"time"
)

func TestSessionManager(t *testing.T) {
	sm := &SessionManager{
		Store:       &fakeStore{},
		Lifetime:    20 * time.Minute,
		TokenLength: 32,
	}

	t.Run("initializing a session", func(t *testing.T) {
		t.Parallel()
		token, err := sm.InitSession()
		assertNoError(t, err)

		_, _, err = sm.RetrieveSession(token)
		assertNoError(t, err)
	})

	t.Run("updating the session data", func(t *testing.T) {
		t.Parallel()
		token, err := sm.InitSession()
		assertNoError(t, err)

		values1 := make(map[string]any)
		values1["test_key"] = "test_value"
		err = sm.putSession(token, values1)
		assertNoError(t, err)

		values2, _, err := sm.RetrieveSession(token)
		assertNoError(t, err)
		assertEqual(t, values2["test_key"], "test_value")
	})

	t.Run("setting a single value", func(t *testing.T) {
		t.Parallel()
		token, err := sm.InitSession()
		assertNoError(t, err)

		err = sm.SetValue(token, "test_key", "test_value")
		assertNoError(t, err)

		value, err := sm.GetValue(token, "test_key")
		assertNoError(t, err)
		assertEqual(t, value, "test_value")
	})

	t.Run("deleting a single value", func(t *testing.T) {
		t.Parallel()
		token, err := sm.InitSession()
		assertNoError(t, err)

		err = sm.SetValue(token, "test_key", "test_value")
		assertNoError(t, err)

		err = sm.DeleteValue(token, "test_key")
		assertNoError(t, err)

		_, err = sm.GetValue(token, "test_key")
		assertError(t, err, ErrPropertyNotFound)
	})

	t.Run("renewing session", func(t *testing.T) {
		t.Parallel()
		token, err := sm.InitSession()
		assertNoError(t, err)

		err = sm.SetValue(token, "test_key", "test_value")
		assertNoError(t, err)

		values1, expiry1, err := sm.RetrieveSession(token)
		assertNoError(t, err)

		time.Sleep(100 * time.Nanosecond)

		newToken, err := sm.RenewToken(token)
		assertNoError(t, err)

		values2, expiry2, err := sm.RetrieveSession(newToken)
		assertNoError(t, err)
		assertEqual(t, values2, values1)
		if expiry2.UnixNano() <= expiry1.UnixNano() {
			t.Error("inadequate expiry time of the new session")
		}
	})

	t.Run("destroying session", func(t *testing.T) {
		t.Parallel()
		token, err := sm.InitSession()
		assertNoError(t, err)

		err = sm.DestroySession(token)
		assertNoError(t, err)

		_, _, err = sm.RetrieveSession(token)
		assertError(t, err, ErrSessionNotFound)
	})
}
