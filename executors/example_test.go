package executors

import (
	. "go.riyazali.net/httpx"
	. "go.riyazali.net/httpx/assertions"
	. "go.riyazali.net/httpx/helpers"
	"net/http"
	"net/http/cookiejar"
)

var t = (TestingT)(nil)

func Example_remoteEndpoint() {
	WithDefaultClient().MakeRequest(
		Get("https://httpbin.org/get"),
	).ExpectIt(t,
		ToHaveStatus(http.StatusOK),
	)
}

func Example_httpHandler() {
	var handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// handler logic here ...
	})

	WithHandler(handler).MakeRequest(
		Get("https://httpbin.org/get"),
	).ExpectIt(t,
		ToHaveStatus(http.StatusOK),
	)
}

func Example_customClient() {
	var jar, _ = cookiejar.New(&cookiejar.Options{})
	WithClient(
		WithNoRedirect(), WithCookieJar(jar),
	).MakeRequest(
		Get(Url("https://httpbin.org/cookies/set", WithQueryParam("a", "1"))),
	).ExpectIt(t,
		ToHaveStatus(http.StatusOK),
		HaveCookie("a"),
	)
}
