package builders

import (
	"bytes"
	"net/http"
	"testing"
)

func assert(t *testing.T, cond bool, msg string, args ...interface{}) {
	t.Helper()
	if !cond {
		t.Errorf(msg, args...)
	}
}

func require(t *testing.T, cond bool, msg string, args ...interface{}) {
	t.Helper()
	if !cond {
		t.Errorf(msg, args...)
		t.FailNow()
	}
}

func newRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/", bytes.NewReader(nil))
	return request
}

func TestWithHeader(t *testing.T) {
	t.Run("with single value", func(t *testing.T) {
		var r = newRequest()
		var err = WithHeader("content-type", "application/json")(r)
		require(t, err == nil, "builder must not return error")
		assert(t, r.Header.Get("content-type") == "application/json", "must set header with correct value on request")
	})

	t.Run("with multiple values", func(t *testing.T) {
		var r = newRequest()
		var err = WithHeader("Accept", "application/json", "application/xml")(r)
		require(t, err == nil, "builder must not return error")

		assert(t, len(r.Header["Accept"]) == 2, "must set all values for header")
	})
}

func TestWithUserAgent(t *testing.T) {
	var r = newRequest()
	var err = WithUserAgent("httpx")(r)
	require(t, err == nil, "builder must not return error")
	assert(t, r.UserAgent() == "httpx", "user agent header must be set properly")
}

func TestWithBasicAuth(t *testing.T) {
	var r = newRequest()
	var err = WithBasicAuth("user", "$ecret")(r)
	require(t, err == nil, "builder must not return error")

	u, p, ok := r.BasicAuth()
	assert(t, ok, "basic auth must be set")
	assert(t, u == "user", "must have correct username")
	assert(t, p == "$ecret", "must have correct password")
}

func TestWithAuthorization(t *testing.T) {
	var r = newRequest()
	var err = WithAuthorization("Bearer", "token")(r)
	require(t, err == nil, "builder must not return error")

	assert(t, r.Header.Get("Authorization") == "Bearer token", "Authorization header must be set")
}

func TestWithHost(t *testing.T) {
	var r = newRequest()
	var err = WithHost("httpbin.org")(r)
	require(t, err == nil, "builder must not return error")
	assert(t, r.Host == "httpbin.org", "host must be overridden")
}
