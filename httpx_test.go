package httpx_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	. "go.cubeq.co/httpx"
	"net/http"
	"reflect"
	"testing"
)

func dummyRequestBuilder(request *http.Request) error {
	return nil
}

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
			return nil, err
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
