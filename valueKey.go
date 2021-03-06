package registry

import (
	"encoding/binary"
	"fmt"
	"io"
)

type valueKey struct {
	rws io.ReadWriteSeeker

	binOffset   int64  // hive bin offset
	valueOffset uint32 // offset of the value relative to binOffset

	signature string // must be equal to "vk"

	nameSize uint16
	dataSize uint32 /* 0 means not set / NULL

	If the MSB 0x80000000 of the data size is set the data offset actually contains the data value.
	A data size of 4 uses all 4 bytes of the data offset
	A data size of 2 uses the last 2 bytes of the data offset (on a little-endian system)
	A data size of 1 uses the last byte (on a little-endian system)
	A data size of 0 represents that the value is not set (or NULL).
	The behavior on a big-endian system is unknown.*/

	dataOffset uint32
	dataType   uint32

	flags uint16

	name string

	data interface{}
}

func newValueKey(rws io.ReadWriteSeeker, binOffset int64, valueOffset uint32) *valueKey {
	return &valueKey{
		binOffset:   binOffset,
		valueOffset: valueOffset,
		rws:         rws,
	}
}

func (vk *valueKey) Read() error {
	r := vk.rws

	_, err := r.Seek(vk.binOffset+int64(vk.valueOffset), io.SeekStart)
	if err != nil {
		return errorW{err: ErrCorruptRegistry, cause: err, function: "valueKey.Read() r.Seek"}
	}
	b := make([]byte, 20)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return errorW{err: ErrCorruptRegistry, cause: err, function: "valueKey.Read() io.ReadFull"}
	}

	vk.signature = string(b[:2])
	vk.nameSize = binary.LittleEndian.Uint16(b[2:4])
	vk.dataSize = binary.LittleEndian.Uint32(b[4:8])
	vk.dataOffset = binary.LittleEndian.Uint32(b[8:12])
	vk.dataType = binary.LittleEndian.Uint32(b[12:16])
	vk.flags = binary.LittleEndian.Uint16(b[16:18])
	// b[18:20] = unknown
	dataSize := b[4:8]
	dataOffset := b[8:12]

	if vk.nameSize == 0 {
		vk.name = "(default)"
	} else {
		b = make([]byte, vk.nameSize)
		_, err := io.ReadFull(r, b)
		if err != nil {
			return errorW{err: ErrCorruptRegistry, cause: err, function: "valueKey.Read() io.ReadFull"}
		}
		vk.name = string(b)
	}

	// If the MSB of the data size is set the data offset actually contains the data value.
	if (dataSize[3]>>7)&1 == 1 {
		b = dataOffset[:]

		switch dataSize[0] {
		case 0:
			vk.data = nil
			vk.dataSize = 0
		case 1:
			vk.data = b[3:]
			vk.dataSize = 1
		case 2:
			vk.data = b[2:]
			vk.dataSize = 2
		default:
			vk.data = b[:]
			vk.dataSize = 4
		}
	} else {
		_, err = r.Seek(vk.binOffset+int64(vk.dataOffset), io.SeekStart)
		if err != nil {
			return errorW{err: ErrCorruptRegistry, cause: err, function: "valueKey.Read() r.Seek"}
		}
		b = make([]byte, vk.dataSize)
		_, err := io.ReadFull(r, b)
		if err != nil {
			return errorW{err: ErrCorruptRegistry, cause: err, function: "valueKey.Read() io.ReadFull"}
		}
		vk.data = b[:]
	}

	if vk.data != nil {
		switch vk.dataType {
		case REG_BINARY: // already []byte
		case REG_SZ, REG_EXPAND_SZ:
			vk.data = stringFromBytes(vk.data.([]byte))
			vk.dataSize = (vk.dataSize - 1) / 2 // 2 byte char to 1 byte char excluding \0
		case REG_DWORD, REG_QWORD:
			vk.data = uint64FromBytesLE(vk.data.([]byte))
		case REG_DWORD_BIG_ENDIAN:
			vk.data = uint32FromBytesBE(vk.data.([]byte))
		case REG_MULTI_SZ:
			vk.data = stringsFromBytes(vk.data.([]byte))
			vk.dataSize = (vk.dataSize - 1) / 2 // 2 byte char to 1 byte char excluding \0
		default:
			return fmt.Errorf("Data type %v not supported yet", Type(vk.dataType))
		}
	}

	return vk.validate()
}

func (vk *valueKey) validate() error {
	if vk.signature != valueKeySig {
		return errorW{err: ErrCorruptRegistry, cause: errBadSignature, function: "valueKey.validate()"}
	}

	return nil
}
