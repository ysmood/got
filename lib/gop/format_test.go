package gop_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"text/template"
	"time"
	"unsafe"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/gop"
)

func TestStyle(t *testing.T) {
	g := got.T(t)
	s := gop.Style{Set: "<s>", Unset: "</s>"}

	g.Eq(gop.S("test", s), "<s>test</s>")
	g.Eq(gop.S("", s), "<s></s>")
	g.Eq(gop.S("", gop.None), "")
}

func TestTokenize(t *testing.T) {
	g := got.T(t)
	ref := "test"
	timeStamp, _ := time.Parse(time.RFC3339Nano, "2021-08-28T08:36:36.807908+08:00")
	fn := func(string) int { return 10 }
	ch1 := make(chan int)
	ch2 := make(chan string, 3)
	ch3 := make(chan struct{})

	v := []interface{}{
		nil,
		[]interface{}{true, false, uintptr(0x17), float32(100.121111133)},
		true, 10, int8(2), int32(100),
		float64(100.121111133),
		complex64(1 + 2i), complex128(1 + 2i),
		[3]int{1, 2},
		ch1,
		ch2,
		ch3,
		fn,
		map[interface{}]interface{}{
			`"test"`: 10,
			"a":      1,
		},
		unsafe.Pointer(&ref),
		struct {
			Int int
			str string
			M   map[int]int
		}{10, "ok", map[int]int{1: 0x20}},
		[]byte("aa\xe2"),
		[]byte("bytes\n\tbytes"),
		[]byte("long long long long string"),
		byte('a'),
		byte(1),
		'天',
		"long long long long string",
		"\ntest",
		"\t\n`",
		&ref,
		(*struct{ Int int })(nil),
		&struct{ Int int }{},
		&map[int]int{1: 2, 3: 4},
		&[]int{1, 2},
		&[2]int{1, 2},
		&[]byte{1, 2},
		timeStamp,
		time.Hour,
		`{"a": 1}`,
		[]byte(`{"a": 1}`),
	}

	check := func(out string, tpl ...string) {
		expected := bytes.NewBuffer(nil)

		t := template.New("")
		g.E(t.Parse(g.Read(g.Open(false, tpl...)).String()))
		g.E(t.Execute(expected, gop.Val{
			"ch1": fmt.Sprintf("0x%x", reflect.ValueOf(ch1).Pointer()),
			"ch2": fmt.Sprintf("0x%x", reflect.ValueOf(ch2).Pointer()),
			"ch3": fmt.Sprintf("0x%x", reflect.ValueOf(ch3).Pointer()),
			"fn":  fmt.Sprintf("0x%x", reflect.ValueOf(fn).Pointer()),
			"ptr": fmt.Sprintf("%v", &ref),
		}))

		g.Eq(out, expected.String())
	}

	out := gop.StripANSI(gop.F(v))

	{
		code := fmt.Sprintf(g.Read(g.Open(false, "fixtures", "compile_check.go.tmpl")).String(), out)
		f := g.Open(true, "tmp", g.RandStr(8), "main.go")
		g.Cleanup(func() { _ = os.Remove(f.Name()) })
		g.Write(code)(f)
		b, err := exec.Command("go", "run", f.Name()).CombinedOutput()
		if err != nil {
			g.Error(string(b))
		}
	}

	check(out, "fixtures", "expected.tmpl")

	out = gop.VisualizeANSI(gop.F(v))
	check(out, "fixtures", "expected_with_color.tmpl")
}

type A struct {
	Int int
	B   *B
}

type B struct {
	s string
	a *A
}

func TestCircularRef(t *testing.T) {
	g := got.T(t)
	a := A{Int: 10}
	b := B{"test", &a}
	a.B = &b

	g.Eq(gop.StripANSI(gop.F(a)), ""+
		"gop_test.A/* len=2 */{\n"+
		"    Int: 10,\n"+
		"    B: &gop_test.B/* len=2 */{\n"+
		"        s: \"test\",\n"+
		"        a: &gop_test.A/* len=2 */{\n"+
		"            Int: 10,\n"+
		"            B: gop.Circular(\"B\").(*gop_test.B),\n"+
		"        },\n"+
		"    },\n"+
		"}")
}

