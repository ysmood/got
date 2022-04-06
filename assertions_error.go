package got

import (
	"strings"

	"github.com/ysmood/got/lib/diff"
	"github.com/ysmood/got/lib/gop"
)

// AssertionErrType enum
type AssertionErrType int

const (
	// AssertionEq type
	AssertionEq AssertionErrType = iota
	// AssertionNeqSame type
	AssertionNeqSame
	// AssertionNeq type
	AssertionNeq
	// AssertionGt type
	AssertionGt
	// AssertionGte type
	AssertionGte
	// AssertionLt type
	AssertionLt
	// AssertionLte type
	AssertionLte
	// AssertionInDelta type
	AssertionInDelta
	// AssertionTrue type
	AssertionTrue
	// AssertionFalse type
	AssertionFalse
	// AssertionNil type
	AssertionNil
	// AssertionNoArgs type
	AssertionNoArgs
	// AssertionNotNil type
	AssertionNotNil
	// AssertionNotNilable type
	AssertionNotNilable
	// AssertionNotNilableNil type
	AssertionNotNilableNil
	// AssertionZero type
	AssertionZero
	// AssertionNotZero type
	AssertionNotZero
	// AssertionRegex type
	AssertionRegex
	// AssertionHas type
	AssertionHas
	// AssertionLen type
	AssertionLen
	// AssertionErr type
	AssertionErr
	// AssertionPanic type
	AssertionPanic
	// AssertionIsInChain type
	AssertionIsInChain
	// AssertionIsKind type
	AssertionIsKind
	// AssertionCount type
	AssertionCount
)

// AssertionCtx holds the context of an assertion
type AssertionCtx struct {
	Type    AssertionErrType
	Details []interface{}
	File    string
	Line    int
}

// AssertionError handler
type AssertionError interface {
	Report(*AssertionCtx) string
}

var _ AssertionError = AssertionErrorReport(nil)

// AssertionErrorReport is used to convert a func to AssertionError
type AssertionErrorReport func(*AssertionCtx) string

// Report interface
func (ae AssertionErrorReport) Report(ac *AssertionCtx) string {
	return ae(ac)
}

type defaultAssertionError struct {
	fns map[AssertionErrType]func(details ...interface{}) string
}

// NewDefaultAssertionError handler
func NewDefaultAssertionError(enableColor, enableDiff bool) AssertionError {
	f := func(v interface{}) string {
		if enableColor {
			return gop.F(v)
		}
		return gop.Plain(v)
	}

	k := func(s string) string {
		s = " ⦗" + s + "⦘ "
		if enableColor {
			return gop.ColorStr(gop.Red, s)
		}
		return s
	}

	fns := map[AssertionErrType]func(details ...interface{}) string{
		AssertionEq: func(details ...interface{}) string {
			x := f(details[0])
			y := f(details[1])
			if enableDiff && hasNewline(x, y) {
				df := ""
				if !enableColor {
					df = diff.Format(diff.TokenizeText(x, y), diff.NoTheme)
				} else {
					df = diff.Diff(x, y)
				}
				return j(x, k("not =="), y, df)
			}
			return j(x, k("not =="), y)
		},
		AssertionNeqSame: func(details ...interface{}) string {
			x := f(details[0])
			y := f(details[1])
			return j(x, k("=="), y)
		},
		AssertionNeq: func(details ...interface{}) string {
			x := f(details[0])
			y := f(details[1])
			return j(x, k("=="), y, k("when converted to the same type"))
		},
		AssertionGt: func(details ...interface{}) string {
			x := f(details[0])
			y := f(details[1])
			return j(x, k("not >"), y)
		},
		AssertionGte: func(details ...interface{}) string {
			x := f(details[0])
			y := f(details[1])
			return j(x, k("not ≥"), y)
		},
		AssertionLt: func(details ...interface{}) string {
			x := f(details[0])
			y := f(details[1])
			return j(x, k("not <"), y)
		},
		AssertionLte: func(details ...interface{}) string {
			x := f(details[0])
			y := f(details[1])
			return j(x, k("not ≤"), y)
		},
		AssertionInDelta: func(details ...interface{}) string {
			x := f(details[0])
			y := f(details[1])
			delta := f(details[2])
			return j(k("delta between"), x, k("and"), y, k("not ≤"), delta)
		},
		AssertionTrue: func(_ ...interface{}) string {
			return k("should be") + f(true)
		},
		AssertionFalse: func(_ ...interface{}) string {
			return k("should be") + f(false)
		},
		AssertionNil: func(details ...interface{}) string {
			last := f(details[0])
			return j(k("last argument"), last, k("should be"), f(nil))
		},
		AssertionNoArgs: func(_ ...interface{}) string {
			return k("no arguments received")
		},
		AssertionNotNil: func(_ ...interface{}) string {
			return k("last argument shouldn't be") + f(nil)
		},
		AssertionNotNilable: func(details ...interface{}) string {
			last := f(details[0])
			return j(k("last argument"), last, k("is not nilable"))
		},
		AssertionNotNilableNil: func(details ...interface{}) string {
			last := f(details[0])
			return j(k("last argument"), last, k("shouldn't be"), f(nil))
		},
		AssertionZero: func(details ...interface{}) string {
			x := f(details[0])
			return j(x, k("should be zero value for its type"))
		},
		AssertionNotZero: func(details ...interface{}) string {
			x := f(details[0])
			return j(x, k("shouldn't be zero value for its type"))
		},
		AssertionRegex: func(details ...interface{}) string {
			pattern := f(details[0])
			str := f(details[1])
			return j(pattern, k("should match"), str)
		},
		AssertionHas: func(details ...interface{}) string {
			container := f(details[0])
			str := f(details[1])
			return j(container, k("should has"), str)
		},
		AssertionLen: func(details ...interface{}) string {
			actual := f(details[0])
			l := f(details[1])
			return k("expect len") + actual + k("to be") + l
		},
		AssertionErr: func(details ...interface{}) string {
			last := f(details[0])
			return j(k("last value"), last, k("should be <error>"))
		},
		AssertionPanic: func(_ ...interface{}) string {
			return k("should panic")
		},
		AssertionIsInChain: func(details ...interface{}) string {
			x := f(details[0])
			y := f(details[1])
			return j(x, k("should in chain of"), y)
		},
		AssertionIsKind: func(details ...interface{}) string {
			x := f(details[0])
			y := f(details[1])
			return j(x, k("should be kind of"), y)
		},
		AssertionCount: func(details ...interface{}) string {
			n := f(details[0])
			count := f(details[1])
			return k("should count") + n + k("times, but got") + count
		},
	}

	return &defaultAssertionError{fns: fns}
}

// Report interface
func (ae *defaultAssertionError) Report(ac *AssertionCtx) string {
	return ae.fns[ac.Type](ac.Details...)
}

func j(args ...string) string {
	if hasNewline(args...) {
		return "\n" + strings.Join(args, "\n\n") + "\n"
	}
	return strings.Join(args, "")
}

func hasNewline(args ...string) bool {
	for _, arg := range args {
		if strings.Contains(arg, "\n") {
			return true
		}
	}
	return false
}
