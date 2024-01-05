package got

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ysmood/got/lib/utils"
)

const snapshotJSONExt = ".json"

type snapshot struct {
	value any
	used  bool
}

func (g G) snapshotsDir() string {
	return filepath.Join(".got", "snapshots", escapeFileName(g.Name()))
}

func (g G) loadSnapshots() {
	paths, err := filepath.Glob(filepath.Join(g.snapshotsDir(), "*"+snapshotJSONExt))
	g.E(err)

	for _, path := range paths {
		g.snapshots.Store(path, snapshot{g.JSON(g.Read(path)), false})
	}

	g.Cleanup(func() {
		if g.Failed() {
			return
		}

		g.snapshots.Range(func(path, data interface{}) bool {
			s := data.(snapshot)
			if !s.used {
				g.E(os.Remove(path.(string)))
			}
			return true
		})
	})
}

// Snapshot asserts that x equals the snapshot with the specified name, name should be unique under the same test case.
// It will create a new snapshot file if the name is not found.
// The snapshot file will be saved to ".got/snapshots/{TEST_NAME}".
// To update the snapshot, just change the name of the snapshot or remove the corresponding snapshot file.
// It will auto-remove the unused snapshot files after the test.
// The snapshot files should be version controlled.
// The format of the snapshot file is json.
func (g G) Snapshot(name string, x interface{}) {
	g.Helper()

	path := filepath.Join(g.snapshotsDir(), escapeFileName(name)+snapshotJSONExt)

	if data, ok := g.snapshots.Load(path); ok {
		s := data.(snapshot)
		xVal := g.JSON(g.ToJSON(x).Bytes())
		if utils.SmartCompare(xVal, s.value) == 0 {
			g.snapshots.Store(path, snapshot{x, true})
		} else {
			g.Assertions.err(AssertionEq, xVal, s.value)
		}
		return
	}

	g.snapshots.Store(path, snapshot{x, true})

	g.Cleanup(func() {
		g.E(os.MkdirAll(g.snapshotsDir(), 0755))
		g.E(os.WriteFile(path, g.ToJSON(x).Bytes(), 0644))
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
