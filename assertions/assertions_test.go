package assertions_test

import (
	"bytes"
	"errors"
	. "go.riyazali.net/httpx/assertions"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func assert(t *testing.T, cond bool, msg string, args ...interface{}) {
	t.Helper()
	if !cond {
		t.Errorf(msg, args...)
	}
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) { return 0, errors.New("test") }

type errorCloser struct {
	io.Reader
}

func (e *errorCloser) Close() error { return errors.New("test: I was asked to error out") }

func TestHaveStatus(t *testing.T) {
	// given
	var writer = httptest.NewRecorder()
	writer.WriteHeader(http.StatusOK)

	var resp = writer.Result()

	// when
	assert(t, ToHaveStatus(http.StatusOK)(resp) == nil, "status must be ok")
	assert(t, ToHaveStatus(http.StatusNotFound)(resp) != nil, "status must not be not found")
}

func TestHaveCookie(t *testing.T) {
	// given
	var writer = httptest.NewRecorder()
	http.SetCookie(writer, &http.Cookie{Name: "a", Value: "1"})

	var resp = writer.Result()

	// when
	assert(t, HaveCookie("a")(resp) == nil, "a must be set")
	assert(t, HaveCookie("b")(resp) != nil, "b must not be set")
}

func TestHaveHeader(t *testing.T) {
	// given
	var writer = httptest.NewRecorder()
	writer.Header().Set("content-type", "application/json")

	var resp = writer.Result()

	// when
	assert(t, HaveHeader("content-type")(resp) == nil, "content-type must be set")
	assert(t, HaveHeader("x-request-id")(resp) != nil, "x-request-id must not be set")
}

func TestBodyBytes(t *testing.T) {
	t.Run("should invoke callback with correct payload", func(t *testing.T) {
		var writer = httptest.NewRecorder()
		_, _ = io.WriteString(writer, "hello world")
		var resp = writer.Result()

		var err = BodyBytes(func(body []byte) error {
			if !bytes.Equal(body, []byte("hello world")) {
				return errors.New("test failed")
			}
			return nil
		})(resp)

		assert(t, err == nil, "body content does not match")
	})

	t.Run("should return error if failed to read response body", func(t *testing.T) {
		// given
		var resp = httptest.NewRecorder().Result()
		resp.Body = ioutil.NopCloser(&errorReader{})

		// when
		var err = BodyBytes(func([]byte) error { return nil })(resp)

		assert(t, err != nil, "must return error if cannot read body")
	})

	t.Run("should return error if callback returns error", func(t *testing.T) {
		// given
		var writer = httptest.NewRecorder()
		var resp = writer.Result()

		// when
		var err = BodyBytes(func([]byte) error { return errors.New("test") })(resp)

		assert(t, err != nil, "must return error if callback returns one")
	})

	t.Run("should error if could not close response body", func(t *testing.T) {
		var writer = httptest.NewRecorder()
		var resp = writer.Result()
		resp.Body = &errorCloser{resp.Body}

		var err = BodyBytes(func(body []byte) error { return nil })(resp)

		assert(t, err != nil, "should return error if fails to close response body")
	})
}

func TestBodyJson(t *testing.T) {
	t.Run("should invoke callback with decoded value", func(t *testing.T) {
		// given
		var writer = httptest.NewRecorder()
		_, _ = io.WriteString(writer, "{\"a\": \"1\"}")
		var resp = writer.Result()

		var err = BodyJson(func(m map[string]string) error {
			if v, ok := m["a"]; !ok || v != "1" {
				return errors.New("test failed")
			}
			return nil
		})(resp)
		assert(t, err == nil, "decoded content must be correct")
	})

	t.Run("should error if callback signature is not correct", func(t *testing.T) {
		assert(t, BodyJson(1)(nil) != nil, "callback must except only functions")
		assert(t, BodyJson(func(a, b int) {}) != nil, "callback should not take more than one argument")
		assert(t, BodyJson(func(a ...int) {}) != nil, "callback should not take variadic argument")
		assert(t, BodyJson(func(a map[string]string) {}) != nil, "callback must return only error")
	})

	t.Run("should error if failed to decode", func(t *testing.T) {
		// given
		var resp = httptest.NewRecorder().Result()
		resp.Body = ioutil.NopCloser(&errorReader{})

		var err = BodyJson(func(m map[string]string) error { return nil })(resp)
		assert(t, err != nil, "must return error if fails to decode object")
	})

	t.Run("should error if callback returns error", func(t *testing.T) {
		// given
		var writer = httptest.NewRecorder()
		_, _ = io.WriteString(writer, "{\"a\": \"1\"}")
		var resp = writer.Result()

		var err = BodyJson(func(map[string]string) error { return errors.New("test") })(resp)
		assert(t, err != nil, "must return error if callback returns one")
	})

	t.Run("should error if could not close response body", func(t *testing.T) {
		var writer = httptest.NewRecorder()
		_, _ = io.WriteString(writer, "{\"a\": \"1\"}")
		var resp = writer.Result()
		resp.Body = &errorCloser{resp.Body}

		var err = BodyBytes(func(body []byte) error { return nil })(resp)

		assert(t, err != nil, "should return error if fails to close response body")
	})
}
