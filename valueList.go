package registry

import (
	"encoding/binary"
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

// newValueList creates a valueList
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

// Read reads offsets for valueList
func (vl *valueList) Read() error {
	r := vl.rws

	_, err := r.Seek(vl.binOffset+int64(vl.valuesListOffset), io.SeekStart)
	if err != nil {
		return errorW{err: ErrCorruptRegistry, cause: err, function: "valueList.Read() r.Seek"}
	}

	b := make([]byte, 4)
	for i := uint32(0); i < vl.numberOfValues; i++ {
		_, err := io.ReadFull(r, b)
		if err != nil {
			return errorW{err: ErrCorruptRegistry, cause: err, function: "valueList.Read() io.ReadFull"}
		}

		vl.offsets[i] = binary.LittleEndian.Uint32(b)
	}

	return nil
}

func (vl *valueList) Len() int {
	return int(vl.numberOfValues)
}

// Value returns the i th value from value list
// If i th value was already read it is immediately returned, otherwise it is called ReadValue
func (vl *valueList) Value(i uint) (*valueKey, error) {
	if uint32(i) >= vl.numberOfValues {
		return nil, ErrOutOfBounds
	}

	var err error
	if vl.values[i] == nil {
		vl.values[i], err = vl.ReadValue(i)
	}
	return vl.values[i], err
}

// ReadValues returns the i th values from value list.
// If i is bigger than the number of values on list ErrOutOfBounds is returned
// If a error occurs reading the value, ErrOutOfBounds or ErrCorruptRegistry is returned
func (vl *valueList) ReadValue(i uint) (*valueKey, error) {
	if uint32(i) >= vl.numberOfValues {
		return nil, ErrOutOfBounds
	}

	vk := newValueKey(vl.rws, vl.binOffset, vl.offsets[i])
	return vk, vk.Read()
}
