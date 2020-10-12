package got

import (
	"testing"
	"time"
)

func TestPanicAfter(t *testing.T) {
	ut := New(t)

	ut.Panic(func() {
		panicWithTrace(1)
	})

	wait := make(chan struct{})

	old := panicWithTrace
	panicWithTrace = func(v interface{}) {
		ut.Eq(v, "TestPanicAfter timeout after 1ns")
		close(wait)
	}
	defer func() { panicWithTrace = old }()

	ut.PanicAfter(1)
	time.Sleep(time.Millisecond)
	<-wait
}
