package gop

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/ysmood/got/lib/utils"
)

// Type of token
type Type int

const (
	// Nil type
	Nil Type = iota
	// Bool type
	Bool
	// Number type
	Number
	// Float type
	Float
	// Complex type
	Complex
	// String type
	String
	// Byte type
	Byte
	// Rune type
	Rune
	// Chan type
	Chan
	// Func type
	Func
	// UnsafePointer type
	UnsafePointer

	// TypeName type
	TypeName

	// ParenOpen type
	ParenOpen
	// ParenClose type
	ParenClose

	// PointerOpen type
	PointerOpen
	// PointerClose type
	PointerClose
	// PointerCyclic type
	PointerCyclic

	// SliceOpen type
	SliceOpen
	// SliceItem type
	SliceItem
	// InlineComma type
	InlineComma
	// Comma type
	Comma
	// SliceClose type
	SliceClose

	// MapOpen type
	MapOpen
	// MapKey type
	MapKey
	// Colon type
	Colon
	// MapClose type
	MapClose

	// StructOpen type
	StructOpen
	// StructKey type
	StructKey
	// StructField type
	StructField
	// StructClose type
	StructClose
)

// Token represents a symbol in value layout
type Token struct {
	Type    Type
	Literal string
}

// Tokenize a random Go value
func Tokenize(v interface{}) []*Token {
	return tokenize(seen{}, []interface{}{}, reflect.ValueOf(v))
}

// ToPtr converts a Go value to its pointer
func ToPtr(interface{}) interface{} {
	return nil
}

// Cyclic reference of the path from the root
func Cyclic(path ...interface{}) interface{} {
	return nil
}

// Base64 from string
func Base64(string) []byte {
	return nil
}

// Time from string
func Time(string) time.Time {
	return time.Time{}
}

// Duration from string
func Duration(string) time.Duration {
	return 0
}

type path []interface{}

func (p path) String() string {
	out := []string{}
	for _, seg := range p {
		out = append(out, fmt.Sprintf("%#v", seg))
	}
	return strings.Join(out, ", ")
}

type seen map[uintptr]path

func (sn seen) cyclic(p path, v reflect.Value) *Token {
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice:
		ptr := v.Pointer()
		if p, has := sn[ptr]; has {
			return &Token{
				PointerCyclic,
				fmt.Sprintf("gop.Cyclic(%s).(%s)", p.String(), v.Type().String()),
			}
		}
		sn[ptr] = p
	}

	return nil
}

func tokenize(sn seen, p path, v reflect.Value) []*Token {
	ts := []*Token{}
	t := &Token{Nil, ""}

	if v.Kind() == reflect.Invalid {
		t.Literal = "nil"
		return append(ts, t)
	} else if r, ok := v.Interface().(rune); ok && unicode.IsGraphic(r) {
		return append(ts, tokenizeRune(t, r))
	} else if b, ok := v.Interface().(byte); ok {
		return append(ts, tokenizeByte(t, b))
	} else if tt, ok := v.Interface().(time.Time); ok {
		return tokenizeTime(tt)
	} else if d, ok := v.Interface().(time.Duration); ok {
		return tokenizeDuration(d)
	}

	if t := sn.cyclic(p, v); t != nil {
		return append(ts, t)
	}

	switch v.Kind() {
	case reflect.Interface:
		ts = append(ts, tokenize(sn, p, v.Elem())...)

	case reflect.Slice, reflect.Array:
		if data, ok := v.Interface().([]byte); ok {
			ts = append(ts, tokenizeBytes(data)...)
			break
		} else {
			ts = append(ts, &Token{TypeName, v.Type().String()})
		}

		ts = append(ts, &Token{SliceOpen, "{"})
		for i := 0; i < v.Len(); i++ {
			p := append(p, i)
			el := v.Index(i)
			ts = append(ts, &Token{SliceItem, ""})
			ts = append(ts, tokenize(sn, p, el)...)
			ts = append(ts, &Token{Comma, ","})
		}
		ts = append(ts, &Token{SliceClose, "}"})

	case reflect.Map:
		ts = append(ts, &Token{TypeName, v.Type().String()})
		ts = append(ts, &Token{MapOpen, "{"})
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return utils.Compare(keys[i], keys[j]) < 0
		})
		for _, k := range keys {
			p := append(p, k)
			ts = append(ts, &Token{MapKey, ""})
			ts = append(ts, tokenize(sn, p, k)...)
			ts = append(ts, &Token{Colon, ":"})
			ts = append(ts, tokenize(sn, p, v.MapIndex(k))...)
			ts = append(ts, &Token{Comma, ","})
		}
		ts = append(ts, &Token{MapClose, "}"})

	case reflect.Struct:
		t := v.Type()

		ts = append(ts, &Token{TypeName, t.String()})
		ts = append(ts, &Token{StructOpen, "{"})
		for i := 0; i < v.NumField(); i++ {
			name := t.Field(i).Name
			ts = append(ts, &Token{StructKey, ""})
			ts = append(ts, &Token{StructField, name})

			f := v.Field(i)
			if !f.CanInterface() {
				f = GetPrivateField(v, i)
			}
			ts = append(ts, &Token{Colon, ":"})
			ts = append(ts, tokenize(sn, append(p, name), f)...)
			ts = append(ts, &Token{Comma, ","})
		}
		ts = append(ts, &Token{StructClose, "}"})

	case reflect.Bool:
		t.Type = Bool
		if v.Bool() {
			t.Literal = "true"
		} else {
			t.Literal = "false"
		}
		ts = append(ts, t)

	case reflect.Int:
		t.Type = Number
		t.Literal = strconv.FormatInt(v.Int(), 10)
		ts = append(ts, t)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Uintptr:

		ts = append(ts, &Token{TypeName, v.Type().Name()})
		ts = append(ts, &Token{ParenOpen, "("})
		t.Type = Number
		t.Literal = fmt.Sprintf("%v", v.Interface())
		ts = append(ts, t)
		ts = append(ts, &Token{ParenClose, ")"})

	case reflect.Complex64:
		ts = append(ts, &Token{TypeName, v.Type().Name()})
		ts = append(ts, &Token{ParenOpen, "("})
		t.Type = Number
		t.Literal = fmt.Sprintf("%v", v.Interface())
		t.Literal = t.Literal[1 : len(t.Literal)-1]
		ts = append(ts, t)
		ts = append(ts, &Token{ParenClose, ")"})

	case reflect.Complex128:
		t.Type = Number
		t.Literal = fmt.Sprintf("%v", v.Interface())
		t.Literal = t.Literal[1 : len(t.Literal)-1]
		ts = append(ts, t)

	case reflect.String:
		t.Type = String
		t.Literal = fmt.Sprintf("%#v", v.Interface())
		ts = append(ts, t)

	case reflect.Chan:
		t.Type = Chan
		if v.Cap() == 0 {
			t.Literal = fmt.Sprintf("make(chan %s)", v.Type().Elem().Name())
		} else {
			t.Literal = fmt.Sprintf("make(chan %s, %d)", v.Type().Elem().Name(), v.Cap())
		}
		ts = append(ts, t)

	case reflect.Func:
		t.Type = Func
		t.Literal = fmt.Sprintf("(%s)(nil)", v.Type().String())
		ts = append(ts, t)

	case reflect.Ptr:
		ts = append(ts, tokenizePtr(sn, p, v)...)

	case reflect.UnsafePointer:
		t.Type = UnsafePointer
		t.Literal = fmt.Sprintf("unsafe.Pointer(uintptr(%v))", v.Interface())
		ts = append(ts, t)
	}

	return ts
}

