// Package httpx provides an expressive framework to test http endpoints and handlers.
package httpx

import (
	"context"
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
func (fn ExecFn) MakeRequest(method, url string, builders ...RequestBuilder) Assertable {
	var err error

	// build a new request and apply customisations
	var request *http.Request
	if request, err = http.NewRequestWithContext(context.Background(), method, url, http.NoBody); err != nil {
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
		for _, fn := range assertions {
			if err = fn(response); err != nil {
				t.Errorf("httpx: assertion: %v", err)
			}
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

// RequestBuilder defines a function that customises the request before it's sent out.
type RequestBuilder func(*http.Request) error

// Assertion defines a function that performs some sort of assertion on the response
// to make sure that request was executed as expected.
type Assertion func(*http.Response) error

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
