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

func BenchmarkParse(b *testing.B) {
	const s = "0000XSNJG0MQJHBF4QX1EFD6Y3"
	for b.Loop() {
		id, err := Parse(s)
		if err != nil {
			b.Fatal(err)
		}
		runtime.KeepAlive(id)
	}
}

func FuzzParse(f *testing.F) {
	f.Add("0000XSNJG0MQJHBF4QX1EFD6Y3")
	f.Add("01ARZ3NDEKTSV4RRFFQ69G5FAV")
	f.Fuzz(func(t *testing.T, s string) {
		id0, err := Parse(s)
		if err != nil {
			t.Skip()
		}
		id1, err := Parse(id0.String())
		if err != nil {
			t.Fatal(err)
		}
		if id0 != id1 {
			t.Fatalf("id0=%v id1=%v", id0, id1)
		}
	})
}

func BenchmarkString(b *testing.B) {
	id := Make()
	for b.Loop() {
		runtime.KeepAlive(id.String())
	}
}
