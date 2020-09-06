package registry

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Key struct
type Key struct {
	nk *namedKey

	registry Registry
}

// OpenSubKey opens the subkey located at path
func (k Key) OpenSubKey(path string) (Key, error) {
	if k.nk.numberOfSubKeys == 0 {
		return Key{}, fmt.Errorf("Key %s does not have subkeys", k.nk.name)
	}

	return k.openSubKey(strings.Split(path, string(separator)))
}

func (k Key) openSubKey(entries []string) (Key, error) {
	// recursive
	// calculate hash "entries[0]"
	// check sublist hash
	// recurse with new key
	// if not found return error
	if len(entries) == 0 {
		return k, nil
	}

	list, err := k.subkeys()
	if err != nil {
		return Key{}, err
	}
	hash := lhSubKeyHash(entries[0])
	for _, sk := range list.allElements() {
		if sk.hashValue == hash {
			newKey := Key{registry: k.registry, nk: &namedKey{binOffset: k.nk.binOffset}}
			k.nk.readSeeker.Seek(list.binOffset+int64(sk.namedKeyOffset), 0)
			newKey.nk.Read(k.nk.readSeeker)
			return newKey.openSubKey(entries[1:])
		}
	}
	return Key{}, ErrNotExist
}

// Close closes open key k.
func (k Key) Close() error {
	// if this key was created by OpenKey function then
	// we must close the registry to properly close the
	// file
	if k.registry.createdByOpenKey {
		return k.registry.Close()
	}
	return nil
}

// GetBinaryValue retrieves the binary value for the specified
// value name associated with an open key k. It also returns the value's type.
// If value does not exist, GetBinaryValue returns ErrNotExist.
// If value is not REG_BINARY, it will return the correct value
// type and ErrUnexpectedType.
func (k Key) GetBinaryValue(name string) (val []byte, valtype uint32, err error) {
	var ok bool
	var value *valueKey

	value, err = k.getValue(name)
	if err != nil {
		return
	}

	valtype = value.dataType
	if valtype != REG_BINARY {
		err = ErrUnexpectedType
		return
	}

	val, ok = value.data.([]byte)
	if !ok {
		err = errors.New("Internal error: value.data is not binary")
	}
	return
}

// GetIntegerValue retrieves the integer value for the specified
// value name associated with an open key k. It also returns the value's type.
// If value does not exist, GetIntegerValue returns ErrNotExist.
// If value is not REG_DWORD or REG_QWORD, it will return the correct value
// type and ErrUnexpectedType.
func (k Key) GetIntegerValue(name string) (val uint64, valtype uint32, err error) {
	var ok bool
	var value *valueKey
	var u uint32

	value, err = k.getValue(name)
	if err != nil {
		return
	}

	valtype = value.dataType
	if valtype != REG_DWORD && // valtype != REG_DWORD_LITTLE_ENDIAN &&
		valtype != REG_DWORD_BIG_ENDIAN &&
		valtype != REG_QWORD { // && valtype != REG_QWORD_LITTLE_ENDIAN
		err = ErrUnexpectedType
		return
	}

	if valtype == REG_DWORD_BIG_ENDIAN {
		u, ok = value.data.(uint32)
		val = uint64(u)
	} else {
		val, ok = value.data.(uint64)
	}

	if !ok {
		err = errors.New("Internal error: value.data is not uint(64|32)")
	}
	return
}

// GetMUIStringValue retrieves the localized string value for
// the specified value name associated with an open key k.
// If the value name doesn't exist or the localized string value
// can't be resolved, GetMUIStringValue returns ErrNotExist.
// GetMUIStringValue panics if the system doesn't support
// regLoadMUIString; use LoadRegLoadMUIString to check if
// regLoadMUIString is supported before calling this function.
func (k Key) GetMUIStringValue(name string) (string, error) {
	return "", errors.New("Unsupported function")
}

// GetStringValue retrieves the string value for the specified
// value name associated with an open key k. It also returns the value's type.
// If value does not exist, GetStringValue returns ErrNotExist.
// If value is not REG_SZ or REG_EXPAND_SZ, it will return the correct value
// type and ErrUnexpectedType.
func (k Key) GetStringValue(name string) (val string, valtype uint32, err error) {
	var ok bool
	var value *valueKey

	value, err = k.getValue(name)
	if err != nil {
		return
	}
	valtype = value.dataType
	if valtype != REG_SZ && valtype != REG_EXPAND_SZ {
		err = ErrUnexpectedType
		return
	}

	val, ok = value.data.(string)
	if !ok {
		err = errors.New("Internal error: value.data is not string")
	}
	return
}