func tokenizeRune(t *Token, r rune) *Token {
	t.Type = Rune
	t.Literal = fmt.Sprintf("'%s'", string(r))
	return t
}

func tokenizeByte(t *Token, b byte) *Token {
	t.Type = Byte
	if unicode.IsGraphic(rune(b)) {
		t.Literal = fmt.Sprintf("byte('%s')", string(b))
	} else {
		t.Literal = fmt.Sprintf("byte(0x%x)", b)
	}
	return t
}

func tokenizeTime(t time.Time) []*Token {
	ts := []*Token{}
	ts = append(ts, &Token{TypeName, "gop.Time"})
	ts = append(ts, &Token{ParenOpen, "("})
	ts = append(ts, &Token{String, `"` + t.Format(time.RFC3339Nano) + `"`})
	ts = append(ts, &Token{ParenClose, ")"})
	return ts
}

func tokenizeDuration(d time.Duration) []*Token {
	ts := []*Token{}
	ts = append(ts, &Token{TypeName, "gop.Duration"})
	ts = append(ts, &Token{ParenOpen, "("})
	ts = append(ts, &Token{String, `"` + d.String() + `"`})
	ts = append(ts, &Token{ParenClose, ")"})
	return ts
}

func tokenizeBytes(data []byte) []*Token {
	ts := []*Token{}

	if utf8.Valid(data) {
		ts = append(ts, &Token{TypeName, "[]byte"})
		ts = append(ts, &Token{ParenOpen, "("})
		ts = append(ts, &Token{String, fmt.Sprintf("%#v", string(data))})
		ts = append(ts, &Token{ParenClose, ")"})
		return ts
	}

	ts = append(ts, &Token{ParenOpen, "gop.Base64("})
	ts = append(ts, &Token{String, fmt.Sprintf("%#v", base64.StdEncoding.EncodeToString(data))})
	ts = append(ts, &Token{ParenClose, ")"})
	return ts
}

func tokenizePtr(sn seen, p path, v reflect.Value) []*Token {
	ts := []*Token{}

	if v.Elem().Kind() == reflect.Invalid {
		ts = append(ts, &Token{Nil, fmt.Sprintf("(%s)(nil)", v.Type().String())})
		return ts
	}

	fn := false

	switch v.Elem().Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		if _, ok := v.Elem().Interface().([]byte); ok {
			fn = true
		}
	default:
		fn = true
	}

	if fn {
		ts = append(ts, &Token{PointerOpen, "gop.ToPtr("})
		ts = append(ts, tokenize(sn, p, v.Elem())...)
		ts = append(ts, &Token{PointerOpen, fmt.Sprintf(").(%s)", v.Type().String())})
	} else {
		ts = append(ts, &Token{TypeName, "&"})
		ts = append(ts, tokenize(sn, p, v.Elem())...)
	}

	return ts
}
