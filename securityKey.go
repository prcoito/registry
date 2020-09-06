package registry

type securityKey struct {
	signature string // must be equal to "sk"

	previousKeyOffset uint32
	nextKeyOffset     uint32

	referenceCount uint32

	ntSecurityDescriptorSize uint32
	ntSecurityDescriptor     string
}

func (sk securityKey) validate() error {
	if sk.signature != securityKeySig {
		return ErrBadSignature
	}

	return nil
}
