package helpers

import "testing"

func TestUrl(t *testing.T) {
	var url = Url("https://httpbin.org",
		WithPath("cookies", "set"),
		WithQueryParam("q", "search term"),
		WithUsernamePassword("user", "password"))

	assert(t, url == "https://user:password@httpbin.org/cookies/set?q=search+term", "must return proper serialized form or Uri")
	assert(t, Url("://httpbin.org") == "", "must return empty url if given base is malformed")
}
