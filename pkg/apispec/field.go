package apispec

import (
	"encoding/json"
	"math"
	"sort"
)

// Field stores information about payload fields.
type Field struct {
	Name     string
	Type     BasicType
	Elem     *Field
	Nested   []*Field
	Optional bool
}

// DeterminePayloadSpec determines fields spec for a given JSON payload
func DeterminePayloadSpec(s string) *Field {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		panic(err)
	}
	return &Field{
		Type:   BasicTypeObject,
		Nested: mapToFields(m),
	}
}

func mapToFields(m map[string]interface{}) []*Field {
	var fields []*Field
	for k, v := range m {
		f := valToField(v)
		f.Name = k
		fields = append(fields, f)
	}
	sort.Slice(fields, func(i int, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields
}

func sliceToField(s []interface{}) *Field {
	// TODO(Jeff): Handle merging fields in the future
	for _, e := range s {
		return valToField(e)
	}
	return &Field{}
}

func mergeFields(old, new *Field) *Field {
	if new.Type != BasicTypeObject {
		return new
	}
	return nil
}

func valToField(v interface{}) *Field {
	f := &Field{}

	switch c := v.(type) {
	case int, int32:
		f.Type = BasicTypeInt
	case float64:
		if c == math.Trunc(c) {
			f.Type = BasicTypeInt
		} else {
			f.Type = BasicTypeFloat
		}
	case string:
		f.Type = BasicTypeString
	case bool:
		f.Type = BasicTypeBool
	case map[string]interface{}:
		f.Type = BasicTypeObject
		f.Nested = mapToFields(c)
	case []interface{}:
		f.Type = BasicTypeArray
		f.Elem = sliceToField(c)
	default:
		panic("programmer: unknown field type")
	}
	return f
}
