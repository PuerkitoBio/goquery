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
	NonStringMapKey      = "a map with non-string key type cannot be unmarshaled"
	MissingValueSelector = "at least one value selector must be passed to use as map index"
)

// CannotUnmarshalError represents an error returned by the goquery Unmarshaler
// and helps consumers in programmatically diagnosing the cause of their error.
type CannotUnmarshalError struct {
	Reason string
	Err    error

	V        reflect.Value
	FldOrIdx interface{}
}

// This type is a mid-level abstraction to help understand the error printing logic
type errChain struct {
	chain []*CannotUnmarshalError
	tail  error
}

// tPath returns the type path in the same string format one might use to access
// the nested value in go code. This should hopefully help make debugging easier.
func (e errChain) tPath() string {
	nest := ""

	for _, err := range e.chain {
		if err.FldOrIdx != nil {
			switch nesting := err.FldOrIdx.(type) {
			case string:
				switch err.V.Type().Kind() {
				case reflect.Map:
					nest += fmt.Sprintf("[%q]", nesting)
				case reflect.Struct:
					nest += fmt.Sprintf(".%s", nesting)
				}
			case int:
				nest += fmt.Sprintf("[%d]", nesting)
			default:
				nest += fmt.Sprintf("[%v]", nesting)
			}
		}
	}

	return nest
}

func (e errChain) last() *CannotUnmarshalError {
	return e.chain[len(e.chain)-1]
}

// Error gives a human-readable error message for debugging purposes.
func (e errChain) Error() string {
	last := e.last()

	// Avoid panic if we cannot get a type name for the Value
	t := "unknown: invalid value"
	if last.V.IsValid() {
		t = last.V.Type().String()
	}

	s := fmt.Sprintf(
		"could not unmarshal into %s%s (type %s): %s",
		e.chain[0].V.Type(),
		e.tPath(),
		t,
		last.Reason,
	)

	// If a generic error was reported elsewhere, report its message last
	if e.tail != nil {
		s = s + ": " + e.tail.Error()
	}

	return s
}

// Traverse e.Err, printing hopefully helpful type info until there are no more
// chained errors.
func (e *CannotUnmarshalError) unwind() *errChain {
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
	return e.unwind().Error()
}
