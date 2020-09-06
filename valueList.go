package registry

import (
	"encoding/binary"
	"errors"
	"io"
)

type valueList struct {
	rws io.ReadWriteSeeker

	binOffset        int64
	valuesListOffset uint32
	numberOfValues   uint32

	offsets []uint32 // Value key offset. The offset value is in bytes and relative from the start of the hive bin data
	values  []*valueKey
}

func newValueList(rws io.ReadWriteSeeker, binOffset int64,
	valuesListOffset, numberOfValues uint32) *valueList {
	return &valueList{
		rws:              rws,
		binOffset:        binOffset,
		numberOfValues:   numberOfValues,
		valuesListOffset: valuesListOffset,
		offsets:          make([]uint32, numberOfValues),
		values:           make([]*valueKey, numberOfValues),
	}
}

// Read reads offsets for valueList and allocates space for values
func (vl *valueList) Read() error {
	r := vl.rws

	_, err := r.Seek(vl.binOffset+int64(vl.valuesListOffset), 0)
	if err != nil {
		return err
	}

	b := make([]byte, 4)
	for i := uint32(0); i < vl.numberOfValues; i++ {
		_, err = r.Read(b)
		if err != nil {
			return err
		}

		vl.offsets[i] = binary.LittleEndian.Uint32(b)
	}

	// vl.values = make([]*valueKey, len(vl.offsets))
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

	vk := newValueKey(vl.rws, vl.binOffset, vl.offsets[i])
	return vk, vk.Read()
}
