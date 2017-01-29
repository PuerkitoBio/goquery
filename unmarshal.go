package goquery

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

// All "Reason" fields within CannotUnmarshalError will be constants and part of
// this list
const (
	NonPointer          = "non-pointer value"
	NilValue            = "destination argument is nil"
	DocumentReadError   = "error reading goquery document"
	ArrayLengthMismatch = "array length does not match document elements found"
)

type CannotUnmarshalError struct {
	V      reflect.Value
	Reason string
	Err    error
}

func (e *CannotUnmarshalError) Error() string {
	if e.Error == nil {
		return fmt.Sprintf("Decode(illegal value type %q); cannot proceed because %q", e.V.Type().String(), e.Reason)
	}
	return fmt.Sprintf("An error occurred while decoding value type %q: %s", e.V.Type().String(), e.Err)
}

// Unmarshaler allows for custom implementations of unmarshaling logic
type Unmarshaler interface {
	UnmarshalSelection(*Selection) error
}

func Unmarshal(bs []byte, v interface{}) error {
	d, err := NewDocumentFromReader(bytes.NewReader(bs))

	if err != nil {
		return &CannotUnmarshalError{Err: err, Reason: DocumentReadError}
	}

	return UnmarshalDocument(d, v)
}

// UnmarshalDoc will unmarshal a goquery.Document into an interface
// appropriately annoated with goquery tags.
func UnmarshalDocument(d *Document, iface interface{}) error {
	v := reflect.ValueOf(iface)

	if iface == nil {
		return &CannotUnmarshalError{V: v, Reason: NilValue}
	}

	if v.Kind() != reflect.Ptr {
		return &CannotUnmarshalError{V: v, Reason: NonPointer}
	}

	u, v := indirect(v)

	if u != nil {
		return u.UnmarshalSelection(d.Selection)
	}

	return unmarshalByType(d.Selection, v)
}

func unmarshalByType(s *Selection, v reflect.Value) error {
	u, v := indirect(v)

	if u != nil {
		return u.UnmarshalSelection(s)
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
		return trySetLiteral(s, v)
	}
	return nil
}

func trySetLiteral(s *Selection, v reflect.Value) error {
	t := v.Type()

	switch t.Kind() {
	case reflect.String:
		v.Set(reflect.ValueOf(s.Text()))
	case reflect.Bool:
		v.Set(reflect.ValueOf(s.Text() == "true"))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := reflect.New(v.Type())
		i, err := strconv.Atoi(s.Text())
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(i))
	}
	return nil
}

func unmarshalStruct(root *Selection, v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("goquery")

		sel := root
		if tag != "" {
			sel = sel.Find(tag)
		}

		err := unmarshalByType(sel, v.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalArray(root *Selection, dest reflect.Value) error {
	if dest.Type().Len() != root.Length() {
		return &CannotUnmarshalError{
			Reason: ArrayLengthMismatch,
			V:      dest,
		}
	}

	for i := 0; i < dest.Type().Len(); i++ {
		err := unmarshalByType(root.Eq(i), dest.Index(i))
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalSlice(root *Selection, slice reflect.Value) error {
	v := slice

	for i := 0; i < root.Length(); i++ {
		eleT := v.Type().Elem()
		newV := reflect.New(eleT)

		err := unmarshalByType(root.Eq(i), newV)

		if err != nil {
			return err
		}

		if eleT.Kind() != reflect.Ptr {
			newV = newV.Elem()
		}

		v = reflect.Append(v, newV)
	}

	slice.Set(v)

	return nil
}

// Stolen mostly from pkg/encoding/json/decode.go and altered for goquery usage
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
