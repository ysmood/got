package diff

import (
	"context"
	"testing"
)

func TestFindHistAnchorCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	xs := []int{1, 2}
	ys := []int{1, 2}
	if m := findHistAnchor(ctx, xs, ys, 0, len(xs), 0, len(ys)); m.length != 0 {
		t.Fatalf("expected empty match on cancelled ctx, got %+v", m)
	}
}
