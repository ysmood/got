// Package example ...
package example

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

// Sum a and b as number
func Sum(a, b string) string {
	an, _ := strconv.ParseInt(a, 10, 32)
	bn, _ := strconv.ParseInt(b, 10, 32)
	return fmt.Sprintf("%d", an+bn)
}

// ServeSum http handler function
func ServeSum(w http.ResponseWriter, r *http.Request) {
	s := Sum(r.URL.Query().Get("a"), r.URL.Query().Get("b"))
	_, _ = w.Write([]byte(s))
}

// Rand generate a random int as string
func Rand(src rand.Source) string {
	return fmt.Sprintf("%d", rand.New(src).Int())
}
