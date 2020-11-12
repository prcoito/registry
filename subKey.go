package registry

import (
	"encoding/binary"
	"errors"
	"io"
)

type subKeyList struct {
	rws io.ReadWriteSeeker

	binOffset int64
	fpOffset  int64

	signature string // must be one of ["lf", "lh", "li", "ri"]

	numberElements uint16

	elements []*subKeyElement
}

func newSubKeyList(rws io.ReadWriteSeeker, binOffset, fpOffset int64) *subKeyList {
	return &subKeyList{
		rws:       rws,
		binOffset: binOffset,
		fpOffset:  fpOffset,
	}
}

func (skl *subKeyList) validate() error {
	if !(skl.signature == subKeyList1Sig ||
		skl.signature == subKeyList2Sig ||
		skl.signature == subKeyList3Sig ||
		skl.signature == subKeyList4Sig) {
		return errorW{err: ErrCorruptRegistry, cause: errBadSignature, function: "subKeyList.validate()"}
	}

	return nil
}

func (skl *subKeyList) Read() (err error) {
	r := skl.rws

	_, err = r.Seek(skl.fpOffset, io.SeekStart)
	if err != nil {
		return errorW{err: ErrCorruptRegistry, cause: err, function: "subKeyList.Read() r.Seek"}
	}

	b := make([]byte, 2)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return errorW{err: ErrCorruptRegistry, cause: err, function: "subKeyList.Read() io.ReadFull"}
	}
	skl.signature = string(b)

	_, err = io.ReadFull(r, b)
	if err != nil {
		return errorW{err: ErrCorruptRegistry, cause: err, function: "subKeyList.Read() io.ReadFull"}
	}

	skl.numberElements = binary.LittleEndian.Uint16(b)

	for i := uint16(0); i < skl.numberElements; i++ {
		el := newSubKeyElement(r, skl.binOffset, skl.binOffset, skl.signature)
		err = el.Read()
		if err != nil {
			return err
		}
		skl.elements = append(skl.elements, el)
	}
	return skl.validate()
}

func (skl *subKeyList) subkeyNames(n int) (names []string, err error) {
	max := int(skl.numberElements)
	if n < 0 {
		n = max
	}
	if n > max {
		n = max
	}

	names = make([]string, n)
	k := 0
	for i := 0; i < n; i++ {
		el := skl.elements[k]
		err := el.ReadElement()
		if err != nil {
			return nil, err
		}
		if el.namedKey != nil {
			names[i] = el.namedKey.name
		} else if el.subKeyList != nil {
			el.subKeyList.rws = skl.rws
			nms, err := el.subKeyList.subkeyNames(max - i)
			if err != nil {
				return nil, err
			}
			for _, v := range nms {
				names[i] = v
				i++
			}
		}
		k++
	}
	return names, nil
}

func (skl *subKeyList) allElements() (el []*subKeyElement, err error) {
	var els []*subKeyElement

	for _, e := range skl.elements {
		err = e.ReadElement()
		if err != nil {
			return
		}
		el = append(el, e)

		if e.subKeyList != nil {
			e.subKeyList.rws = skl.rws
			els, err = e.subKeyList.allElements()
			if err != nil {
				return
			}
			el = append(el, els...)
		}
	}
	return
}

type subKeyElement struct {
	rws            io.ReadWriteSeeker
	binOffset      int64
	hiveDataOffset int64

	signature string

	namedKeyOffset uint32 // set if lf, lh or li
	hashValue      uint32 // different than 0 if lf or lh
	namedKey       *namedKey

	subKeyListOffset uint32 // set if ri
	subKeyList       *subKeyList
}

func newSubKeyElement(rws io.ReadWriteSeeker, binOffset, dataOffset int64, sig string) *subKeyElement {
	return &subKeyElement{
		rws:            rws,
		binOffset:      binOffset,
		hiveDataOffset: dataOffset,
		signature:      sig,
	}
}

func (el *subKeyElement) Read() error {
	r := el.rws

	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return errorW{err: ErrCorruptRegistry, cause: err, function: "subKeyElement.Read() io.ReadFull"}
	}

	switch el.signature {
	case "lf", "lh":
		el.namedKeyOffset = binary.LittleEndian.Uint32(buf)
		_, err := io.ReadFull(r, buf)
		if err != nil {
			return errorW{err: ErrCorruptRegistry, cause: err, function: "subKeyElement.Read() io.ReadFull"}
		}
		el.hashValue = binary.LittleEndian.Uint32(buf)
	case "li":
		el.namedKeyOffset = binary.LittleEndian.Uint32(buf)
	case "ri":
		el.subKeyListOffset = binary.LittleEndian.Uint32(buf)
	}

	return nil
}

func (el *subKeyElement) ReadElement() error {
	switch el.signature {
	case "lf", "lh":
		el.namedKey = newNamedKey(
			el.rws,
			el.binOffset,
			el.hiveDataOffset+int64(el.namedKeyOffset),
		)
		err := el.namedKey.Read()
		if err != nil {
			return err
		}
		hash := lhSubKeyHash(el.namedKey.name)
		if hash != el.hashValue {
			return errorW{err: ErrCorruptRegistry, cause: errInvalidHash, function: "subKeyElement.ReadElement() hash comparision"}
		}
	case "ri":
		el.subKeyList = newSubKeyList(
			el.rws,
			el.binOffset,
			el.hiveDataOffset+int64(el.subKeyListOffset),
		)
		err := el.subKeyList.Read()
		if err != nil {
			return err
		}
	default:
		return errors.New("Unsupported element type " + el.signature)
	}

	return nil
}
