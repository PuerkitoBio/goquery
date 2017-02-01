package goquery

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// Now included in GoQuery is the ability to declaratively unmarshal your HTML
// j
//
// Map, Slice, Array may use value selector if element type is primitive

// Unmarshaler allows for custom implementations of unmarshaling logic
type Unmarshaler interface {
	UnmarshalHTML([]*html.Node) error
}

type valFunc func(*Selection) string

type goqueryTag string

func (tag goqueryTag) selector(which int) string {
	arr := strings.Split(string(tag), ",")
	if which > len(arr)-1 {
		return ""
	}
	return arr[which]
}

var (
	textVal valFunc = func(s *Selection) string {
		return strings.TrimSpace(s.Text())
	}
	htmlVal = func(s *Selection) string {
		str, _ := s.Html()
		return strings.TrimSpace(str)
	}

	vfCache = map[goqueryTag]valFunc{}
)

func attrFunc(attr string) valFunc {
	return func(s *Selection) string {
		str, _ := s.Attr(attr)
		return str
	}
}

func (tag goqueryTag) valFunc() valFunc {
	if fn := vfCache[tag]; fn != nil {
		return fn
	}

	srcArr := strings.Split(string(tag), ",")
	if len(srcArr) < 2 {
		vfCache[tag] = textVal
		return textVal
	}

	src := srcArr[1]

	var f valFunc
	switch {
	case src[0] == '[':
		// [someattr] will return value of .Attr("someattr")
		attr := src[1 : len(src)-1]
		f = attrFunc(attr)
	case src == "html":
		f = htmlVal
	default:
		f = textVal
	}

	vfCache[tag] = f
	return f
}

// popVal should allow us to handle arbitrarily nested maps as well as the
// cleanly handling the possiblity of map[literal]literal by just delegating
// back to `unmarshalByType`.
func (tag goqueryTag) popVal() goqueryTag {
	arr := strings.Split(string(tag), ",")
	if len(arr) < 2 {
		return tag
	}
	newA := []string{arr[0]}
	newA = append(newA, arr[2:]...)

	return goqueryTag(strings.Join(newA, ","))
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

	return unmarshalByType(s, v, "")
}

func unmarshalByType(s *Selection, v reflect.Value, f goqueryTag) error {
	u, v := indirect(v)

	if u != nil {
		return wrapUnmErr(u.UnmarshalHTML(s.Nodes), v)
	}

	// Handle special cases where we can just set the value directly
	switch val := v.Interface().(type) {
	case []*html.Node:
		val = append(val, s.Nodes...)
		v.Set(reflect.ValueOf(val))
		return nil
	}

	t := v.Type()

	switch t.Kind() {
	case reflect.Struct:
		return unmarshalStruct(s, v)
	case reflect.Slice:
		return unmarshalSlice(s, v, f)
	case reflect.Array:
		return unmarshalArray(s, v, f)
	case reflect.Map:
		return unmarshalMap(s, v, f)
	default:
		vf := f.valFunc()
		err := unmarshalLiteral(vf(s), v)
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

func unmarshalLiteral(s string, v reflect.Value) error {
	t := v.Type()

	switch t.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		i, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(i)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(s, 64)
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
		tag := goqueryTag(t.Field(i).Tag.Get("goquery"))

		sel := s
		if tag != "" {
			selStr := tag.selector(0)
			sel = sel.Find(selStr)
		}

		err := unmarshalByType(sel, v.Field(i), tag)
		if err != nil {
			return &CannotUnmarshalError{
				Reason:   TypeConversionError,
				Err:      err,
				V:        v,
				FldOrIdx: t.Field(i).Name,
			}
		}
	}
	return nil
}

func unmarshalArray(s *Selection, v reflect.Value, f goqueryTag) error {
	if v.Type().Len() != len(s.Nodes) {
		return &CannotUnmarshalError{
			Reason: ArrayLengthMismatch,
			V:      v,
		}
	}

	for i := 0; i < v.Type().Len(); i++ {
		err := unmarshalByType(s.Eq(i), v.Index(i), f)
		if err != nil {
			return &CannotUnmarshalError{
				Reason:   TypeConversionError,
				Err:      err,
				V:        v,
				FldOrIdx: i,
			}
		}
	}

	return nil
}

func unmarshalSlice(s *Selection, v reflect.Value, f goqueryTag) error {
	slice := v
	eleT := v.Type().Elem()

	for i := 0; i < s.Length(); i++ {
		newV := reflect.New(eleT)

		err := unmarshalByType(s.Eq(i), newV, f)

		if err != nil {
			return &CannotUnmarshalError{
				Reason:   TypeConversionError,
				Err:      err,
				V:        v,
				FldOrIdx: i,
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

func childrenUntilMatch(s *Selection, sel string) *Selection {
	orig := s
	s = s.Children()
	for s.Length() != 0 && s.Filter(sel).Length() == 0 {
		s = s.Children()
	}
	if s.Length() == 0 {
		return orig
	}
	return s.Filter(sel)
}

func unmarshalMap(s *Selection, v reflect.Value, f goqueryTag) error {
	// Make new map here because indirect for some reason doesn't help us out
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}

	if f.selector(1) == "" {
		// We need minimum one value selector to determine the map key
		return &CannotUnmarshalError{
			Reason: MissingValueSelector,
			V:      v,
		}
	}

	keyT := v.Type().Key()

	if keyT.Kind() != reflect.String {
		return &CannotUnmarshalError{
			Reason: NonStringMapKey,
			V:      v,
		}
	}

	kf := f.valFunc()
	eleT := v.Type().Elem()

	switch eleT.Kind() {
	case reflect.Slice, reflect.Array, reflect.Struct:
	default:
		s = childrenUntilMatch(s, f.selector(1))
		f = f.popVal()
	}

	var err error
	s.EachWithBreak(func(_ int, subS *Selection) bool {
		newV := reflect.New(eleT)

		err = unmarshalByType(subS, newV, f)
		if err != nil {
			return false
		}

		key := kf(subS)

		if eleT.Kind() != reflect.Ptr {
			newV = newV.Elem()
		}

		v.SetMapIndex(reflect.ValueOf(key), newV)

		return true
	})

	if err != nil {
		return &CannotUnmarshalError{
			Reason: TypeConversionError,
			Err:    err,
			V:      v,
		}
	}

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
