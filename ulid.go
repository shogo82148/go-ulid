package ulid

import (
	"crypto/rand"
	"time"
)

/*
A ULID is a 16 byte Universally Unique Lexicographically Sortable Identifier

	The components are encoded as 16 octets.
	Each component is encoded with the MSB first (network byte order).

	0                   1                   2                   3
	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|                      32_bit_uint_time_high                    |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|     16_bit_uint_time_low      |       16_bit_uint_random      |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|                       32_bit_uint_random                      |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|                       32_bit_uint_random                      |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/
type ULID [16]byte

// Make returns a ULID with the current time in Unix milliseconds and a random component.
func Make() ULID {
	id := ULID{}
	id.SetTime(time.Now().UnixMilli())
	if _, err := rand.Read(id[6:]); err != nil {
		panic(err)
	}
	return id
}

// SetTime sets the time component of the ULID to the given Unix milliseconds.
func (id *ULID) SetTime(ms int64) {
	(*id)[0] = byte(ms >> 40)
	(*id)[1] = byte(ms >> 32)
	(*id)[2] = byte(ms >> 24)
	(*id)[3] = byte(ms >> 16)
	(*id)[4] = byte(ms >> 8)
	(*id)[5] = byte(ms)
}

const encoding = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

func (id ULID) text() [26]byte {
	var buf [26]byte

	// Optimized unrolled loop ahead.
	h := uint64(id[0])<<56 | uint64(id[1])<<48 | uint64(id[2])<<40 |
		uint64(id[3])<<32 | uint64(id[4])<<24 | uint64(id[5])<<16 |
		uint64(id[6])<<8 | uint64(id[7])
	l := uint64(id[8])<<56 | uint64(id[9])<<48 | uint64(id[10])<<40 |
		uint64(id[11])<<32 | uint64(id[12])<<24 | uint64(id[13])<<16 |
		uint64(id[14])<<8 | uint64(id[15])

	// 10 bytes timestamp
	buf[0] = encoding[(h>>61)&0x1f]
	buf[1] = encoding[(h>>56)&0x1f]
	buf[2] = encoding[(h>>51)&0x1f]
	buf[3] = encoding[(h>>46)&0x1f]
	buf[4] = encoding[(h>>41)&0x1f]
	buf[5] = encoding[(h>>36)&0x1f]
	buf[6] = encoding[(h>>31)&0x1f]
	buf[7] = encoding[(h>>26)&0x1f]
	buf[8] = encoding[(h>>21)&0x1f]
	buf[9] = encoding[(h>>16)&0x1f]

	// 16 bytes random
	buf[10] = encoding[(h>>11)&0x1f]
	buf[11] = encoding[(h>>6)&0x1f]
	buf[12] = encoding[(h>>1)&0x1f]
	buf[13] = encoding[(h<<4|l>>60)&0x1f]
	buf[14] = encoding[(l>>55)&0x1f]
	buf[15] = encoding[(l>>50)&0x1f]
	buf[16] = encoding[(l>>45)&0x1f]
	buf[17] = encoding[(l>>40)&0x1f]
	buf[18] = encoding[(l>>35)&0x1f]
	buf[19] = encoding[(l>>30)&0x1f]
	buf[20] = encoding[(l>>25)&0x1f]
	buf[21] = encoding[(l>>20)&0x1f]
	buf[22] = encoding[(l>>15)&0x1f]
	buf[23] = encoding[(l>>10)&0x1f]
	buf[24] = encoding[(l>>5)&0x1f]
	buf[25] = encoding[l&0x1f]

	return buf
}

func (id ULID) String() string {
	buf := id.text()
	return string(buf[:])
}

func (id ULID) MarshalText() ([]byte, error) {
	buf := id.text()
	return buf[:], nil
}

func (id ULID) AppendText(b []byte) ([]byte, error) {
	buf := id.text()
	return append(b, buf[:]...), nil
}
