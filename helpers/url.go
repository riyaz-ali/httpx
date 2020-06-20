package helpers

import (
	"net/url"
	"strings"
)

// Url is a handy builder method to build a new url.
func Url(base string, opts ...func(*url.URL)) string {
	var u *url.URL
	if u, _ = url.Parse(base); u == nil {
		return ""
	}
	for _, opt := range opts {
		opt(u)
	}
	return u.String()
}

func WithPath(part string, parts ...string) func(*url.URL) {
	return func(u *url.URL) {
		parts = append([]string{u.Path, part}, parts...)
		u.Path = strings.Join(parts, "/")
	}
}

func WithQueryParam(key, value string) func(*url.URL) {
	return func(u *url.URL) {
		q := u.Query()
		q.Add(key, value)
		u.RawQuery = q.Encode()
	}
}

func WithUsernamePassword(username, password string) func(*url.URL) {
	return func(u *url.URL) {
		u.User = url.UserPassword(username, password)
	}
}
