package assertions

import (
	"go.riyazali.net/httpx"
	"io"
	"net/http"
)

// CheckClose calls Close on the given io.Closer. If the given *error points to
// nil, it will be assigned the error returned by Close. Otherwise, any error
// returned by Close will be ignored. CheckClose is usually called with defer.
//
// taken from https://git.io/Jfb2b
func checkClose(closer io.Closer, err *error) {
	if cerr := closer.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}

// failed returns a noop assertion that always returns an error
func failed(err error) httpx.Assertion {
	return func(*http.Response) error {
		return err
	}
}
