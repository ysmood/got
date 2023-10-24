package got

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ysmood/gop"
)

const snapshotExt = ".got-snap"

type snapshot struct {
	value string
	used  bool
}

func (g G) snapshotsDir() string {
	return filepath.Join(".got", "snapshots", escapeFileName(g.Name()))
}

func (g G) loadSnapshots() {
	paths, err := filepath.Glob(filepath.Join(g.snapshotsDir(), "*"+snapshotExt))
	g.E(err)

	for _, path := range paths {
		g.snapshots.Store(path, snapshot{g.Read(path).String(), false})
	}

	g.Cleanup(func() {
		g.snapshots.Range(func(path, data interface{}) bool {
			s := data.(snapshot)
			if !s.used {
				g.E(os.Remove(path.(string)))
			}
			return true
		})
	})
}

// Snapshot asserts that x equals the snapshot with the specified name, name should be unique under the same test.
// It will create a new snapshot file if the name is not found.
// The snapshot file will be saved to ".got/snapshots/{TEST_NAME}/{name}.got-snap".
// To update the snapshot, just change the name of the snapshot or remove the corresponding snapshot file.
// It will auto-remove the unused snapshot files after the test.
// The snapshot files should be version controlled.
func (g G) Snapshot(name string, x interface{}) {
	g.Helper()

	path := filepath.Join(g.snapshotsDir(), escapeFileName(name)+snapshotExt)

	xs := gop.Plain(x)

	if data, ok := g.snapshots.Load(path); ok {
		s := data.(snapshot)
		if xs == s.value {
			g.snapshots.Store(path, snapshot{xs, true})
		} else {
			g.Assertions.err(AssertionSnapshot, xs, s.value)
		}
		return
	}

	g.snapshots.Store(path, snapshot{xs, true})

	g.Cleanup(func() {
		g.E(os.MkdirAll(g.snapshotsDir(), 0755))
		g.E(os.WriteFile(path, []byte(xs), 0644))
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
