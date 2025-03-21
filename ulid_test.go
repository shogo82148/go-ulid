package ulid

import (
	"runtime"
	"testing"
)

func TestMake(t *testing.T) {
	seen := make(map[ULID]struct{}, 0)
	for range 1000 {
		id := Make()
		if _, ok := seen[id]; ok {
			t.Fatalf("duplicate ULID: %v", id)
		}
		seen[id] = struct{}{}
	}
}

func BenchmarkMake(b *testing.B) {
	for b.Loop() {
		runtime.KeepAlive(Make())
	}
}

func BenchmarkString(b *testing.B) {
	id := Make()
	for b.Loop() {
		runtime.KeepAlive(id.String())
	}
}
