package diff_test

import (
	"strings"
	"testing"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/diff"
)

type T struct {
	got.G
}

func Test(t *testing.T) {
	got.Each(t, func(t *testing.T) T {
		return T{got.NewWith(t, got.NoColor().NoDiff())}
	})
}

func (t T) Format() {
	ts := diff.Tokenize(
		strings.ReplaceAll("a b c d f g h j q z", " ", "\n"),
		strings.ReplaceAll("a b c d e f g i j k r x y z", " ", "\n"),
	)

	df := diff.Format(ts, diff.NoTheme)

	t.Eq(df, ""+
		"01 01   a\n"+
		"02 02   b\n"+
		"03 03   c\n"+
		"04 04   d\n"+
		"   05 + e\n"+
		"05 06   f\n"+
		"06 07   g\n"+
		"07    - h\n"+
		"   08 + i\n"+
		"08 09   j\n"+
		"09    - q\n"+
		"   10 + k\n"+
		"   11 + r\n"+
		"   12 + x\n"+
		"   13 + y\n"+
		"10 14   z\n"+
		"")
}
