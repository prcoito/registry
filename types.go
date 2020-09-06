package registry

const separator = '\\'

// Signatures
const (
	registrySig    = "regf"
	binHeaderSig   = "hbin"
	namedKeySig    = "nk"
	securityKeySig = "sk"
	valueKeySig    = "vk"
	dataBlockSig   = "db"
	subKeyList1Sig = "lf"
	subKeyList2Sig = "lh"
	subKeyList3Sig = "li"
	subKeyList4Sig = "ri"
)

const (
	// Named Key flags

	nk_KEY_IS_VOLATILE   = 0x0001 // Is volatile key
	nk_KEY_HIVE_EXIT     = 0x0002 // Is mount point (of another Registry hive)
	nk_KEY_HIVE_ENTRY    = 0x0004 // Is root key (of current Registry hive)
	nk_KEY_NO_DELETE     = 0x0008 // Cannot be deleted
	nk_KEY_SYM_LINK      = 0x0010 // Is symbolic link key
	nk_KEY_COMP_NAME     = 0x0020 // Name is an ASCII string. Otherwise the name is an Unicode (UTF-16 little-endian) string
	nk_KEY_PREFEF_HANDLE = 0x0040 // Is predefined handle
	nk_KEY_VIRT_MIRRORED = 0x0080 // Unknown
	nk_KEY_VIRT_TARGET   = 0x0100 // Unknown
	nk_KEY_VIRTUAL_STORE = 0x0200 // Unknown
	nk_Unknown_1         = 0x1000 // Unknown
	nk_Unknown_2         = 0x4000 // Unknown

	// Value Key flags

	vk_VALUE_COMP_NAME = 0x0001 // Name is an ASCII string. Otherwise the name is an Unicode (UTF-16 little-endian) string
)

const (
	// REG_NONE is the undefined type
	REG_NONE uint32 = 0x00000000
	// REG_SZ is a UTF-16 little-endian string with optional end-of-string character
	REG_SZ uint32 = 0x00000001
	// REG_EXPAND_SZ is a string that contains expandable (environment) variables like %PATH%. Either in ASCII or Unicode with an end-of-string character
	REG_EXPAND_SZ uint32 = 0x00000002
	// REG_BINARY is binary data
	REG_BINARY uint32 = 0x00000003
	// REG_DWORD is a 32-bit unsigned little-endian (double word) integer
	REG_DWORD uint32 = 0x00000004
	// REG_DWORD_LITTLE_ENDIAN is a 32-bit unsigned little-endian (double word) integer
	REG_DWORD_LITTLE_ENDIAN uint32 = 0x00000004
	// REG_DWORD_BIG_ENDIAN is a 32-bit unsigned big-endian (double word) integer
	REG_DWORD_BIG_ENDIAN uint32 = 0x00000005
	// REG_LINK is a string that contains a symbolic link UTF-16 little-endian string with end-of-string character
	REG_LINK uint32 = 0x00000006
	/* REG_MULTI_SZ is a array of strings
	Array of UTF-16 little-endian strings with end-of-string character, where the array is terminated by an empty string
	Note that the termination empty string is not always present
	*/
	REG_MULTI_SZ uint32 = 0x00000007
	// REG_RESOURCE_LIST Unknown (List of hardware resources of used by a physical device driver)
	REG_RESOURCE_LIST uint32 = 0x00000008
	// REG_FULL_RESOURCE_DESCRIPTOR Unknown (List of hardware resources of controlled by a physical device driver)
	REG_FULL_RESOURCE_DESCRIPTOR uint32 = 0x00000009
	// REG_RESOURCE_REQUIREMENTS_LIST Unknown (List of hardware resources of available to a physical device driver)
	REG_RESOURCE_REQUIREMENTS_LIST uint32 = 0x0000000a
	// REG_QWORD is a 64-bit unsigned little-endian (quad word) integer
	REG_QWORD uint32 = 0x0000000b
	// REG_QWORD_LITTLE_ENDIAN 64-bit unsigned little-endian (quad word) integer
	REG_QWORD_LITTLE_ENDIAN uint32 = 0x0000000b
)

var sType = map[uint32]string{
	REG_NONE:                       "REG_NONE",
	REG_SZ:                         "REG_SZ",
	REG_EXPAND_SZ:                  "REG_EXPAND_SZ",
	REG_BINARY:                     "REG_BINARY",
	REG_DWORD_LITTLE_ENDIAN:        "REG_DWORD_LITTLE_ENDIAN",
	REG_DWORD_BIG_ENDIAN:           "REG_DWORD_BIG_ENDIAN",
	REG_LINK:                       "REG_LINK",
	REG_MULTI_SZ:                   "REG_MULTI_SZ",
	REG_RESOURCE_LIST:              "REG_RESOURCE_LIST",
	REG_FULL_RESOURCE_DESCRIPTOR:   "REG_FULL_RESOURCE_DESCRIPTOR",
	REG_RESOURCE_REQUIREMENTS_LIST: "REG_RESOURCE_REQUIREMENTS_LIST",
	REG_QWORD_LITTLE_ENDIAN:        "REG_QWORD_LITTLE_ENDIAN",
}

// Type returns the string representation of type t. If t is not valid it is returned a empty string
func Type(d uint32) string {
	return sType[d]
}

// known class names
var knownClassNames = []string{
	"activeds.dll ",
	"Class",
	"cygnus",
	"Cygwin",
	"DefaultClass ",
	"DynDRootClass ",
	"GenericClass",
	"OS2SS",
	"progman ",
	"REG_SZ",
	"Shell",
	"TCPMon",
}
