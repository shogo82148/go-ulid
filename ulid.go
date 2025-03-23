package ulid

import (
	"bytes"
	"crypto/rand"
	"errors"
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

// Zero is the zero value for a ULID.
var Zero ULID

// Errors returned by the Parse function.
var ErrInvalidSize = errors.New("ulid: invalid size")

// Errors returned by the Parse function.
var ErrInvalidCharacter = errors.New("ulid: invalid character")

// Errors returned by the Parse function.
var ErrOverflow = errors.New("ulid: overflow")

// EncodedSize is the size of a ULID when encoded to text.
const EncodedSize = 26

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
	id[0] = byte(ms >> 40)
	id[1] = byte(ms >> 32)
	id[2] = byte(ms >> 24)
	id[3] = byte(ms >> 16)
	id[4] = byte(ms >> 8)
	id[5] = byte(ms)
}

// Time returns the time component of the ULID as Unix milliseconds.
func (id ULID) Time() int64 {
	return int64(id[0])<<40 | int64(id[1])<<32 | int64(id[2])<<24 |
		int64(id[3])<<16 | int64(id[4])<<8 | int64(id[5])
}

// MarshalBinary implements the [encoding.BinaryMarshaler] interface.
func (id ULID) MarshalBinary() ([]byte, error) {
	ret := make([]byte, len(id))
	copy(ret, id[:])
	return ret, nil
}

// UnmarshalBinary implements the [encoding.BinaryUnmarshaler] interface.
func (id *ULID) UnmarshalBinary(data []byte) error {
	if len(data) != len(id) {
		return ErrInvalidSize
	}
	copy(id[:], data)
	return nil
}

// AppendBinary implements the [encoding.BinaryAppender] interface.
func (id ULID) AppendBinary(b []byte) ([]byte, error) {
	return append(b, id[:]...), nil
}

// Parse parses a ULID from a string.
func Parse(s string) (ULID, error) {
	return parse(s)
}

type bs interface {
	[]byte | string
}

func parse[T bs](s T) (ULID, error) {
	if len(s) != EncodedSize {
		return ULID{}, ErrInvalidSize
	}

	// Use an optimized unrolled loop to decode a base32 ULID.
	// The MSB(Most Significant Bit) is reserved for detecting invalid indexes.
	//
	// For example, in normal case, the bit layout of uint64(dec[v[0]])<<45 becomes:
	//
	//     | 63 | 62 | 61 | 60 | 59 | 58 | 57 | 56 | 55 | 54 | 53 | 52 | 51 | 50 | 49 | 48 | 47 | 46 | 45 | 44 | ... |
	//     |----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|-----|
	//     |  0 |  0 |  0 |  0 |  0 |  0 |  0 |  0 |  0 |  0 |  0 |  0 |  0 |  0 |  x |  x |  x |  x |  x |  0 | ... |
	//
	// and the MSB is set to 0.
	//
	// If the character is not part of the base32 character set, the layout becomes:
	//
	//     | 63 | 62 | 61 | 60 | 59 | 58 | 57 | 56 | 55 | 54 | 53 | 52 | 51 | 50 | 49 | 48 | 47 | 46 | 45 | 44 | ... |
	//     |----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|-----|
	//     |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  1 |  0 | ... |
	//
	// and the MSB is set to 1.
	h := uint64(dec[s[0]])<<45 |
		uint64(dec[s[1]])<<40 |
		uint64(dec[s[2]])<<35 |
		uint64(dec[s[3]])<<30 |
		uint64(dec[s[4]])<<25 |
		uint64(dec[s[5]])<<20 |
		uint64(dec[s[6]])<<15 |
		uint64(dec[s[7]])<<10 |
		uint64(dec[s[8]])<<5 |
		uint64(dec[s[9]])
	m := uint64(dec[s[10]])<<35 |
		uint64(dec[s[11]])<<30 |
		uint64(dec[s[12]])<<25 |
		uint64(dec[s[13]])<<20 |
		uint64(dec[s[14]])<<15 |
		uint64(dec[s[15]])<<10 |
		uint64(dec[s[16]])<<5 |
		uint64(dec[s[17]])
	l := uint64(dec[s[18]])<<35 |
		uint64(dec[s[19]])<<30 |
		uint64(dec[s[20]])<<25 |
		uint64(dec[s[21]])<<20 |
		uint64(dec[s[22]])<<15 |
		uint64(dec[s[23]])<<10 |
		uint64(dec[s[24]])<<5 |
		uint64(dec[s[25]])

	// Check if all the characters in a base32 encoded ULID are part of the
	// expected base32 character set.
	if (h|m|l)&(1<<63) != 0 {
		return ULID{}, ErrInvalidCharacter
	}

	if s[0] > '7' {
		return ULID{}, ErrOverflow
	}

	var id ULID

	// 6 bytes timestamp
	id[0] = byte(h >> 40)
	id[1] = byte(h >> 32)
	id[2] = byte(h >> 24)
	id[3] = byte(h >> 16)
	id[4] = byte(h >> 8)
	id[5] = byte(h)

	// 10 bytes random
	id[6] = byte(m >> 32)
	id[7] = byte(m >> 24)
	id[8] = byte(m >> 16)
	id[9] = byte(m >> 8)
	id[10] = byte(m)
	id[11] = byte(l >> 32)
	id[12] = byte(l >> 24)
	id[13] = byte(l >> 16)
	id[14] = byte(l >> 8)
	id[15] = byte(l)

	return id, nil
}

