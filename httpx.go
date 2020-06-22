// Package httpx provides an expressive framework to test http endpoints and handlers.
package httpx // import "go.riyazali.net/httpx"

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
)

// ExecFn defines a function that can take an http.Request and return an http.Response (and optionally, an error).
//
// This is the core type defined by this package and instances of this type does the
// actual heavy lifting work of making the request, receiving responses and more.
// How the actual execution is done is left to the implementation. Some may make actual
// http calls to remote endpoint whereas others would call in-memory http.Handler.
//
// The core package provides two implementations that works with net/http package.
type ExecFn func(*http.Request) (*http.Response, error)

// MakeRequest is the primary entry point into the framework.
//
// This method builds a request object, apply the given customisations / builders to it and
// then pass it to the ExecFn for execution returning an Assertable which you can then use to perform
// assertions on the response etc.
//
// The core library provides certain general purpose builders. See RequestBuilder and it's implementations
// for more details on builders and how you can create a custom builder.
func (fn ExecFn) MakeRequest(factory RequestFactory, builders ...RequestBuilder) Assertable {
	var err error

	// build a new request and apply customisations
	var request *http.Request
	if request, err = factory(); err != nil {
		return fail("httpx: failed to create request: %v", err)
	}

	for _, fn := range builders {
		if err = fn(request); err != nil {
			return fail("httpx: builder: %v", err)
		}
	}

	// execute the request
	var response *http.Response
	if response, err = fn(request); err != nil {
		return fail("httpx: failed to execute request: %v", err)
	}

	// return an Assertable to run assertions on response
	return func(t TestingT, assertions ...Assertion) {
		t.Helper()
		defer response.Body.Close() // make sure to close the original response body always

		// prepare a bytes.Buffer that'd allow us to seek to start after every assertion
		// so that we can have multiple assertions that could read response's body
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(response.Body); err != nil {
			t.Errorf("httpx: failed to read body into buffer: %v", err)
			t.FailNow()
		}
		var reader = bytes.NewReader(buf.Bytes())
		response.Body = ioutil.NopCloser(reader)

		for _, fn := range assertions {
			if err = fn(response); err != nil {
				t.Errorf("httpx: assertion: %v", err)
			}
			_, _ = reader.Seek(0, io.SeekStart) // safe to ignore return values
		}
	}
}

// Assertable defines a function that can take a slice of assertions and apply it on http.Response.
//
// Although exported, user's won't be able to do much with this type. Instead they should use
// the ExpectIt(...) method to allow fluent chaining with MakeRequest(...).
type Assertable func(TestingT, ...Assertion)

// ExpectIt allows us to implement fluent chaining with MakeRequest(...).
// Use this method instead of directly invoking the Assertable to improve readability of your code.
func (a Assertable) ExpectIt(t TestingT, assertions ...Assertion) {
	t.Helper()
	a(t, assertions...)
}

// RequestFactory defines a function capable of creating http.Request instances.
// Use of this type allows us to decouple MakeRequest(...) from the actual underlying
// mechanism of building an http.Request. Implementations of this type could (say)
// create instances configured for a PaaS (like Google App Engine) and more.
//
// The core library provides a default implementation which should be sufficient for most use cases.
type RequestFactory func() (*http.Request, error)

// Using returns a RequestFactory which is a wrapper over the default http.NewRequest method
func Using(method, url string, body io.Reader) RequestFactory {
	return func() (*http.Request, error) {
		return http.NewRequestWithContext(context.Background(), method, url, body)
	}
}

// Get is a shorthand method to create a RequestFactory with http.MethodGet
func Get(url string) RequestFactory {
	return Using(http.MethodGet, url, nil)
}

// Post is a shorthand method to create a RequestFactory with http.MethodPost
func Post(url string, body io.Reader) RequestFactory {
	return Using(http.MethodPost, url, body)
}

// Put is a shorthand method to create a RequestFactory with http.MethodPut
func Put(url string, body io.Reader) RequestFactory {
	return Using(http.MethodPut, url, body)
}

// Delete is a shorthand method to create a RequestFactory with http.MethodDelete
func Delete(url string) RequestFactory {
	return Using(http.MethodDelete, url, nil)
}

// TestingT allows us to decouple our code from the actual testing.T type.
// Most end user shouldn't care about it.
type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Helper()
}

// fail returns a no-op Assertable that allows us to break out of MakeRequest(...) quicker.
func fail(format string, args ...interface{}) Assertable {
	return func(t TestingT, _ ...Assertion) {
		t.Helper()
		t.Errorf(format, args...)
		t.FailNow() // doesn't return
	}
}
