package got_test

import (
	"os"
	"path/filepath"
	"testing"

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

	m := &mock{t: t, name: t.Name()}
	gm := got.New(m)
	gm.Snapshot("a", "no")
	m.check(`"no" ⦗not ==⦘ "ok"`)
}

func TestSnapshotsCreate(t *testing.T) {
	err := os.RemoveAll(filepath.Join(".snapshots", "TestSnapshotsCreate.txt"))
	if err != nil {
		panic(err)
	}

	g := got.T(t)

	g.Snapshot("a", "ok")
}

func TestSnapshotsNotUsed(t *testing.T) {
	m := &mock{t: t, name: t.Name()}
	got.New(m)
	m.cleanup()

	if m.msg != "snapshot `a` is not used" {
		t.Fail()
	}
}
