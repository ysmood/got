package got_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ysmood/got"
)

func TestSnapshots(t *testing.T) {
	m := &mock{t: t, name: t.Name()}
	g := got.T(m)

	g.Snapshot("a", "ok")
	g.Snapshot("b", 1)

	g.Snapshot("a", "no")
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
