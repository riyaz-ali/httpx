package httpx

import (
	"net/http"
	"net/http/httptest"
	"time"
)

// WithDefaultClient wraps the http.DefaultClient into an ExecFn and returns it
func WithDefaultClient() ExecFn {
	return http.DefaultClient.Do
}

// WithClient returns an ExecFn that wraps an http.Client.
// Use opts to customise the http.Client.
func WithClient(opts ...func(*http.Client)) ExecFn {
	var client = &http.Client{}
	for _, fn := range opts {
		fn(client)
	}
	return client.Do
}

// WithTimeout configures a timeout on the given http.Client
func WithTimeout(d time.Duration) func(*http.Client) {
	return func(c *http.Client) {
		c.Timeout = d
	}
}

// WithCookies set the given cookie jar on the http.Client.
// The jar is consulted for cookies on requests made by the client and store cookies
// from responses. If you don't need that, you can also set cookies on individual request.
func WithCookieJar(jar http.CookieJar) func(*http.Client) {
	return func(client *http.Client) {
		client.Jar = jar
	}
}

// WithTransport sets the given transport to use with the client.
// Use this to set custom transport that, for example, does TLS client authentication and more.
func WithTransport(t *http.Transport) func(*http.Client) {
	return func(c *http.Client) {
		c.Transport = t
	}
}

// WithHandler wraps the given http.Handler and returns an ExecFn that invokes
// the handler on request and return the response. This ExecFn doesn't need to make network round-trip
// and can be used to implement unit tests for http endpoints in your application.
func WithHandler(handler http.Handler) ExecFn {
	return func(request *http.Request) (*http.Response, error) {
		var recorder = httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		return recorder.Result(), nil
	}
}

// WithHandlerFn wraps the given http.HandlerFunc and returns an ExecFn.
// See WithHandler(...) for more details.
func WithHandlerFn(fn http.HandlerFunc) ExecFn {
	return WithHandler(fn)
}
