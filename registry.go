package registry

import (
	"os"
)

// Open opens a registry file
func Open(f string) (Registry, error) {
	fp, err := os.Open(f)
	if err != nil {
		return Registry{}, err
	}

	h := &header{fp: fp}

	err = h.Read()
	if err != nil {
		return Registry{}, err
	}

	bins, err := getHiveBins(fp)
	if err != nil {
		return Registry{}, err
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

	return Registry{
		header:   h,
		hiveBins: bins,
		root:     root,
	}, nil
}

// Registry struct
type Registry struct {
	header *header

	root *namedKey

	hiveBins []bin

	createdByOpenKey bool
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
	k := Key{registry: r, nk: r.root}
	if path == "" {
		return k, nil
	}
	return k.OpenSubKey(path)
}

// Close closes registry file
func (r Registry) Close() error {
	if r.header != nil && r.header.fp != nil {
		return r.header.fp.Close()
	}
	return nil
}
