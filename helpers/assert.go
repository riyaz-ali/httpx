package helpers

import (
	"errors"
	"fmt"
	"strings"
)

// AssertThat asserts that the given condition evaluated to true and returns an error if not.
// This method allows us to write more compact assertions, such as,
//    BodyJson(func(x X) error {
//        return AssertThat(x.Y != nil, "Y is nil")
//    })
// instead of writing (more cumbersome one),
//    BodyJson(func(x X) error {
//        if x.Y == nil {
//            return errors.New("Y is nil")
//        }
//        return nil
//    })
func AssertThat(cond bool, msg string, args ...interface{}) error {
	if !cond {
		return fmt.Errorf(msg, args...)
	}
	return nil
}

// Multiple is a handy way to combine multiple assertions into one and write more compact and elegant test cases.
// It takes in a slice of errors and returns an error which combines all the error messages into one, ignoring nil errors.
//
//    BodyJson(func(x X) error {
//        return Multiple(
//            AssertThat(x.Y != nil, "Y is nil"),
//            AssertThat(x.Z != nil, "Z is nil"),
//        )
//    })
func Multiple(errs ...error) error {
	var buf strings.Builder
	var n = 0
	for _, e := range errs {
		if e == nil {
			continue
		}
		n++
		buf.WriteByte('-')
		buf.WriteByte(' ')
		buf.WriteString(e.Error())
		buf.WriteByte('\n')
	}
	if n > 1 {
		return errors.New(fmt.Sprintf("multiple errors: \n%s", buf.String()))
	} else if n == 1 {
		return errors.New(buf.String()[2:])
	} else {
		return nil
	}
}
