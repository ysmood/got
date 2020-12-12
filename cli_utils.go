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
// Return error if coverage is less than min, min is a percentage value.
func EnsureCoverage(path string, min float64) error {
	out, err := exec.Command("go", "tool", "cover", "-func="+path).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	list := strings.Split(strings.TrimSpace(string(out)), "\n")
	total := 0.0
	for _, l := range list {
		covStr := regexp.MustCompile(`(\d+.\d+)%\z`).FindStringSubmatch(l)[1]

		cov, _ := strconv.ParseFloat(string(covStr), 64)
		total += cov
	}
	total = total / float64(len(list))

	if compareFloat(total, min) < 0 {
		return fmt.Errorf("[lint] Test coverage %f%% must >= %f%%\n%s", total, min, out)
	}

	return nil
}

func compareFloat(a, b float64) int {
	return int(a*10000) - int(b*10000) // to avoid machine epsilon
}
