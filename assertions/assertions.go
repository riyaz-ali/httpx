// Package assertions provides set of handy assertions for httpx.
package assertions // import "go.riyazali.net/httpx/assertions"

import (
	"encoding/json"
	"fmt"
	"go.riyazali.net/httpx"
	. "go.riyazali.net/httpx/helpers"
	"io/ioutil"
	"net/http"
	"reflect"
)

// use with type.Implements(...) to see if type implements error
var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// ToHaveStatus returns an assertions that checks whether the request status matches the given status or not
func ToHaveStatus(status int) httpx.Assertion {
	return func(response *http.Response) error {
		return AssertThat(response.StatusCode == status,
			fmt.Sprintf("status: returned status (%d) not equal to expected status (%d)", response.StatusCode, status))
	}
}

// BodyJson returns an assertion that un-marshal the response body and invoke the given callback with the decoded value.
// The given callback must be function with following signature,
//    func cb(x X) error
// where X can be any type that json.Decoder supports.
func BodyJson(cb interface{}) httpx.Assertion {
	// extract type and value of callback
	var t = reflect.TypeOf(cb)
	var v = reflect.ValueOf(cb)

	// do some sanity checks
	if t.Kind() != reflect.Func {
		return failed(fmt.Errorf("json: given callback is not a function"))
	} else if t.NumIn() != 1 || t.IsVariadic() {
		return failed(fmt.Errorf("json: callback must only accept single argument"))
	} else if t.NumOut() != 1 || !t.Out(0).Implements(errorInterface) {
		return failed(fmt.Errorf("json: callback must only return single value of type error"))
	}

	// return an Assertable that json decodes the response and invokes callback
	return func(response *http.Response) (err error) {
		defer checkClose(response.Body, &err)

		// create a new instance of callback's 0th argument
		var arg0 = t.In(0)
		var obj = reflect.New(arg0)

		if err := json.NewDecoder(response.Body).Decode(obj.Interface()); err != nil {
			return fmt.Errorf("json: failed to decode response body: %v", err)
		}

		// invoke callback
		var ret = v.Call([]reflect.Value{obj.Elem()})
		if err, ok := ret[0].Interface().(error); ok {
			return fmt.Errorf("json: %v", err)
		}
		return nil
	}
}

// BodyBytes returns an Assertion that reads the response body and invokes the given cb with it
func BodyBytes(cb func([]byte) error) httpx.Assertion {
	return func(response *http.Response) (err error) {
		defer checkClose(response.Body, &err)

		if body, err := ioutil.ReadAll(response.Body); err != nil {
			return fmt.Errorf("body: failed to read response body: %v", err)
		} else {
			if err := cb(body); err != nil {
				return fmt.Errorf("body: %v", err)
			}
			return nil
		}
	}
}

// WithCookie returns an assertion which extracts the cookie and invokes
// the given handler function with it.
func WithCookie(name string, hn func(*http.Cookie) error) httpx.Assertion {
	return func(response *http.Response) (err error) {
		var cookie *http.Cookie
		for cookies, i := response.Cookies(), 0; i < len(cookies); i++ {
			if cookies[i].Name == name {
				cookie = cookies[i]
			}
		}
		if err := hn(cookie); err != nil {
			return fmt.Errorf("cookie: %v", err)
		}
		return nil
	}
}

// HaveCookie returns an assertion that just checks whether a cookie with the
// given name was set in the response or not.
func HaveCookie(name string) httpx.Assertion {
	return WithCookie(name, func(c *http.Cookie) error {
		return AssertThat(c != nil, fmt.Sprintf("cookie with name '%s' not set", name))
	})
}

// WithHeader returns an assertion which extracts the header value and invokes
// the given handler with all the header values found in the response.
func WithHeader(name string, hn func(string) error) httpx.Assertion {
	return func(response *http.Response) error {
		var header = response.Header.Get(name)
		if err := hn(header); err != nil {
			return fmt.Errorf("header: %v", err)
		}
		return nil
	}
}

// HaveHeader returns an assertion that just checks whether a header with the
// given name was set in the response or not.
func HaveHeader(name string) httpx.Assertion {
	return WithHeader(name, func(header string) error {
		return AssertThat(len(header) > 0, fmt.Sprintf("header with name '%s' not found", name))
	})
}
