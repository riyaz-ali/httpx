package executors_test

import (
	. "go.riyazali.net/httpx/executors"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"testing"
	"time"
)

func assert(t *testing.T, cond bool, msg string, args ...interface{}) {
	t.Helper()
	if !cond {
		t.Errorf(msg, args...)
	}
}

func TestCustomClient(t *testing.T) {
	var jar, _ = cookiejar.New(&cookiejar.Options{})

	WithClient(
		WithTimeout(1*time.Second),
		WithCookieJar(jar),
		WithTransport(http.DefaultTransport.(*http.Transport)),
		WithNoRedirect(),

		// following builder is used to do the assertions
		// as we can't get access to the custom client outside of this scope
		func(client *http.Client) {
			assert(t, client.Timeout == (1*time.Second), "timeout must be set")
			assert(t, client.Jar == jar, "cookie jar must be set")
			assert(t, client.Transport == http.DefaultTransport.(*http.Transport), "transport must be overridden")
			assert(t, client.CheckRedirect != nil, "custom redirect must be set")
			assert(t, client.CheckRedirect(nil, nil) == http.ErrUseLastResponse, "redirect must return http.ErrUseLastResponse")
		},
	)
}

func TestDefaultClient(t *testing.T) {
	assert(t,
		reflect.ValueOf(WithDefaultClient()).Pointer() == reflect.ValueOf(http.DefaultClient.Do).Pointer(),
		"must use default http client")
}

func TestWithHandler(t *testing.T) {
	var called bool
	var handler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		called = true
		assert(t, writer != nil, "writer must not be nil")
	}

	_, _ = WithHandlerFn(handler)(&http.Request{})
	assert(t, called, "handler must be invoked")
}
