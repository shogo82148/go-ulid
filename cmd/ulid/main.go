package main

import (
	"fmt"
	"os"
	"time"

	"github.com/shogo82148/go-ulid"
)

const rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"

func main() {
	if len(os.Args) < 2 {
		// no argument: generate a new ULID
		id := ulid.Make()
		fmt.Println(id.String())
	} else {
		// one argument: parse the ULID and print the time component
		id, err := ulid.Parse(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		t := time.UnixMilli(id.Time())
		fmt.Println(t.Format(rfc3339Milli))
	}
}
