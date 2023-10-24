package got_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ysmood/gop"
	"github.com/ysmood/got"
)

func TestSnapshots(t *testing.T) {
	g := got.T(t)

	type C struct {
		Val int
	}

	g.Snapshot("a", "ok")
	g.Snapshot("b", 1)
	g.Snapshot("c", C{10})

	g.Run("sub", func(g got.G) {
		g.Snapshot("d", "ok")
	})

	m := &mock{t: t, name: t.Name()}
	gm := got.New(m)
	gm.Snapshot("a", "ok")
	gm.Snapshot("a", "no")
	m.check(`"no" ⦗not ==⦘ "ok"`)

	gm.Snapshot("a", "no\nno")
	g.Has(m.msg, "diff chunk")
	m.reset()

	gm.ErrorHandler = got.NewDefaultAssertionError(gop.ThemeNone, nil)
	gm.Snapshot("a", "no")
	m.checkWithStyle(true, `"no" ⦗not ==⦘ "ok"`)
}

func TestSnapshotsCreate(t *testing.T) {
	path := filepath.FromSlash(".got/snapshots/TestSnapshotsCreate/a.got-snap")
	err := os.RemoveAll(path)
	if err != nil {
		panic(err)
	}

	g := got.T(t)

	g.Cleanup(func() {
		g.True(g.PathExists(path))
	})

	g.Snapshot("a", "ok")
}

func TestSnapshotsNotUsed(t *testing.T) {
	path := filepath.FromSlash(".got/snapshots/TestSnapshotsNotUsed/a.got-snap")

	g := got.T(t)
	g.WriteFile(path, []byte(`1`))

	m := &mock{t: t, name: t.Name()}
	got.New(m)
	m.cleanup()

	g.False(g.PathExists(path))
}
