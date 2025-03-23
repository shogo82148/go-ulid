package ulid

import (
	"bytes"
	"runtime"
	"testing"
)

func TestMake(t *testing.T) {
	seen := make(map[ULID]struct{}, 0)
	for range 10000 {
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

func TestSetTime(t *testing.T) {
	var id ULID
	id.SetTime(0x1563e3ab5d3)
	if id != (ULID{0x01, 0x56, 0x3e, 0x3a, 0xb5, 0xd3, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) {
		t.Fatalf("id=%x", [16]byte(id))
	}
}

func TestTime(t *testing.T) {
	id := ULID{0x01, 0x56, 0x3e, 0x3a, 0xb5, 0xd3, 0xd6, 0x76, 0x4c, 0x61, 0xef, 0xb9, 0x93, 0x02, 0xbd, 0x5b}
	if id.Time() != 0x1563e3ab5d3 {
		t.Fatalf("time=%x", id.Time())
	}
}

func TestMarshalBinary(t *testing.T) {
	id := ULID{0x01, 0x56, 0x3e, 0x3a, 0xb5, 0xd3, 0xd6, 0x76, 0x4c, 0x61, 0xef, 0xb9, 0x93, 0x02, 0xbd, 0x5b}
	data, err := id.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x01, 0x56, 0x3e, 0x3a, 0xb5, 0xd3, 0xd6, 0x76, 0x4c, 0x61, 0xef, 0xb9, 0x93, 0x02, 0xbd, 0x5b}
	if !bytes.Equal(data, want) {
		t.Fatalf("data=%x", data)
	}
}

func TestUnmarshalBinary(t *testing.T) {
	data := []byte{0x01, 0x56, 0x3e, 0x3a, 0xb5, 0xd3, 0xd6, 0x76, 0x4c, 0x61, 0xef, 0xb9, 0x93, 0x02, 0xbd, 0x5b}
	var id ULID
	if err := id.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
	want := ULID{0x01, 0x56, 0x3e, 0x3a, 0xb5, 0xd3, 0xd6, 0x76, 0x4c, 0x61, 0xef, 0xb9, 0x93, 0x02, 0xbd, 0x5b}
	if id != want {
		t.Fatalf("want %v, got %v", want, id)
	}
}

func TestAppendBinary(t *testing.T) {
	id := ULID{0x01, 0x56, 0x3e, 0x3a, 0xb5, 0xd3, 0xd6, 0x76, 0x4c, 0x61, 0xef, 0xb9, 0x93, 0x02, 0xbd, 0x5b}
	data, err := id.AppendBinary(nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x01, 0x56, 0x3e, 0x3a, 0xb5, 0xd3, 0xd6, 0x76, 0x4c, 0x61, 0xef, 0xb9, 0x93, 0x02, 0xbd, 0x5b}
	if !bytes.Equal(data, want) {
		t.Fatalf("data=%x", data)
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
