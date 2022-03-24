package got_test

import (
	"testing"

	"github.com/ysmood/got"
)

func TestEnsureCoverage(t *testing.T) {
	g := setup(t)
	g.Nil(got.EnsureCoverage("fixtures/coverage/cov.txt", 100))

	g.Err(got.EnsureCoverage("fixtures/coverage/cov.txt", 120))
	g.Err(got.EnsureCoverage("not-exists", 100))
}
