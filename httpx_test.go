package httpx_test

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	. "go.cubeq.co/httpx"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func dummyAssertion(response *http.Response) error {
	return nil
}

// TestingT implementation that logs it's method calls
type reporter map[string]int

func (r reporter) Errorf(_ string, _ ...interface{}) {
	r["Errorf"] = r["Errorf"] + 1
}
func (r reporter) FailNow() {
	r["FailNow"] = r["FailNow"] + 1
}
func (r reporter) Helper() {
	r["Helper"] = r["Helper"] + 1
}

func TestAssertable_ExpectIt(t *testing.T) {
	var a = Assertable(func(_ TestingT, assertions ...Assertion) {
		assert.True(t, len(assertions) == 1)
		assert.True(t, reflect.ValueOf(dummyAssertion).Pointer() == reflect.ValueOf(assertions[0]).Pointer())
	})

	a.ExpectIt(make(reporter), dummyAssertion)
}

func TestExecFn_MakeRequest(t *testing.T) {
	// helper function to build noop ExecFn
	var execer = func(err error) ExecFn {
		return func(*http.Request) (*http.Response, error) {
			var body = ioutil.NopCloser(bytes.NewBuffer(nil))
			return &http.Response{StatusCode: http.StatusOK, Body: body}, err
		}
	}

	t.Run("should fail if cannot build request", func(t *testing.T) {
		r := make(reporter)
		execer(nil).MakeRequest(Using("n/a", "", nil)).ExpectIt(r)
		assert.Equal(t, 1, r["Errorf"])
		assert.Equal(t, 1, r["FailNow"])
	})

	t.Run("should fail if builder returns error", func(t *testing.T) {
		r := make(reporter)
		execer(nil).MakeRequest(Get("http://example.com"), func(*http.Request) error {
			return errors.New("test")
		}).ExpectIt(r)
		assert.Equal(t, 1, r["Errorf"])
		assert.Equal(t, 1, r["FailNow"])
	})

	t.Run("should fail if could not execute request", func(t *testing.T) {
		r := make(reporter)
		execer(errors.New("test")).MakeRequest(Get("http://example.com")).ExpectIt(r)
		assert.Equal(t, 1, r["Errorf"])
		assert.Equal(t, 1, r["FailNow"])
	})

	t.Run("should fail if assertion fails", func(t *testing.T) {
		r := make(reporter)
		execer(nil).MakeRequest(Get("http://example.com")).ExpectIt(r, func(*http.Response) error {
			return errors.New("test")
		})
		assert.Equal(t, 1, r["Errorf"])
		assert.Equal(t, 0, r["FailNow"])
	})

	t.Run("post request factory", func(t *testing.T) {
		r := make(reporter)
		execer(nil).MakeRequest(Post("https://example.com", nil)).ExpectIt(r)
		assert.Equal(t, 0, r["Errorf"])
		assert.Equal(t, 0, r["FailNow"])
	})

	t.Run("put request factory", func(t *testing.T) {
		r := make(reporter)
		execer(nil).MakeRequest(Put("https://example.com", nil)).ExpectIt(r)
		assert.Equal(t, 0, r["Errorf"])
		assert.Equal(t, 0, r["FailNow"])
	})

	t.Run("delete request factory", func(t *testing.T) {
		r := make(reporter)
		execer(nil).MakeRequest(Delete("https://example.com")).ExpectIt(r)
		assert.Equal(t, 0, r["Errorf"])
		assert.Equal(t, 0, r["FailNow"])
	})
}

type errorReader struct {
	error
}

func (e errorReader) Read([]byte) (int, error) {
	return 0, e.error
}

func TestReadResponseBodyMultipleTimes(t *testing.T) {
	var execer = func(body io.Reader) ExecFn {
		return func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(body)}, nil
		}
	}

	t.Run("propagates error properly", func(t *testing.T) {
		r := make(reporter)
		execer(errorReader{errors.New("test")}).
			MakeRequest(Get("https://example.com")).ExpectIt(r)
		assert.Equal(t, 1, r["Errorf"])
		assert.Equal(t, 1, r["FailNow"])
	})

	t.Run("multiple assertions reading response should work", func(t *testing.T) {
		readBody := Assertion(func(response *http.Response) error {
			_, _ = ioutil.ReadAll(response.Body)
			return response.Body.Close()
		})

		r := make(reporter)
		execer(bytes.NewReader(nil)).
			MakeRequest(Get("https://example.com")).ExpectIt(r, readBody, readBody)
		assert.Equal(t, 0, r["Errorf"])
		assert.Equal(t, 0, r["FailNow"])
	})
}
