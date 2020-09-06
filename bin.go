package registry

import (
	"encoding/binary"
	"fmt"
	"io"
)

type bin struct {
	offset int64
	header *binHeader
	cell   *binCell
}

type binHeader struct {
	hiveOffset uint32
	hiveSize   uint32
	size       uint8
}

func (bh *binHeader) Validate(header []byte) error {
	if string(header[:4]) != binHeaderSig || len(header) != 32 {
		return ErrInvalidBinHeader
	}

	bh.hiveOffset = binary.LittleEndian.Uint32(header[4:8])
	bh.hiveSize = binary.LittleEndian.Uint32(header[8:12])

	return nil
}

type binCell struct {
	size int32
	data interface{}
}

func (c *binCell) String() string {
	return fmt.Sprintf("{size: %d; data: %+v}", c.size, c.data)
}

func (c *binCell) Read(r io.ReadSeeker) error {
	err := binary.Read(r, binary.LittleEndian, &c.size)

	var signature [2]byte
	_, err = r.Read(signature[:])
	if err != nil {
		return err
	}

	_, err = r.Seek(-2, 1)
	if err != nil {
		return err
	}

	offset, _ := r.Seek(0, 1)
	switch string(signature[:]) {
	case "nk":
		nk := &namedKey{fpOffset: offset}
		c.data = nk
		err = nk.Read(r)
	}

	return err
}

func getHiveBins(fp io.ReadSeeker) ([]bin, error) {
	_, err := fp.Seek(4096, 0) // goto end of registry header
	if err != nil {
		return nil, err
	}

	// var size uint32 = 1
	bins := make([]bin, 0)

	// for size == 1 {
	offset, _ := fp.Seek(0, 1)
	b := bin{header: &binHeader{size: 4}, cell: &binCell{}, offset: offset}
	header := make([]byte, 32)
	_, err = fp.Read(header)
	if err != nil {
		return nil, err
	}

	err = b.header.Validate(header)
	if err != nil {
		return nil, err
	}

	err = b.cell.Read(fp)
	if err != nil {
		return nil, err
	}

	bins = append(bins, b)

	// }

	return bins, nil
}
