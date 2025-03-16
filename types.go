package selenium

// KeyType represents the type of key that can be used
type KeyType string

const (
	// KeyTypeString represents a string key
	KeyTypeString KeyType = "string"
	// KeyTypeInt represents an integer key
	KeyTypeInt KeyType = "int"
	// KeyTypeFloat represents a float key
	KeyTypeFloat KeyType = "float"
)

// Key represents a strongly typed key value
type Key struct {
	Type  KeyType
	Value interface{}
}

// NewStringKey creates a new string key
func NewStringKey(value string) Key {
	return Key{Type: KeyTypeString, Value: value}
}

// NewIntKey creates a new integer key
func NewIntKey(value int) Key {
	return Key{Type: KeyTypeInt, Value: value}
}

// NewFloatKey creates a new float key
func NewFloatKey(value float64) Key {
	return Key{Type: KeyTypeFloat, Value: value}
}

// WaitExcTypes represents a slice of error types that should be waited for
type WaitExcTypes []error

// StdType represents the type of standard stream
type StdType string

const (
	// StdTypeFile represents a file descriptor
	StdTypeFile StdType = "file"
	// StdTypeString represents a string value
	StdTypeString StdType = "string"
	// StdTypeInt represents an integer file descriptor
	StdTypeInt StdType = "int"
)

// SubprocessStd represents a standard stream configuration
type SubprocessStd struct {
	Type  StdType
	Value interface{}
}

// NewFileStd creates a new file-based standard stream
func NewFileStd(file interface{}) SubprocessStd {
	return SubprocessStd{Type: StdTypeFile, Value: file}
}

// NewStringStd creates a new string-based standard stream
func NewStringStd(value string) SubprocessStd {
	return SubprocessStd{Type: StdTypeString, Value: value}
}

// NewIntStd creates a new integer-based standard stream
func NewIntStd(fd int) SubprocessStd {
	return SubprocessStd{Type: StdTypeInt, Value: fd}
}
