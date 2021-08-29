package gop_test

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/gop"
)

type T struct {
	got.G
}

func Test(t *testing.T) {
	got.Each(t, T{})
}

func (t T) Tokenize() {
	ref := "test"
	timeStamp, _ := time.Parse(time.RFC3339Nano, "2021-08-28T08:36:36.807908+08:00")

	v := []interface{}{
		nil,
		[]interface{}{true, false, uintptr(0x17), float32(100.121111133)},
		true, 10, int8(2), int32(100),
		float64(100.121111133), complex(1, 2),
		[3]int{1, 2},
		make(chan int),
		make(chan string, 3),
		func(string) int { return 10 },
		map[interface{}]interface{}{
			"test": 10,
			"a":    1,
		},
		unsafe.Pointer(&ref),
		struct {
			Int int
			str string
			M   map[int]int
		}{10, "ok", map[int]int{1: 0x20}},
		[]byte("aa\xe2"),
		[]byte("bytes\n\tbytes"),
		byte('a'),
		byte(1),
		'天',
		"\ntest",
		&ref,
		(*struct{ Int int })(nil),
		&struct{ Int int }{},
		&map[int]int{1: 2, 3: 4},
		&[]int{1, 2},
		&[2]int{1, 2},
		&[]byte{1, 2},
		timeStamp,
		time.Hour,
	}

	out := gop.StripColor(gop.F(v))

	t.Eq(out, `[]interface {}{
    nil,
    []interface {}{
        true,
        false,
        uintptr(23),
        float32(100.12111),
    },
    true,
    10,
    int8(2),
    'd',
    float64(100.121111133),
    complex128((1+2i)),
    [3]int{
        1,
        2,
        0,
    },
    make(chan int),
    make(chan string, 3),
    (func(string) int)(nil),
    map[interface {}]interface {}{
        "a": 1,
        "test": 10,
    },
    unsafe.Pointer(uintptr(`+fmt.Sprintf("%v", &ref)+`)),
    struct { Int int; str string; M map[int]int }{
        Int: 10,
        str: "ok",
        M: map[int]int{
            1: 32,
        },
    },
    gop.Base64("YWHi"),
    []byte("" +
        "bytes\n" +
        "\tbytes"),
    byte('a'),
    byte(0x1),
    '天',
    "" +
        "\n" +
        "test",
    gop.ToPtr("test").(*string),
    (*struct { Int int })(nil),
    &struct { Int int }{
        Int: 0,
    },
    &map[int]int{
        1: 2,
        3: 4,
    },
    &[]int{
        1,
        2,
    },
    &[2]int{
        1,
        2,
    },
    &[]byte("\x01\x02"),
    time.Parse(time.RFC3339Nano, "`+timeStamp.Format(time.RFC3339Nano)+`"),
    Duration(1h0m0s),
}`)
}

type A struct {
	Int int
	B   *B
}

type B struct {
	s string
	a *A
}

func (t T) Cyclic() {
	a := A{Int: 10}
	b := B{"test", &a}
	a.B = &b

	ts := gop.Tokenize(a)

	t.Has(gop.Format(ts, gop.NoTheme), `gop.Cyclic(`)
}

func (t T) Plain() {
	t.Eq(gop.Plain(10), "10")
}

func (t T) P() {
	gop.Stdout = io.Discard
	_, _ = gop.P("test")
	gop.Stdout = os.Stdout
}

func (t T) Others() {
	gop.ToPtr(nil)
	_ = gop.Cyclic(0)
	_ = gop.Base64("")
}

func (t T) GetPrivateFieldErr() {
	t.Panic(func() {
		gop.GetPrivateField(reflect.ValueOf(1), 0)
	})
}

func (t T) Lab() {

}
