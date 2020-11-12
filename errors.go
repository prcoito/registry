package registry

import (
	"errors"
)

var (
	// ErrBadRegistry is the error returned on OpenKey or Open if the passed file is not a
	// valid registry file
	ErrBadRegistry = errors.New("Invalid registry file")

	// ErrUnexpectedType is returned when is requested a value with type different from Key value type
	ErrUnexpectedType = errors.New("Unexpected value type requested")
	// ErrNotExist is returned when a Key or Value with the given path does not exist
	ErrNotExist = errors.New("Key/Value does not exist")
	// ErrShortBuffer is returned when the passed buffer to GetValue is too short to hold all content
	ErrShortBuffer = errors.New("Passed buffer too short")

	// ErrOutOfBounds is the error returned in functions that received an index and the received index
	// is out of bounds
	ErrOutOfBounds = errors.New("Index out of bounds")

	// ErrCorruptRegistry is returned when there is data corruption or when a read operation fails
	ErrCorruptRegistry = errors.New("Corrupt registry file")
)

var (
	// errBadSequenceNumber is one cause for ErrBadRegistry
	errBadSequenceNumber = errors.New("Invalid sequence number")
	// errBadSignature is one cause for Read errors
	errBadSignature = errors.New("Bad header signature")
	// errInvalidXOR is one cause for ErrBadRegistry
	errInvalidXOR = errors.New("Invalid XOR. Corrupted registry file header")

	// errInvalidBinHeader is returned if the bin header is not valid
	errInvalidBinHeader = errors.New("Invalid Bin header")

	errRootNotFound = errors.New("registry root key not found")

	errInvalidHash = errors.New("Element hash invalid")
)

type errorW struct {
	err      error
	cause    error
	function string
}

func (e errorW) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e errorW) errorMsg() string {
	if e.cause == nil {
		return "Error at function " + e.function + ": " + e.Error() + " caused by unknown"
	}
	return "Error at function " + e.function + ": " + e.Error() + " caused by " + e.cause.Error()
}
