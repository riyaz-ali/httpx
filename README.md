# httpx

[![Go v1.13](https://img.shields.io/badge/v1.13-blue.svg?labelColor=a8bfc0&color=5692c7&logoColor=fff&style=for-the-badge&logo=Go)](https://golang.org/doc/go1.13)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?labelColor=a8bfc0&color=5692c7&logoColor=fff&style=for-the-badge)](https://pkg.go.dev/go.riyazali.net/httpx)
![No added preservatives](https://img.shields.io/badge/No-Added_Preservatives-blue.svg?labelColor=a8bfc0&color=5692c7&logoColor=fff&style=for-the-badge)

Simple :bowtie:, elegant :snowflake: and concise :dart: tests for those HTTP endpoints in Go!

## Overview

**`httpx`** provides a minimal, yet powerful, function-driven framework to write simple and concise tests for HTTP, that reads like poem :notes:

The primary motivation is to enable developers to write self-describing and concise test cases, that also serves as documentation.

Using **`httpx`**, you can write your test cases like,

```go
import (
  . "go.riyazali.net/httpx"
  "testing"
)

func TestHttpbinGet(t *testing.T) {
  WithDefaultClient().
    MakeRequest(
      Get(Url("https://httpbin.org/get", WithQueryParam("a", "1"))),
      WithHeader("Accept", "application/json"),
    ).
    ExpectIt(t,
      ToHaveStatus(http.StatusOK),
    )
}
```

**`httpx`** is quite versatile and allows you to use the same set of methods whether you are testing a remote endpoint (using `WithClient(...)`) or a local handler (using `WithHandler(...)`) :tada:

## Usage

To use it, first install it using 

```shell
go get -u go.riyazali.net/httpx
```

**`httpx`** defines a set of core types all of which plugs into one another to provide a simple and cohesive framework.

The primary entry point into the framework is the [`ExecFn`](https://pkg.go.dev/go.riyazali.net/httpx?tab=doc#ExecFn) which defines any function capable of executing an `http` transaction (take an `http.Request` and return an `http.Response`). The core library provides two `Execfn` functions, `WithClient(..)` and `WithHandler(..)` (and some variations of those two).

After creating the `ExecFn`, you would normally want to call it's [`MakeRequest`](https://pkg.go.dev/go.riyazali.net/httpx?tab=doc#ExecFn.MakeRequest) method, which would allow you to craft an `http.Request` and invoke the corresponding `ExecFn` on it. See [`RequestFactory`](https://pkg.go.dev/go.riyazali.net/httpx?tab=doc#RequestFactory) and [`RequestBuilder`](https://pkg.go.dev/go.riyazali.net/httpx?tab=doc#RequestBuilder) for details on how to craft a custom request, the set of built-in `factories` and `builders` the library provides and how you can write your own!

Once you've crafted a request, `MakeRequest` would return an [`Assertable`](https://pkg.go.dev/go.riyazali.net/httpx?tab=doc#Assertable) which you could then use to perform _assertions_ on the received `http.Response`. See [`Assertion`](https://pkg.go.dev/go.riyazali.net/httpx?tab=doc#Assertion) for more details.

And that's pretty much it! that's the core of the library! roughly `~120` lines of code (give or take :wink:)

### Example

This example demonstrates how you would use **`httpx`** to make a `GET` request to [`httpbin.org`](https://httbin.org) and carry out assertions on the received `http.Response`.

```golang
import (
  . "go.cubeq.co/httpx"
  . "go.cubeq.co/httpx/helpers"
  "net/http"
  "testing"
)

func TestRemote(t *testing.T) {
  type X struct {
    Args map[string]string `json:"args"`
  }

  WithDefaultClient().
    MakeRequest(
      Get(Url("https://httpbin.org/get", WithQueryParam("a", "1"))),
      WithHeader("Accept", "application/json"),
    ).
    ExpectIt(t,
      ToHaveStatus(http.StatusOK),
      BodyJson(func(x X) error {
        return Multiple(
          AssertThat(x.Args["a"] != "", "'a' is empty"),
          AssertThat(x.Args["b"] == "", "'b' is present"),
        )
      }),
    )
}
```

For more examples and samples, make sure to checkout the [`godoc` here](https://pkg.go.dev/go.riyazali.net/httpx)
