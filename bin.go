package registry

import (
	"encoding/binary"
	"fmt"
	"io"
)

type bin struct {
	rws io.ReadWriteSeeker

	offset int64
	header *binHeader
	cell   *binCell
}

func newBin(r io.ReadWriteSeeker, offset int64) bin {
	return bin{
		rws:    r,
		header: &binHeader{size: 4},
		cell:   &binCell{rws: r, binOffset: offset},
		offset: offset,
	}
}

type binHeader struct {
	hiveOffset uint32
	hiveSize   uint32
	size       uint8
}

func (bh *binHeader) Validate(header []byte) error {
	if string(header[:4]) != binHeaderSig {
		return errBadSignature
	}
	if len(header) != 32 {
		return errInvalidBinHeader
	}

	bh.hiveOffset = binary.LittleEndian.Uint32(header[4:8])
	bh.hiveSize = binary.LittleEndian.Uint32(header[8:12])

	return nil
}

type binCell struct {
	binOffset int64
	rws       io.ReadWriteSeeker

	size int32
	data interface{}
}

func (c *binCell) Read() error {
	rws := c.rws

	err := binary.Read(rws, binary.LittleEndian, &c.size)
	if err != nil {
		return err
	}

	var signature [2]byte
	_, err = io.ReadFull(rws, signature[:])
	if err != nil {
		return err
	}

	_, err = rws.Seek(-2, io.SeekCurrent)
	if err != nil {
		return err
	}

	offset, _ := rws.Seek(0, io.SeekCurrent)
	sig := string(signature[:])
	switch sig {
	case "nk":
		nk := newNamedKey(rws, int64(c.binOffset), offset)
		c.data = nk
		err = nk.Read()
	default:
		return fmt.Errorf("Cell with %v not supported yet", sig)
	}

	return err
}

// TODO: this only gets the first bin
// Either create a new "getRootBin" function or fix this
func getHiveBins(rs io.ReadWriteSeeker) ([]bin, error) {
	_, err := rs.Seek(4096, io.SeekStart) // goto end of registry header
	if err != nil {
		return nil, err
	}

	// var size uint32 = 1
	bins := make([]bin, 0)

	// for size == 1 {
	offset, _ := rs.Seek(0, io.SeekCurrent)
	b := newBin(rs, offset)
	header := make([]byte, 32)
	_, err = io.ReadFull(rs, header)
	if err != nil {
		return nil, err
	}

	err = b.header.Validate(header)
	if err != nil {
		return nil, err
	}

	err = b.cell.Read()
	if err != nil {
		return nil, err
	}

	bins = append(bins, b)

	// }

	return bins, nil
}
