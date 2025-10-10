package benchmark

import (
	"crypto/rand"
	"runtime"
	"testing"
	"time"

	oklog "github.com/oklog/ulid/v2"
	"github.com/shogo82148/go-ulid"
)

func BenchmarkOklogMake(b *testing.B) {
	for b.Loop() {
		runtime.KeepAlive(oklog.MustNew(oklog.Timestamp(time.Now()), rand.Reader))
	}
}

func BenchmarkShogoMake(b *testing.B) {
	for b.Loop() {
		runtime.KeepAlive(ulid.Make())
	}
}

func BenchmarkOklogString(b *testing.B) {
	id := oklog.Make()
	for b.Loop() {
		runtime.KeepAlive(id.String())
	}
}

func BenchmarkShogoString(b *testing.B) {
	id := ulid.Make()
	for b.Loop() {
		runtime.KeepAlive(id.String())
	}
}

func BenchmarkOklogParse(b *testing.B) {
	s := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	for b.Loop() {
		_, err := oklog.ParseStrict(s)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkShogoParse(b *testing.B) {
	s := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	for b.Loop() {
		_, err := ulid.Parse(s)
		if err != nil {
			b.Fatal(err)
		}
	}
}
