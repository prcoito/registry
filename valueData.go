package registry

type valueData struct {
	signature string // must be "db"

	numberSegments uint16
	dbOffset       uint16 // Data block (segment) list offset. The offset value is in bytes and relative from the start of the hive bin data.

	padding uint16 // due to 8 byte alignment of cell size. Sometimes contains remnant data
}

func (vd valueData) validate() error {
	if vd.signature != dataBlockSig {
		return errorW{err: ErrCorruptRegistry, cause: errBadSignature, function: "valueData.validate()"}
	}

	return nil
}

type dataBlockSegmentList struct {
	entries []uint16
}
