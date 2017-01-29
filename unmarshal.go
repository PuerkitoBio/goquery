package goquery

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	NonPointer        = "non-pointer value"
	NilValue          = "destination argument is nil"
	UnimplementedType = "value type unimplemented"
)

type CannotUnmarshalError struct {
	V      reflect.Value
	Reason string
}

func (e *CannotUnmarshalError) Error() string {
	return fmt.Sprintf("Decode(illegal value %s) cannot proceed: %s", e.V.Type().String(), e.Reason)
}

// Unmarshaler allows for custom implementations of unmarshaling logic
type Unmarshaler interface {
	Unmarshal(*Selection) error
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
		return u.Unmarshal(d.Selection)
	}

	return unmarshalByType(d.Selection, v)
}

func unmarshalByType(s *Selection, v reflect.Value) error {
	u, v := indirect(v)

	if u != nil {
		return u.Unmarshal(s)
	}

	t := v.Type()
	switch t.Kind() {
	case reflect.Struct:
		return unmarshalStruct(s, v)
	case reflect.Slice:
		return unmarshalSlice(s, v)
	// case reflect.Array:
	// 	return unmarshalArray(s, v)
	case reflect.String:
		v.Set(reflect.ValueOf(s.Text()))
		return nil
	case reflect.Int:
		i, err := strconv.Atoi(s.Text())
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(i))
		return nil
	}

	return &CannotUnmarshalError{V: v, Reason: UnimplementedType}
}

func unmarshalStruct(root *Selection, v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		tag := fld.Tag.Get("goquery")
		sel := root.Find(tag)

		err := unmarshalByType(sel, v.Field(i))
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
