package gop

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"unicode"
	"unicode/utf8"
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
	return tokenize(seen{}, reflect.ValueOf(v))
}

// ToPtr converts a Go value to its pointer
func ToPtr(interface{}) interface{} {
	return nil
}

// Cyclic reference
func Cyclic(uintptr) interface{} {
	return nil
}

// Base64 from string
func Base64(string) []byte {
	return nil
}

type seen map[interface{}]struct{}

func tokenize(sn seen, v reflect.Value) []*Token {
	ts := []*Token{}
	t := &Token{Nil, ""}

	if v.Kind() == reflect.Invalid {
		t.Literal = "nil"
		ts = append(ts, t)
		return ts
	} else if r, ok := v.Interface().(rune); ok && unicode.IsGraphic(r) {
		ts = append(ts, tokenizeRune(t, r))
		return ts
	} else if b, ok := v.Interface().(byte); ok {
		ts = append(ts, tokenizeByte(t, b))
		return ts
	}

	switch v.Kind() {
	case reflect.Interface:
		ts = append(ts, tokenize(sn, v.Elem())...)

	case reflect.Slice, reflect.Array:
		if data, ok := v.Interface().([]byte); ok {
			ts = append(ts, tokenizeBytes(data)...)
			break
		} else {
			ts = append(ts, &Token{TypeName, v.Type().String()})
		}

		ts = append(ts, &Token{SliceOpen, "{"})
		for i := 0; i < v.Len(); i++ {
			el := v.Index(i)
			ts = append(ts, &Token{SliceItem, ""})
			ts = append(ts, tokenize(sn, el)...)
			ts = append(ts, &Token{Comma, ","})
		}
		ts = append(ts, &Token{SliceClose, "}"})

	case reflect.Map:
		ts = append(ts, &Token{TypeName, v.Type().String()})
		ts = append(ts, &Token{MapOpen, "{"})
		it := v.MapRange()
		for it.Next() {
			ts = append(ts, &Token{MapKey, ""})
			ts = append(ts, tokenize(sn, it.Key())...)
			ts = append(ts, &Token{Colon, ":"})
			ts = append(ts, tokenize(sn, it.Value())...)
			ts = append(ts, &Token{Comma, ","})
		}
		ts = append(ts, &Token{MapClose, "}"})

	case reflect.Struct:
		t := v.Type()

		ts = append(ts, &Token{TypeName, t.String()})
		ts = append(ts, &Token{StructOpen, "{"})
		for i := 0; i < v.NumField(); i++ {
			ts = append(ts, &Token{StructKey, ""})
			ts = append(ts, &Token{StructField, t.Field(i).Name})

			f := v.Field(i)
			if !f.CanInterface() {
				f = GetPrivateField(v, i)
			}
			ts = append(ts, &Token{Colon, ":"})
			ts = append(ts, tokenize(sn, f)...)
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
		reflect.Complex64, reflect.Complex128,
		reflect.Uintptr:

		ts = append(ts, &Token{TypeName, v.Type().Name()})
		ts = append(ts, &Token{ParenOpen, "("})
		t.Type = Number
		t.Literal = fmt.Sprintf("%v", v.Interface())
		ts = append(ts, t)
		ts = append(ts, &Token{ParenClose, ")"})

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
		ts = append(ts, tokenizePtr(sn, v)...)

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

func tokenizePtr(sn seen, v reflect.Value) []*Token {
	ts := []*Token{}

	if v.Elem().Kind() == reflect.Invalid {
		ts = append(ts, &Token{Nil, fmt.Sprintf("(%s)(nil)", v.Type().String())})
		return ts
	}

	if _, has := sn[v.Interface()]; has {
		ts = append(ts, &Token{
			PointerCyclic,
			fmt.Sprintf("gop.Cyclic(0x%x).(%s)", v.Pointer(), v.Type().String()),
		})
		return ts
	}

	sn[v.Interface()] = struct{}{}

	switch v.Elem().Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		ts = append(ts, &Token{TypeName, "&"})
		ts = append(ts, tokenize(sn, v.Elem())...)
	default:
		ts = append(ts, &Token{PointerOpen, "gop.ToPtr("})
		ts = append(ts, tokenize(sn, v.Elem())...)
		ts = append(ts, &Token{PointerOpen, fmt.Sprintf(").(%s)", v.Type().String())})
	}

	return ts
}
