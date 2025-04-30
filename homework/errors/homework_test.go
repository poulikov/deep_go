package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

// MultiError provides an obsolete way to join a few errors; much better way is using errors.Join or "%w" verb.
//
// Deprecated: see errors.Join or fmt.Errorf with multiple %w verbs
type MultiError struct {
	err []error
}

func (e *MultiError) Error() string {
	if len(e.err) == 0 {
		return "<nil>"
	}
	msg := make([]string, 0, len(e.err))
	for _, emsg := range e.err {
		msg = append(msg, emsg.Error())
	}
	return fmt.Sprintf("%d errors occured:\n\t* %s\n", len(e.err), strings.Join(msg, "\t* "))
}

func Append(err error, errs ...error) *MultiError {
	if me, ok := err.(*MultiError); ok {
		me.err = append(me.err, errs...)
		return me
	}
	me := &MultiError{err: errs}
	return me
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
