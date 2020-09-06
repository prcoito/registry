package registry

import (
	"encoding/binary"
	"fmt"
	"time"
	"unicode"
	"unicode/utf16"
)

func date(i uint64) time.Time {
	return time.Unix(0, int64(i)*int64(time.Nanosecond))
}

func stringFromBytes(u []byte) string {
	b := make([]uint16, len(u)/2)
	for i := 0; i < len(u); i += 2 {
		b[i/2] = (uint16(u[i+1]) << 8) + uint16(u[i])
	}
	if b[len(u)/2-1] == 0 {
		b = b[:len(u)/2-1]
	}
	return string(utf16.Decode(b))
}

func stringsFromBytes(u []byte) (r []string) {
	str := make([]uint16, 0)
	for i := 0; i < len(u); i += 2 {
		c := (uint16(u[i+1]) << 8) + uint16(u[i]) // utf16-LE to rune
		if c == 0 && len(str) > 0 {               // end of string
			r = append(r, string(utf16.Decode(str)))
			str = make([]uint16, 0)
		} else { // append to cur string
			str = append(str, c)
		}
	}
	return
}

func uint64FromBytesLE(u []byte) uint64 {
	b := make([]byte, 8) // make sure b is uint64
	copy(b, u)
	return binary.LittleEndian.Uint64(b)
}

func uint32FromBytesBE(u []byte) uint32 {
	b := make([]byte, 4) // make sure b is uint32
	copy(b, u)
	return binary.BigEndian.Uint32(b)
}

func bytesFromUint64LE(u uint64, dataType uint32) ([]byte, error) {
	switch dataType {
	case REG_DWORD:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(u))
		return b, nil
	case REG_QWORD:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, u)
		return b, nil
	default:
		return nil, fmt.Errorf("Invalid data type %v", Type(dataType))
	}
}

// bytesFromUint32BE returns the []byte representation of uint32
// ONLY USE FOR REG_DWORD_BIG_ENDIAN
func bytesFromUint32BE(u uint32) ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, u)
	return b, nil
}

func bytesFromStrings(strs []string) []byte {
	r := make([]byte, 0)
	for _, s := range strs {
		r = append(r, []byte(s)...)
		r = append(r, 0)
	}
	return r
}

func dataSizeFromType(u uint32) int {
	switch u {
	case REG_DWORD_LITTLE_ENDIAN, REG_DWORD_BIG_ENDIAN:
		return 4 // 32 bit
	case REG_QWORD:
		return 8 // 64 bit
	default:
		return 0
	}
}

func lhSubKeyHash(str string) uint32 {
	var hashValue uint32 = 0
	for idx := 0; idx < len(str); idx++ {
		hashValue *= 37
		hashValue += uint32(unicode.ToUpper(rune(str[idx])))
	}
	return hashValue
}
