package got

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type snapshots struct {
	list *sync.Map
	file io.ReadWriter
}

type snapshot struct {
	Name  string
	Value interface{}
}

// snapshotsFilePath returns the path of the snapshot file for current test.
func (g G) snapshotsFilePath() string {
	return filepath.Join(".snapshots", escapeFileName(g.Name())+".txt")
}

func (g G) loadSnapshots() {
	p := g.snapshotsFilePath()

	g.snapshots.list = &sync.Map{}

	if g.PathExists(p) {
		f, err := os.OpenFile(p, os.O_RDWR, 0755)
		g.E(err)
		g.snapshots.file = f
	} else {
		return
	}

	dec := json.NewDecoder(g.snapshots.file)

	for {
		var data snapshot
		err := dec.Decode(&data)
		if err == io.EOF {
			return
		}
		g.E(err)
		g.snapshots.list.Store(data.Name, data.Value)
	}
}

// Snapshot asserts that x equals the snapshot with the specified name, name should be unique.
// It can only compare JSON serializable types.
// It will create a new snapshot if the name is not found.
// The snapshot will be saved to ".snapshots/{TEST_NAME}" beside the test file,
// TEST_NAME is the current test name.
// To update the snapshot, just delete the corresponding file.
func (g G) Snapshot(name string, x interface{}) {
	g.Helper()

	if y, ok := g.snapshots.list.Load(name); ok {
		g.Eq(x, y)
		return
	}

	g.snapshots.list.Store(name, x)

	g.Cleanup(func() {
		if g.snapshots.file == nil {
			p := g.snapshotsFilePath()

			err := os.MkdirAll(filepath.Dir(p), 0755)
			g.E(err)

			f, err := os.Create(p)
			g.E(err)
			g.snapshots.file = f
		}

		enc := json.NewEncoder(g.snapshots.file)
		enc.SetIndent("", "  ")
		g.E(enc.Encode(snapshot{name, x}))
	})
}

func escapeFileName(fileName string) string {
	// Define the invalid characters for both Windows and Unix
	invalidChars := `< > : " / \ | ? *`

	// Replace the invalid characters with an underscore
	regex := "[" + regexp.QuoteMeta(invalidChars) + "]"
	escapedFileName := regexp.MustCompile(regex).ReplaceAllString(fileName, "_")

	// Remove any leading or trailing spaces or dots
	escapedFileName = strings.Trim(escapedFileName, " .")

	// Remove consecutive dots
	escapedFileName = regexp.MustCompile(`\.{2,}`).ReplaceAllString(escapedFileName, ".")

	return escapedFileName
}
