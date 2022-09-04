// Package main ...
package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type snippet struct {
	Prefix      string   `json:"prefix,omitempty"`
	Body        []string `json:"body,omitempty"`
	Description string   `json:"description,omitempty"`
}

type snippets map[string]snippet

func main() {
	b, err := ioutil.ReadFile("lib/example/03_setup_test.go")
	if err != nil {
		log.Fatal(err)
	}

	setup := strings.Replace(string(b), "example_test", "${0:example}_test", -1)

	s := snippets{
		"gop print": {
			Prefix: "gp",
			Body:   []string{"gop.P($0)"},
		},
		"got test function": {
			Prefix: "gt",
			Body: strings.Split(`
func Test$1(t *testing.T) {
	g := got.T(t)

	${0:g.Eq(1, 1)}
}
`, "\n"),
		},
		"got test function with setup": {
			Prefix: "gts",
			Body: strings.Split(`
func Test$1(t *testing.T) {
	g := setup(t)

	${0:g.Eq(1, 1)}
}
`, "\n"),
		},
		"got setup": {
			Prefix: "gsetup",
			Body:   strings.Split(string(setup), "\n"),
		},
	}

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err = enc.Encode(s)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("lib/got-vscode-extension/snippets.json", buf.Bytes(), 0764)
	if err != nil {
		log.Fatal(err)
	}
}
