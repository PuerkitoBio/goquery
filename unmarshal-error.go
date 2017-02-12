package goquery

import (
	"fmt"
	"reflect"
)

// All "Reason" fields within CannotUnmarshalError will be constants and part of
// this list
const (
	nonPointer           = "non-pointer value"
	nilValue             = "destination argument is nil"
	documentReadError    = "error reading goquery document"
	arrayLengthMismatch  = "array length does not match document elements found"
	customUnmarshalError = "a custom Unmarshaler implementation threw an error"
	typeConversionError  = "a type conversion error occurred"
	nonStringMapKey      = "a map with non-string key type cannot be unmarshaled"
	missingValueSelector = "at least one value selector must be passed to use as map index"
)

// cannotUnmarshalError represents an error returned by the goquery Unmarshaler
// and helps consumers in programmatically diagnosing the cause of their error.
type cannotUnmarshalError struct {
	Err      error
	Val      string
	FldOrIdx interface{}

	v      reflect.Value
	reason string
}

// This type is a mid-level abstraction to help understand the error printing logic
type errChain struct {
	chain []*cannotUnmarshalError
	val   string
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
				switch err.v.Type().Kind() {
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

func (e errChain) last() *cannotUnmarshalError {
	return e.chain[len(e.chain)-1]
}

// Error gives a human-readable error message for debugging purposes.
func (e errChain) Error() string {
	last := e.last()

	// Avoid panic if we cannot get a type name for the Value
	t := "unknown: invalid value"
	if last.v.IsValid() {
		t = last.v.Type().String()
	}

	msg := "could not unmarshal "

	if e.val != "" {
		msg += fmt.Sprintf("value %q ", e.val)
	}

	msg += fmt.Sprintf(
		"into '%s%s' (type %s): %s",
		e.chain[0].v.Type(),
		e.tPath(),
		t,
		last.reason,
	)

	// If a generic error was reported elsewhere, report its message last
	if e.tail != nil {
		msg = msg + ": " + e.tail.Error()
	}

	return msg
}

// Traverse e.Err, printing hopefully helpful type info until there are no more
// chained errors.
func (e *cannotUnmarshalError) unwind() *errChain {
	str := &errChain{chain: []*cannotUnmarshalError{}}
	for {
		str.chain = append(str.chain, e)

		if e.Val != "" {
			str.val = e.Val
		}

		// Terminal error was of type *CannotUnmarshalError and had no children
		if e.Err == nil {
			return str
		}

		if e2, ok := e.Err.(*cannotUnmarshalError); ok {
			e = e2
			continue
		}

		// Child error was not a *CannotUnmarshalError; print its message
		str.tail = e.Err
		return str
	}
}

func (e *cannotUnmarshalError) Error() string {
	return e.unwind().Error()
}
