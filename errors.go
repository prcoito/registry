package registry

import "errors"

// Errors: TODO: proper doc
var (
	ErrBadRegistry       = errors.New("Invalid registry file")
	ErrBadSequenceNumber = errors.New("Invalid sequence number")
	ErrBadBinHeader      = errors.New("Invalid bin header")
	ErrBadSignature      = errors.New("Bad header signature")
	ErrInvalidXOR        = errors.New("Invalid XOR. Corrupted registry file header")

	ErrInvalidBinHeader = errors.New("Invalid Bin header")

	ErrUnexpectedType = errors.New("Unexpected value type requested")
	ErrNotExist       = errors.New("Key/Value does not exist")
	ErrShortBuffer    = errors.New("Passed buffer too short")
)
