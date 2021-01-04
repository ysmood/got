package got

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// EnsureCoverage via report file generated from, for example:
//     go test -coverprofile=coverage.txt
// Return error if any functions's coverage is less than min, min is a percentage value.
func EnsureCoverage(path string, min float64) error {
	out, err := exec.Command("go", "tool", "cover", "-func="+path).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	list := strings.Split(strings.TrimSpace(string(out)), "\n")
	rejected := []string{}
	for _, l := range list {
		if strings.HasPrefix(l, "total:") {
			continue
		}

		covStr := regexp.MustCompile(`(\d+.\d+)%\z`).FindStringSubmatch(l)[1]

		cov, _ := strconv.ParseFloat(string(covStr), 64)
		if compareFloat(cov, min) < 0 {
			rejected = append(rejected, l)
		}
	}

	if len(rejected) > 0 {
		return fmt.Errorf(
			"[lint] Test coverage for these functions should be greater than %.2f%%:\n%s",
			min,
			strings.Join(rejected, "\n"),
		)
	}
	return nil
}

func compareFloat(a, b float64) int {
	return int(a*10000) - int(b*10000) // to avoid machine epsilon
}
