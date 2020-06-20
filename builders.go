package httpx

import (
	"fmt"
	"net/http"
)

// RequestBuilder defines a function that customises the request before it's sent out.
type RequestBuilder func(*http.Request) error

// WithHeader takes in a header name and one or more values and returns a RequestBuilder.
// The first value is added using header.Add(...) method whereas remaining values are
// added using header.Set(...). See net/http.Header more info.
func WithHeader(name, v string, values ...string) RequestBuilder {
	return func(request *http.Request) error {
		request.Header.Set(name, v)
		for _, val := range values {
			request.Header.Add(name, val)
		}
		return nil
	}
}

// WithUserAgent returns a RequestBuilder that adds a user agent to outgoing request.
func WithUserAgent(name string) RequestBuilder {
	return WithHeader("User-Agent", name)
}

// WithBasicAuth sets HTTP Basic auth on the request.
func WithBasicAuth(username, password string) RequestBuilder {
	return func(request *http.Request) error {
		request.SetBasicAuth(username, password)
		return nil
	}
}

// WithAuthorization adds an Authorization header with the given scheme and credentials
func WithAuthorization(scheme, credentials string) RequestBuilder {
	return WithHeader("Authorization", fmt.Sprintf("%s %s", scheme, credentials))
}

// WithHost changes the host value used by the request.
// By default outgoing requests use the value from url.Host for Host header. Setting this overrides
// the default behaviour and changes the Host header sent in the request.
func WithHost(host string) RequestBuilder {
	return func(request *http.Request) error {
		request.Host = host
		return nil
	}
}