// GetStringsValue retrieves the []string value for the specified
// value name associated with an open key k. It also returns the value's type.
// If value does not exist, GetStringsValue returns ErrNotExist.
// If value is not REG_MULTI_SZ, it will return the correct value
// type and ErrUnexpectedType.
func (k Key) GetStringsValue(name string) (val []string, valtype uint32, err error) {
	var ok bool
	var value *valueKey

	value, err = k.getValue(name)
	if err != nil {
		return nil, 0, err
	}
	valtype = value.dataType
	if valtype != REG_MULTI_SZ {
		return nil, value.dataType, ErrUnexpectedType
	}

	val, ok = value.data.([]string)
	if !ok {
		return nil, 0, errors.New("Internal error: value.data is not []string")
	}
	return
}

// GetValue retrieves the type and data for the specified value associated
// with an open key k. It fills up buffer buf and returns the retrieved
// byte count n.
// If buf is too small to fit the stored value it returns
// ErrShortBuffer error along with the required buffer size n (no data copied).
// If no buffer is provided, it returns the actual buffer size n.
// If no buffer is provided, GetValue returns the value's type only.
// If the value does not exist, the error returned is ErrNotExist.
//
// GetValue is a low level function. If value's type is known, use the appropriate
// Get*Value function instead.
func (k Key) GetValue(name string, buf []byte) (n int, valtype uint32, err error) {
	v, err := k.getValue(name)
	if err != nil {
		return 0, 0, err
	}
	n = int(v.dataSize)
	valtype = v.dataType
	if buf == nil {
		return
	}
	if n > len(buf) {
		err = ErrShortBuffer
		return
	}
	switch v.data.(type) {
	case string, []byte:
		// REVIEW: should string be converted to UTF-16LE ?
		n = copy(buf, v.data.([]byte))
	case uint64:
		var b []byte
		if v.dataType == REG_DWORD_BIG_ENDIAN {
			b, err = bytesFromUint32BE(v.data.(uint32))
		} else {
			b, err = bytesFromUint64LE(v.data.(uint64), v.dataType)
		}
		n = copy(buf, b)
	case []string:
		// REVIEW: should string be converted to UTF-16LE ?
		n = copy(buf, bytesFromStrings(v.data.([]string)))
	}
	return
}

func (k Key) getValue(name string) (*valueKey, error) {
	list := k.nk.values
	for i := 0; i < list.Len(); i++ {
		value, err := list.Value(uint(i))
		if err != nil || value.name == name {
			return value, err
		}
	}
	return nil, ErrNotExist
}

// ReadSubKeyNames returns the names of subkeys of key k.
// The parameter n controls the number of returned names,
// analogous to the way os.File.Readdirnames works.
func (k Key) ReadSubKeyNames(n int) ([]string, error) {
	if n == 0 || k.nk.numberOfSubKeys == 0 { // TODO: check behavior on microsoft api
		return []string{}, nil
	}

	list, err := k.subkeys()
	if err != nil {
		return nil, err
	}

	max := int(k.nk.numberOfSubKeys)
	if n < 0 {
		n = max
	}
	if n > max {
		n = max
	}

	names := make([]string, n)
	j := 0
	for i := 0; i < n; i++ {
		el := list.elements[j]
		err := el.ReadElement(k.nk.readSeeker)
		if err != nil {
			return nil, err
		}
		if el.namedKey != nil { // current element is a named key
			names[i] = el.namedKey.name
		} else if el.subKeyList != nil { // current element is a subkey list
			el.subKeyList.readSeeker = k.nk.readSeeker
			n, err := el.subKeyList.subkeyNames(n - i)
			if err != nil {
				return nil, err
			}
			for _, v := range n { // add received subkeys
				names[i] = v
				i++
			}
		}
		j++
	}
	sort.Strings(names)
	return names, nil
}

func (k Key) subkeys() (*subKeyList, error) {
	if k.registry.header == nil {
		return nil, fmt.Errorf("Nil pointers")
	}
	fp := k.registry.header.fp

	fp.Seek(k.nk.binOffset+int64(k.nk.subKeysListOffset), 0)

	list := &subKeyList{binOffset: k.nk.binOffset, readSeeker: k.nk.readSeeker}
	return list, list.Read(fp)
}

// ReadValueNames returns the value names of key k.
// The parameter n controls the number of returned names,
// analogous to the way os.File.Readdirnames works.
func (k Key) ReadValueNames(n int) ([]string, error) {
	if n == 0 || k.nk.numberOfValues == 0 { // TODO: check behavior on microsoft api
		return []string{}, nil
	}

	list := k.nk.values
	max := len(list.values)
	if n < 0 {
		n = max
	}
	if n > max {
		n = max
	}

	names := make([]string, n)

	for i := 0; i < n; i++ {
		value, err := list.Value(uint(i))
		if err != nil {
			return nil, err
		}
		names[i] = value.name
	}
	sort.Strings(names)
	return names, nil
}
