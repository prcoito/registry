package registry

import (
	"bytes"
	"encoding/binary"
	"io"
)

// file header
type header struct {
	rws io.ReadWriteSeeker

	buf []byte

	lastModification uint64

	major uint32
	minor uint32

	fileType uint32 /* 0x0000 is normal file
	0x0001 is transaction log*/

	rootOffset uint32

	binSize uint32

	xor []byte
}

func newHeader(rws io.ReadWriteSeeker) *header {
	return &header{
		rws: rws,
	}
}

func (h *header) Read() error {
	h.buf = make([]byte, 4096)

	_, err := io.ReadFull(h.rws, h.buf)
	if err != nil {
		return errorW{err: ErrCorruptRegistry, cause: err, function: "header.Read() io.ReadFull"}
	}

	h.lastModification = binary.LittleEndian.Uint64(h.buf[12:20])

	// header versions
	h.major = binary.LittleEndian.Uint32(h.buf[20:24])
	h.minor = binary.LittleEndian.Uint32(h.buf[24:28])

	h.fileType = binary.LittleEndian.Uint32(h.buf[28:32])

	h.rootOffset = binary.LittleEndian.Uint32(h.buf[36:40])

	h.binSize = binary.LittleEndian.Uint32(h.buf[40:44])

	h.xor = h.buf[508:512]

	// h.buf[512:3576] = reserved

	return h.validate()
}

// validate reads header and validates it
func (h *header) validate() error {
	// header magic number
	if string(h.buf[:4]) != registrySig {
		return errBadSignature
	}

	if binary.LittleEndian.Uint32(h.buf[4:8]) != binary.LittleEndian.Uint32(h.buf[8:12]) {
		return errBadSequenceNumber
	}

	// calculate xor from previous bytes
	calculatedXOR := make([]byte, 4)
	for i, b := range h.buf[:508] {
		calculatedXOR[i&3] ^= b
	}

	if bytes.Compare(calculatedXOR, h.xor) != 0 {
		return errInvalidXOR
	}

	return nil
}
