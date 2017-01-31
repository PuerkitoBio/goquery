package goquery

import (
	"fmt"
	"reflect"
)

// All "Reason" fields within CannotUnmarshalError will be constants and part of
// this list
const (
	NonPointer           = "non-pointer value"
	NilValue             = "destination argument is nil"
	DocumentReadError    = "error reading goquery document"
	ArrayLengthMismatch  = "array length does not match document elements found"
	CustomUnmarshalError = "a custom Unmarshaler implementation threw an error"
	TypeConversionError  = "a type conversion error occurred"
)

// CannotUnmarshalError represents an error returned by the goquery Unmarshaler
// and helps consumers in programmatically diagnosing the cause of their error.
type CannotUnmarshalError struct {
	V      reflect.Value
	Reason string
	Err    error
}

// This type is a mid-level abstraction to help understand the error printing logic
type errChain struct {
	chain []*CannotUnmarshalError
	tail  error
}

func (e errChain) Error() string {
	s := ""
	for _, err := range e.chain {
		// Avoid panic if we cannot get a type name for the Value
		t := "unknown type: invalid value"
		if err.V.IsValid() {
			t = err.V.Type().String()
		}

		s += fmt.Sprintf(": (%s) %s", t, err.Reason)
	}

	if e.tail == nil {
		return s
	}

	return s + ": " + e.tail.Error()
}

// Traverse e.Err, printing hopefully helpful type info until there are no more
// chained errors.
func (e *CannotUnmarshalError) unwindReason() *errChain {
	str := &errChain{chain: []*CannotUnmarshalError{}}
	for {
		str.chain = append(str.chain, e)

		// Terminal error was of type *CannotUnmarshalError and had no children
		if e.Err == nil {
			return str
		}

		if e2, ok := e.Err.(*CannotUnmarshalError); ok {
			e = e2
			continue
		}

		// Child error was not a *CannotUnmarshalError; print its message
		str.tail = e.Err
		return str
	}
}

func (e *CannotUnmarshalError) Error() string {
	return "an error occurred while unmarshaling" + e.unwindReason().Error()
}
