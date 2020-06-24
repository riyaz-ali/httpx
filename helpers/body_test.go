package helpers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestEmptyBody(t *testing.T) {
	assert(t, EmptyBody() == http.NoBody, "must use http.NoBody magic variable")
}

func TestSerializeJson(t *testing.T) {
	var body, _ = ioutil.ReadAll(SerializeJson(map[string]string{"a": "1"}))
	assert(t, bytes.Equal(body, []byte("{\"a\":\"1\"}\n")), "must return json representation of object")
}
