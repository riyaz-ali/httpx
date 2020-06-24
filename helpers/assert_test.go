package helpers

import (
	"errors"
	"fmt"
	"testing"
)

func assert(t *testing.T, cond bool, msg string, args ...interface{}) {
	t.Helper()
	if !cond {
		t.Errorf(msg, args...)
	}
}

func TestAssertThat(t *testing.T) {
	assert(t, AssertThat(true, "test: %d %d %d", 1, 2, 3) == nil, "must not return error if condition evaluates to true")
	assert(t, AssertThat(false, "test: %d %d %d", 1, 2, 3) != nil, "must return error if condition evaluates to false")
	assert(t,
		AssertThat(false, "test: %d %d %d", 1, 2, 3).Error() == fmt.Errorf("test: %d %d %d", 1, 2, 3).Error(),
		"must format error using format string and variadic arguments")
}

func TestMultiple(t *testing.T) {
	assert(t, Multiple(nil, nil) == nil, "must return nil if all arguments are nil")
	assert(t, Multiple(errors.New("test")).Error() == "test\n", "must return error with single error message")
	assert(t,
		Multiple(errors.New("test1"), errors.New("test2")).Error() == "multiple errors: \n- test1\n- test2\n",
		"must return error with multiple error messages")
}
