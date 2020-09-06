package registry

import (
	"encoding/binary"
	"errors"
)

type valueList struct {
	nk *namedKey

	offsets []uint32 // Value key offset. The offset value is in bytes and relative from the start of the hive bin data
	values  []*valueKey
}

// Read reads offsets for valueList and allocates space for values
func (vl *valueList) Read() error {
	r := vl.nk.readSeeker
	_, err := r.Seek(vl.nk.binOffset+int64(vl.nk.valuesListOffset), 0)
	if err != nil {
		return err
	}
	b := make([]byte, 4)
	for i := uint32(0); i < vl.nk.numberOfValues; i++ {
		r.Read(b)
		offset := binary.LittleEndian.Uint32(b)
		vl.offsets = append(vl.offsets, offset)
	}

	vl.values = make([]*valueKey, len(vl.offsets))
	return nil
}

func (vl *valueList) Len() int {
	return len(vl.offsets)
}

func (vl *valueList) Value(i uint) (*valueKey, error) {
	if int(i) >= len(vl.offsets) {
		return nil, errors.New("index out of bounds")
	}

	var err error
	if vl.values[i] == nil {
		vl.values[i], err = vl.ReadValue(i)
	}
	return vl.values[i], err
}

func (vl *valueList) ReadValue(i uint) (*valueKey, error) {
	if int(i) >= len(vl.offsets) {
		return nil, errors.New("index out of bounds")
	}

	r := vl.nk.readSeeker
	r.Seek(vl.nk.binOffset+int64(vl.offsets[i]), 0)
	vk := &valueKey{binOffset: vl.nk.binOffset, valueOffset: vl.offsets[i]}
	return vk, vk.Read(r)
}

type valueData struct {
	signature string // must be "db"

	numberSegments uint16
	dbOffset       uint16 // Data block (segment) list offset. The offset value is in bytes and relative from the start of the hive bin data.

	padding uint16 // due to 8 byte alignment of cell size. Sometimes contains remnant data
}

func (vd valueData) validate() error {
	if vd.signature != dataBlockSig {
		return ErrBadSignature
	}

	return nil
}

type dataBlockSegmentList struct {
	entries []uint16
}
