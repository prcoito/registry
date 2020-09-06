package registry

import (
	"encoding/binary"
	"errors"
	"io"
)

type subKeyList struct {
	readSeeker io.ReadSeeker

	binOffset int64

	signature string // must be one of ["lf", "lh", "li", "ri"]

	numberElements uint16

	elements []*subKeyElement
}

func (skl *subKeyList) validate() error {
	if !(skl.signature == subKeyList1Sig ||
		skl.signature == subKeyList2Sig ||
		skl.signature == subKeyList3Sig ||
		skl.signature == subKeyList4Sig) {
		return ErrBadSignature
	}

	return nil
}

func (skl *subKeyList) Read(r io.ReadSeeker) (err error) {
	b := make([]byte, 2)
	_, err = r.Read(b)
	if err != nil {
		return err
	}
	skl.signature = string(b)

	_, err = r.Read(b)
	if err != nil {
		return err
	}

	skl.numberElements = binary.LittleEndian.Uint16(b)

	for i := uint16(0); i < skl.numberElements; i++ {
		el := newSubKeyElement(skl.binOffset, skl.binOffset, skl.signature)
		err = el.Read(r)
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
		err := el.ReadElement(skl.readSeeker)
		if err != nil {
			return nil, err
		}
		if el.namedKey != nil {
			names[i] = el.namedKey.name
		} else if el.subKeyList != nil {
			el.subKeyList.readSeeker = skl.readSeeker
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

func (skl *subKeyList) allElements() (el []*subKeyElement) {
	for _, e := range skl.elements {
		e.ReadElement(skl.readSeeker)
		el = append(el, e)
		if e.subKeyList != nil {
			e.subKeyList.readSeeker = skl.readSeeker
			el = append(el, e.subKeyList.allElements()...)
		}
	}
	return
}

type subKeyElement struct {
	binOffset      int64
	hiveDataOffset int64

	signature string

	namedKeyOffset uint32 // set if lf, lh or li
	hashValue      uint32 // different than 0 if lf or lh
	namedKey       *namedKey

	subKeyListOffset uint32 // set if ri
	subKeyList       *subKeyList
}

func newSubKeyElement(binOffset, dataOffset int64, sig string) *subKeyElement {
	return &subKeyElement{
		binOffset:      binOffset,
		hiveDataOffset: dataOffset,
		signature:      sig,
	}
}

func (el *subKeyElement) Read(r io.ReadSeeker) error {
	buf := make([]byte, 4)
	r.Read(buf)
	switch el.signature {
	case "lf", "lh":
		el.namedKeyOffset = binary.LittleEndian.Uint32(buf)
		r.Read(buf)
		el.hashValue = binary.LittleEndian.Uint32(buf)
	case "li":
		el.namedKeyOffset = binary.LittleEndian.Uint32(buf)
	case "ri":
		el.subKeyListOffset = binary.LittleEndian.Uint32(buf)
	}

	// fmt.Printf("subKeyElement: %+v\n", el)
	return nil
}

func (el *subKeyElement) ReadElement(r io.ReadSeeker) error {
	switch el.signature {
	case "lf", "lh":
		r.Seek(el.hiveDataOffset+int64(el.namedKeyOffset), 0)
		el.namedKey = &namedKey{}
		err := el.namedKey.Read(r)
		if err != nil {
			return err
		}
		hash := lhSubKeyHash(el.namedKey.name)
		if hash != el.hashValue {
			return errors.New("Element hash invalid")
		}
	case "ri":
		r.Seek(el.hiveDataOffset+int64(el.subKeyListOffset), 0)
		el.subKeyList = &subKeyList{binOffset: el.binOffset}
		err := el.subKeyList.Read(r)
		if err != nil {
			return err
		}
	default:
		return errors.New("Unsupported element type " + el.signature)
	}

	return nil
}
