package httpx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

// Assertion defines a function that performs some sort of assertion on the response
// to make sure that request was executed as expected.
type Assertion func(*http.Response) error

// use with type.Implements(...) to see if type implements error
var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// ToHaveStatus returns an assertions that checks whether the request status matches the given status or not
func ToHaveStatus(status int) Assertion {
	return func(response *http.Response) error {
		if response.StatusCode != status {
			return fmt.Errorf("returned status (%d) not equal to expected status (%d)", response.StatusCode, status)
		}
		return nil
	}
}

// BodyJson returns an assertion that un-marshal the response body and invoke the given callback with the decoded value.
// The given callback must be function with following signature,
//    func cb(x X) error
// where X can be any type that json.Decoder supports.
func BodyJson(cb interface{}) Assertion {
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

// BodyBytes returns an Assertion that reads the response body
func BodyBytes(cb func([]byte) error) Assertion {
	return func(response *http.Response) (err error) {
		defer checkClose(response.Body, &err)

		if body, err := ioutil.ReadAll(response.Body); err != nil {
			return err
		} else {
			return cb(body)
		}
	}
}

// failed returns a noop assertion that always returns an error
func failed(err error) Assertion {
	return func(*http.Response) error {
		return err
	}
}
