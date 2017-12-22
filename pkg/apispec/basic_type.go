package apispec

// BasicType encapsulates baisc types types
type BasicType string

// BasicType
const (
	BasicTypeString BasicType = "string"
	BasicTypeFloat  BasicType = "float"
	BasicTypeInt    BasicType = "int"
	BasicTypeBool   BasicType = "bool"
	BasicTypeArray  BasicType = "array"
	BasicTypeObject BasicType = "object"
)

// DetermineType attempts to determine the type of the string
func DetermineType(s string) BasicType {
	if s == "true" || s == "false" {
		return BasicTypeBool
	}
	return BasicTypeString
}
