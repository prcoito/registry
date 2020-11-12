package registry

import (
	"io"
	"os"
)

// Registry struct
type Registry struct {
	fp  *os.File // set if Open/OpenKey
	rws io.ReadWriteSeeker

	header *header

	root *namedKey

	hiveBins []bin

	createdByOpenKey bool
}

// Open opens a registry file
func Open(f string) (Registry, error) {
	fp, err := os.Open(f)
	if err != nil {
		return Registry{}, err
	}

	h := newHeader(fp)

	err = h.Read()
	if err != nil {
		return Registry{}, errorW{function: "Open h.Read", err: ErrBadRegistry, cause: err}
	}

	bins, err := getHiveBins(fp)
	if err != nil {
		return Registry{}, errorW{function: "Open getHiveBins", err: ErrBadRegistry, cause: err}
	}

	var root *namedKey
	for _, bin := range bins {
		nk, ok := bin.cell.data.(*namedKey)
		if ok && nk.name == "ROOT" {
			root = nk
			root.binOffset = bin.offset + int64(bin.header.size)
			break
		}
	}

	if root == nil {
		return Registry{}, errorW{function: "Open findRoot", err: ErrBadRegistry, cause: errRootNotFound}
	}

	return Registry{
		header:   h,
		hiveBins: bins,
		root:     root,
		fp:       fp,
		rws:      fp,
	}, nil
}

// OpenKey opens a new key in file located at path
func OpenKey(file, path string) (Key, error) {
	registry, err := Open(file)
	if err != nil {
		return Key{}, err
	}
	registry.createdByOpenKey = true
	return registry.OpenKey(path)
}

// OpenKey opens a new key located at path
// If path is empty, it is returned the root key
func (r Registry) OpenKey(path string) (Key, error) {
	k := newKey(r, r.rws, r.root)
	if path == "" {
		return k, nil
	}
	return k.OpenSubKey(path)
}

// Close closes registry file
func (r Registry) Close() error {
	if r.fp != nil {
		return r.fp.Close()
	}
	return nil
}
