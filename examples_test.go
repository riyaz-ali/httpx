package httpx

import (
	"errors"
	"net/http"
)

var t = (TestingT)(nil)

var WithExecFn = func() ExecFn {
	return func(*http.Request) (*http.Response, error) { return nil, nil }
}

func ExampleAssertion_customAssertion() {
	WithExecFn().MakeRequest(
		Get("/versions"),
	).ExpectIt(t,
		func(r *http.Response) error {
			if r.StatusCode != http.StatusOK {
				return errors.New("status not ok")
			}
			return nil
		},
	)
}

func ExampleRequestBuilder_customRequestBuilder() {
	WithExecFn().MakeRequest(
		Get("/versions"),
		func(request *http.Request) error {
			request.AddCookie(&http.Cookie{Name: "a", Value: "1"})
			return nil
		},
	).ExpectIt(t)
}
