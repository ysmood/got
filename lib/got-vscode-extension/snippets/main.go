package main

import (
	"encoding/json"
	"os"
	"strings"
)

type snippet struct {
	Prefix      string   `json:"prefix,omitempty"`
	Body        []string `json:"body,omitempty"`
	Description string   `json:"description,omitempty"`
}

type snippets map[string]snippet

func main() {
	b, _ := os.ReadFile("lib/example/setup_test.go")

	setup := strings.Replace(string(b), "example_test", "${0:example}_test", -1)

	s := snippets{
		"got test function": {
			Prefix: "gt",
			Body: []string{`
func Test$1(t *testing.T) {
	g := setup(t)

	${0:g.Eq(1, 1)}
}`},
		},
		"got setup": {
			Prefix: "gts",
			Body:   []string{string(setup)},
		},
	}

	b, _ = json.Marshal(s)

	_ = os.WriteFile("lib/got-vscode-extension/snippets.json", b, 0764)
}
