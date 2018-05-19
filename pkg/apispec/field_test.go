package apispec

import (
	"testing"

	"github.com/jeffizhungry/sherlock/pkg/testutils"
)

func TestDeterminePayloadSpec(t *testing.T) {
	testcases := map[string]struct {
		s        string
		expected *Field
	}{
		"should be able to parse boolean field": {
			s: `{"basic":true}`,
			expected: &Field{
				Type:   BasicTypeObject,
				Nested: []*Field{{Name: "basic", Type: BasicTypeBool}},
			},
		},
		"should be able to parse integer field": {
			s: `{"basic":123}`,
			expected: &Field{
				Type:   BasicTypeObject,
				Nested: []*Field{{Name: "basic", Type: BasicTypeInt}},
			},
		},
		"should be able to parse float field": {
			s: `{"basic":10.123}`,
			expected: &Field{
				Type:   BasicTypeObject,
				Nested: []*Field{{Name: "basic", Type: BasicTypeFloat}},
			},
		},
		"should be able to parse string field": {
			s: `{"basic":"needs"}`,
			expected: &Field{
				Type:   BasicTypeObject,
				Nested: []*Field{{Name: "basic", Type: BasicTypeString}},
			},
		},
		"should be able to parse nested object fields": {
			s: `{"basic":{"nested":"object"}}`,
			expected: &Field{
				Type: BasicTypeObject,
				Nested: []*Field{
					{
						Name: "basic",
						Type: BasicTypeObject,
						Nested: []*Field{
							{
								Name: "nested",
								Type: BasicTypeString,
							},
						},
					},
				},
			},
		},
		"should be able to parse nested array fields of type string": {
			s: `{"basic":["string"]}`,
			expected: &Field{
				Type: BasicTypeObject,
				Nested: []*Field{
					{
						Name: "basic",
						Type: BasicTypeArray,
						Elem: &Field{Type: BasicTypeString},
					},
				},
			},
		},
		"should be able to parse nested array fields of type object": {
			s: `{"basic":[{"nested":"object"}]}`,
			expected: &Field{
				Type: BasicTypeObject,
				Nested: []*Field{
					{
						Name: "basic",
						Type: BasicTypeArray,
						Elem: &Field{
							Type: BasicTypeObject,
							Nested: []*Field{
								{
									Name: "nested",
									Type: BasicTypeString,
								},
							},
						},
					},
				},
			},
		},
		// "should be able to parse boolean field": {
		// 	s:        `{"basic":true}`,
		// 	expected: []*Field{{Name: "basic", Type: BasicTypeBool}},
		// },
	}
	for msg, tc := range testcases {
		actual := DeterminePayloadSpec(tc.s)
		testutils.AssertSliceEqual(t, tc.expected, actual, msg)
	}
}
