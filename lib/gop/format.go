package gop

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Stdout is the default stdout for gop.P .
var Stdout io.Writer = os.Stdout

const indentUnit = "    "

// DefaultTheme colors for Sprint
var DefaultTheme = func(t Type) Color {
	switch t {
	case TypeName:
		return Cyan
	case Bool, Chan:
		return Blue
	case Rune, Byte, String:
		return Yellow
	case Number:
		return Green
	case Func:
		return Magenta
	case Comment:
		return White
	case Nil:
		return Red
	default:
		return None
	}
}

// NoTheme colors for Sprint
var NoTheme = func(t Type) Color {
	return None
}

// F is a shortcut for Format with color
func F(v interface{}) string {
	return Format(Tokenize(v), nil)
}

// P pretty print the values
func P(values ...interface{}) error {
	list := []interface{}{}
	for _, v := range values {
		list = append(list, F(v))
	}

	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()
	cwd, _ := os.Getwd()
	file, _ = filepath.Rel(cwd, file)
	tpl := ColorStr(DefaultTheme(Comment), "// %s %s:%d (%s)\n")
	_, _ = fmt.Fprintf(Stdout, tpl, time.Now().Format(time.RFC3339Nano), file, line, fn)

	_, err := fmt.Fprintln(Stdout, list...)
	return err
}

// Plain is a shortcut for Format with plain color
func Plain(v interface{}) string {
	return Format(Tokenize(v), NoTheme)
}

// Format a list of tokens
func Format(ts []*Token, theme func(Type) Color) string {
	if theme == nil {
		theme = DefaultTheme
	}

	out := ""
	depth := 0
	for i, t := range ts {
		if oneOf(t.Type, SliceOpen, MapOpen, StructOpen) {
			depth++
		}
		if i < len(ts)-1 && oneOf(ts[i+1].Type, SliceClose, MapClose, StructClose) {
			depth--
		}

		color := theme(t.Type)

		switch t.Type {
		case SliceOpen, MapOpen, StructOpen:
			out += ColorStr(color, t.Literal) + "\n"
		case SliceItem, MapKey, StructKey:
			out += strings.Repeat(indentUnit, depth)
		case Colon, InlineComma, Chan:
			out += ColorStr(color, t.Literal) + " "
		case Comma:
			out += ColorStr(color, t.Literal) + "\n"
		case SliceClose, MapClose, StructClose:
			out += strings.Repeat(indentUnit, depth) + ColorStr(color, t.Literal)
		case String:
			out += ColorStr(color, readableStr(depth, t.Literal))
		default:
			out += ColorStr(color, t.Literal)
		}
	}

	return out
}

func oneOf(t Type, list ...Type) bool {
	for _, el := range list {
		if t == el {
			return true
		}
	}
	return false
}

// To make multi-line string block more human readable.
// Split newline into two strings, convert "\t" into tab.
// Such as foramt string: "line one \n\t line two" into:
//     "line one \n" +
//     "	 line two"
func readableStr(depth int, s string) string {
	if ((len(s) > LongStringLen) || strings.Contains(s, "\n") || strings.Contains(s, `"`)) && !strings.Contains(s, "`") {
		return "`" + s + "`"
	}

	s = fmt.Sprintf("%#v", s)
	s, _ = replaceEscaped(s, 't', "	")

	indent := strings.Repeat(indentUnit, depth+1)
	if n, has := replaceEscaped(s, 'n', "\\n\" +\n"+indent+"\""); has {
		return "\"\" +\n" + indent + n
	}

	return s
}

// We use a simple state machine to replace escaped char like "\n"
func replaceEscaped(s string, escaped rune, new string) (string, bool) {
	type State int
	const (
		init State = iota
		prematch
		match
	)

	state := init
	out := ""
	buf := ""
	has := false

	onInit := func(r rune) {
		state = init
		out += buf + string(r)
		buf = ""
	}

	onPrematch := func() {
		state = prematch
		buf = "\\"
	}

	onEscape := func() {
		state = match
		out += new
		buf = ""
		has = true
	}

	for _, r := range s {
		switch state {
		case prematch:
			switch r {
			case escaped:
				onEscape()
			default:
				onInit(r)
			}

		case match:
			switch r {
			case '\\':
				onPrematch()
			default:
				onInit(r)
			}

		default:
			switch r {
			case '\\':
				onPrematch()
			default:
				onInit(r)
			}
		}
	}

	return out, has
}
