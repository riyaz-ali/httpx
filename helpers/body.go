package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// EmptyBody returns an io.ReadCloser with no bytes. Read always returns EOF
// and Close always returns nil. Use it when you don't have / want to send any
// body with outgoing request. You can also set the body to nil, but using EmptyBody(..)
// makes your test case more readable. Internally, this method just returns the http.NoBody magic
// value. Refer to docs in net/http to know more about it.
func EmptyBody() io.ReadCloser {
	return http.NoBody
}

// SerializeJson serializes the given object using encoding/json and returns an io.ReadCloser.
func SerializeJson(obj interface{}) io.ReadCloser {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(obj)
	return ioutil.NopCloser(&buf)
}