func TestCircularNilRef(t *testing.T) {
	arr := []A{{}, {}}

	got.T(t).Eq(gop.StripANSI(gop.F(arr)), `[]gop_test.A/* len=2 cap=2 */{
    gop_test.A/* len=2 */{
        Int: 0,
        B: (*gop_test.B)(nil),
    },
    gop_test.A/* len=2 */{
        Int: 0,
        B: (*gop_test.B)(nil),
    },
}`)
}

func TestCircularMap(t *testing.T) {
	g := got.T(t)
	a := map[int]interface{}{}
	a[0] = a

	ts := gop.Tokenize(a)

	g.Eq(gop.Format(ts, gop.ThemeNone), ""+
		"map[int]interface {}{\n"+
		"    0: gop.Circular().(map[int]interface {}),\n"+
		"}")
}

func TestCircularSlice(t *testing.T) {
	g := got.New(t)
	a := [][]interface{}{{nil}, {nil}}
	a[0][0] = a[1]
	a[1][0] = a[0][0]

	ts := gop.Tokenize(a)

	g.Eq(gop.Format(ts, gop.ThemeNone), ""+
		"[][]interface {}/* len=2 cap=2 */{\n"+
		"    gop.Arr/* len=1 cap=1 */{\n"+
		"        gop.Arr/* len=1 cap=1 */{\n"+
		"            gop.Circular(0, 0).(gop.Arr),\n"+
		"        },\n"+
		"    },\n"+
		"    gop.Circular(0, 0).(gop.Arr),\n"+
		"}")
}

func TestPlain(t *testing.T) {
	g := got.T(t)
	g.Eq(gop.Plain(10), "10")
}

func TestP(t *testing.T) {
	gop.Stdout = ioutil.Discard
	_ = gop.P("test")
	gop.Stdout = os.Stdout
}

func TestConvertors(t *testing.T) {
	g := got.T(t)
	g.Nil(gop.Circular(""))

	s := g.RandStr(8)
	g.Eq(gop.Ptr(s).(*string), &s)

	bs := base64.StdEncoding.EncodeToString([]byte(s))

	g.Eq(gop.Base64(bs), []byte(s))
	now := time.Now()
	g.Eq(gop.Time(now.Format(time.RFC3339Nano), 1234), now)
	g.Eq(gop.Duration("10m"), 10*time.Minute)

	g.Eq(gop.JSONStr(nil, "[1, 2]"), "[1, 2]")
	g.Eq(gop.JSONBytes(nil, "[1, 2]"), []byte("[1, 2]"))
}

func TestGetPrivateFieldErr(t *testing.T) {
	g := got.T(t)
	g.Panic(func() {
		gop.GetPrivateField(reflect.ValueOf(1), 0)
	})
	g.Panic(func() {
		gop.GetPrivateFieldByName(reflect.ValueOf(1), "test")
	})
}

func TestFixNestedStyle(t *testing.T) {
	g := got.T(t)

	s := gop.S(" 0 "+gop.S(" 1 "+
		gop.S(" 2 "+
			gop.S(" 3 ", gop.Cyan)+
			" 4 ", gop.Blue)+
		" 5 ", gop.Red)+" 6 ", gop.BgRed)
	fmt.Println(gop.VisualizeANSI(s))
	out := gop.VisualizeANSI(gop.FixNestedStyle(s))
	g.Eq(out, `<41> 0 <31> 1 <39><34> 2 <39><36> 3 <39><34> 4 <39><31> 5 <39> 6 <49>`)

	gop.FixNestedStyle("test")
}

func TestStripANSI(t *testing.T) {
	g := got.T(t)
	g.Eq(gop.StripANSI(gop.S("test", gop.Red)), "test")
}

func TestTheme(t *testing.T) {
	g := got.T(t)
	g.Eq(gop.ThemeDefault(gop.Error), []gop.Style{gop.Underline, gop.Red})
}
