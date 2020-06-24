package httpx

import (
	. "go.riyazali.net/httpx/assertions"
	. "go.riyazali.net/httpx/executors"
	. "go.riyazali.net/httpx/helpers"
	"net/http"
)

func Example_remoteEndpoint() {
	WithDefaultClient().MakeRequest(
		Get("https://httpbin.org/get"),
	).ExpectIt(t,
		ToHaveStatus(http.StatusOK),
	)
}

func Example_httpHandler() {
	WithHandler(handler).MakeRequest(
		Get("https://httpbin.org/get"),
	).ExpectIt(t,
		ToHaveStatus(http.StatusOK),
	)
}

func Example_customClient() {
	WithClient(
		WithNoRedirect(),
		WithCookieJar(jar),
	).MakeRequest(
		Get(Url("https://httpbin.org/cookies/set", WithQueryParam("a", "1"))),
	).ExpectIt(t,
		ToHaveStatus(http.StatusOK),
		HaveCookie("a"),
	)
}

func ExampleAssertion_customAssertion() {
	WithClient().MakeRequest(
		Get("/versions"),
	).ExpectIt(t,
		func(r *http.Response) error {
			return AssertThat(r.StatusCode == http.StatusOK, "status not ok")
		},
	)
}

func ExampleRequestBuilder_customRequestBuilder() {
	WithClient().MakeRequest(
		Get("/versions"),
		func(request *http.Request) error {
			request.AddCookie(&http.Cookie{Name: "a", Value: "1"})
			return nil
		},
	).ExpectIt(t, ToHaveStatus(http.StatusOK))
}
