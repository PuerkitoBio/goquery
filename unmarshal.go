package goquery

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"golang.org/x/net/html"
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

// Traverse e.Err, printing hopefully helpful type info until there are no more
// chained errors.
func (e *CannotUnmarshalError) unwindReason() string {
	if e == nil {
		return ""
	}

	str := ""
	ok := true
	for ok {
		// Avoid panic if we cannot get a type name for the Value
		t := "unknown type: invalid value"
		if e.V.IsValid() {
			t = e.V.Type().String()
		}

		str += fmt.Sprintf(": (%s) %s", t, e.Reason)

		// Terminal error was of type *CannotUnmarshalError and had no children
		if e.Err == nil {
			return str
		}

		e, ok = e.Err.(*CannotUnmarshalError)
	}

	// Child error was not a *CannotUnmarshalError; print its message
	return fmt.Sprintf("%s: %s", str, e.Error())
}

func (e *CannotUnmarshalError) Error() string {
	return "an error occurred while unmarshaling" + e.unwindReason()
}

// Unmarshaler allows for custom implementations of unmarshaling logic
type Unmarshaler interface {
	UnmarshalHTML([]*html.Node) error
}

// Unmarshal takes a byte slice and a destination pointer to any interface{},
// and unmarshals the document into the destination based on the `goquery`
// struct tags.
func Unmarshal(bs []byte, v interface{}) error {
	d, err := NewDocumentFromReader(bytes.NewReader(bs))

	if err != nil {
		return &CannotUnmarshalError{Err: err, Reason: DocumentReadError}
	}

	return UnmarshalSelection(d.Selection, v)
}

func wrapUnmErr(err error, v reflect.Value) error {
	if err == nil {
		return nil
	}

	return &CannotUnmarshalError{
		V:      v,
		Reason: CustomUnmarshalError,
		Err:    err,
	}
}

// UnmarshalSelection will unmarshal a goquery.Selection into an interface
// appropriately annoated with goquery tags.
func UnmarshalSelection(s *Selection, iface interface{}) error {
	v := reflect.ValueOf(iface)

	// Must come before v.IsNil() else IsNil panics on NonPointer value
	if v.Kind() != reflect.Ptr {
		return &CannotUnmarshalError{V: v, Reason: NonPointer}
	}

	if iface == nil || v.IsNil() {
		return &CannotUnmarshalError{V: v, Reason: NilValue}
	}

	u, v := indirect(v)

	if u != nil {
		return wrapUnmErr(u.UnmarshalHTML(s.Nodes), v)
	}

	return unmarshalByType(s, v)
}

func unmarshalByType(s *Selection, v reflect.Value) error {
	u, v := indirect(v)

	if u != nil {
		return wrapUnmErr(u.UnmarshalHTML(s.Nodes), v)
	}

	t := v.Type()

	switch t.Kind() {
	case reflect.Struct:
		return unmarshalStruct(s, v)
	case reflect.Slice:
		return unmarshalSlice(s, v)
	case reflect.Array:
		return unmarshalArray(s, v)
	default:
		err := unmarshalLiteral(s, v)
		if err != nil {
			return &CannotUnmarshalError{
				V:      v,
				Reason: TypeConversionError,
				Err:    err,
			}
		}
		return nil
	}
}

func unmarshalLiteral(s *Selection, v reflect.Value) error {
	t := v.Type()

	switch t.Kind() {
	case reflect.String:
		v.SetString(s.Text())
	case reflect.Bool:
		i, err := strconv.ParseBool(s.Text())
		if err != nil {
			return err
		}
		v.SetBool(i)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s.Text(), 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(s.Text(), 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(s.Text(), 64)
		if err != nil {
			return err
		}
		v.SetFloat(i)
	}
	return nil
}

func unmarshalStruct(s *Selection, v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("goquery")

		sel := s
		if tag != "" {
			sel = sel.Find(tag)
		}

		err := unmarshalByType(sel, v.Field(i))
		if err != nil {
			return &CannotUnmarshalError{
				Reason: TypeConversionError,
				Err:    err,
				V:      v,
			}
		}
	}
	return nil
}

func unmarshalArray(s *Selection, v reflect.Value) error {
	if v.Type().Len() != len(s.Nodes) {
		return &CannotUnmarshalError{
			Reason: ArrayLengthMismatch,
			V:      v,
		}
	}

	for i := 0; i < v.Type().Len(); i++ {
		err := unmarshalByType(s.Eq(i), v.Index(i))
		if err != nil {
			return &CannotUnmarshalError{
				Reason: TypeConversionError,
				Err:    err,
				V:      v,
			}
		}
	}

	return nil
}

func unmarshalSlice(s *Selection, v reflect.Value) error {

	slice := v

	for i := 0; i < s.Length(); i++ {
		eleT := v.Type().Elem()
		newV := reflect.New(eleT)

		err := unmarshalByType(s.Eq(i), newV)

		if err != nil {
			return &CannotUnmarshalError{
				Reason: TypeConversionError,
				Err:    err,
				V:      v,
			}
		}

		if eleT.Kind() != reflect.Ptr {
			newV = newV.Elem()
		}

		v = reflect.Append(v, newV)
	}

	slice.Set(v)
	return nil
}

// Stolen mostly from pkg/encoding/json/decode.go and removed some cases
// (handling `null`) that goquery doesn't need to handle.
func indirect(v reflect.Value) (Unmarshaler, reflect.Value) {
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(Unmarshaler); ok {
				return u, reflect.Value{}
			}
		}
		v = v.Elem()
	}
	return nil, v
}
