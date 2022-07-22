package base

import (
	"context"
	"testing"
)

func TestRepo(t *testing.T) {
	r := NewRepo("https://github.com/gole-dev/gole-layout.git")
	if err := r.Clone(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := r.CopyTo(context.Background(), "/tmp/test_gole_repo", "github.com/gole-dev/gole-layout", nil); err != nil {
		t.Fatal(err)
	}
}
