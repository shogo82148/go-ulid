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

func TestParse(t *testing.T) {
	t.Run("valid ulid", func(t *testing.T) {
		id, err := Parse("01ARZ3NDEKTSV4RRFFQ69G5FAV")
		if err != nil {
			t.Fatal(err)
		}
		if id != (ULID{0x01, 0x56, 0x3e, 0x3a, 0xb5, 0xd3, 0xd6, 0x76, 0x4c, 0x61, 0xef, 0xb9, 0x93, 0x02, 0xbd, 0x5b}) {
			t.Fatalf("id=%x", [16]byte(id))
		}
	})

	t.Run("invalid size", func(t *testing.T) {
		_, err := Parse("01ARZ3NDEKTSV4RRFFQ69G5FA")
		if err != ErrInvalidSize {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("invalid character", func(t *testing.T) {
		_, err := Parse("01ARZ3NDEKTSV4RRFFQ69G5FA!")
		if err != ErrInvalidCharacter {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("overflow", func(t *testing.T) {
		_, err := Parse("80000000000000000000000000")
		if err != ErrOverflow {
			t.Fatalf("err=%v", err)
		}
	})
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