// We use -1 (all bits are set to 1) as sentinel value for invalid indexes.
// The reason for using -1 is that, even when cast, it does not lose the property that all bits are set to 1.
// e.g. uint64(int8(-1)) == 0xFFFFFFFFFFFFFFFF
var dec = [...]int8{
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, 0x00, 0x01,
	0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, -1, -1,
	-1, -1, -1, -1, -1, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E,
	0x0F, 0x10, 0x11, -1, 0x12, 0x13, -1, 0x14, 0x15, -1,
	0x16, 0x17, 0x18, 0x19, 0x1A, -1, 0x1B, 0x1C, 0x1D, 0x1E,
	0x1F, -1, -1, -1, -1, -1, -1, 0x0A, 0x0B, 0x0C,
	0x0D, 0x0E, 0x0F, 0x10, 0x11, -1, 0x12, 0x13, -1, 0x14,
	0x15, -1, 0x16, 0x17, 0x18, 0x19, 0x1A, -1, 0x1B, 0x1C,
	0x1D, 0x1E, 0x1F, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1,
}

const encoding = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

func (id ULID) text() [26]byte {
	var buf [26]byte

	// Optimized unrolled loop ahead.
	// Combining 32 bit loads allows the same code to be used
	// for 32 and 64 bit platforms.
	a := uint32(id[0])<<16 |
		uint32(id[1])<<8 |
		uint32(id[2])
	b := uint32(id[2])<<24 |
		uint32(id[3])<<16 |
		uint32(id[4])<<8 |
		uint32(id[5])
	c := uint32(id[6])<<24 |
		uint32(id[7])<<16 |
		uint32(id[8])<<8 |
		uint32(id[9])
	d := uint32(id[9])<<24 |
		uint32(id[10])<<16 |
		uint32(id[11])<<8 |
		uint32(id[12])
	e := uint32(id[12])<<24 |
		uint32(id[13])<<16 |
		uint32(id[14])<<8 |
		uint32(id[15])

	// 10 bytes timestamp
	buf[0] = encoding[(a>>21)&0x1f]
	buf[1] = encoding[(a>>16)&0x1f]
	buf[2] = encoding[(a>>11)&0x1f]
	buf[3] = encoding[(a>>6)&0x1f]
	buf[4] = encoding[(a>>1)&0x1f]
	buf[5] = encoding[(b>>20)&0x1f]
	buf[6] = encoding[(b>>15)&0x1f]
	buf[7] = encoding[(b>>10)&0x1f]
	buf[8] = encoding[(b>>5)&0x1f]
	buf[9] = encoding[b&0x1f]

	// 16 bytes random
	buf[10] = encoding[(c>>27)&0x1f]
	buf[11] = encoding[(c>>22)&0x1f]
	buf[12] = encoding[(c>>17)&0x1f]
	buf[13] = encoding[(c>>12)&0x1f]
	buf[14] = encoding[(c>>7)&0x1f]
	buf[15] = encoding[(c>>2)&0x1f]
	buf[16] = encoding[(d>>21)&0x1f]
	buf[17] = encoding[(d>>16)&0x1f]
	buf[18] = encoding[(d>>11)&0x1f]
	buf[19] = encoding[(d>>6)&0x1f]
	buf[20] = encoding[(d>>1)&0x1f]
	buf[21] = encoding[(e>>20)&0x1f]
	buf[22] = encoding[(e>>15)&0x1f]
	buf[23] = encoding[(e>>10)&0x1f]
	buf[24] = encoding[(e>>5)&0x1f]
	buf[25] = encoding[e&0x1f]

	return buf
}

func (id ULID) String() string {
	buf := id.text()
	return string(buf[:])
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (id ULID) MarshalText() ([]byte, error) {
	buf := id.text()
	return buf[:], nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (id *ULID) UnmarshalText(data []byte) error {
	id2, err := parse(data)
	if err != nil {
		return err
	}
	*id = id2
	return nil
}

// AppendText implements the [encoding.TextAppender] interface.
func (id ULID) AppendText(b []byte) ([]byte, error) {
	buf := id.text()
	return append(b, buf[:]...), nil
}

// IsZero returns true if the ULID is the zero value.
func (id ULID) IsZero() bool {
	return id == Zero
}

// Compare returns an integer comparing two ULIDs lexicographically.
// The result will be 0 if id==other, -1 if id < other, and +1 if id > other.
func (id ULID) Compare(other ULID) int {
	return bytes.Compare(id[:], other[:])
}
