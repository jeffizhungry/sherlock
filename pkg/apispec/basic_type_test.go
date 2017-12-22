package apispec

// func TestDetermineType(t *testing.T) {
// 	testcases := map[string]struct {
// 		s        string
// 		expected BasicType
// 	}{
// 		"should be able to determine type bool (true)": {
// 			s:        "true",
// 			expected: BasicTypeBool,
// 		},
// 		"should be able to determine type bool (false)": {
// 			s:        "true",
// 			expected: BasicTypeBool,
// 		},
// 		"should be able to default to string if unknown": {
// 			s:        "9123.sjkflajTRUE{}",
// 			expected: BasicTypeString,
// 		},
// 		"should be able to determine integer type": {
// 			s:        "9123",
// 			expected: BasicTypeInt,
// 		},
// 		"should be able to determine float types": {
// 			s:        "9000.000",
// 			expected: BasicTypeFloat,
// 		},
// 		"should be able to default to string if invalid numeric value": {
// 			s:        "9000.000.000",
// 			expected: BasicTypeString,
// 		},
// 	}
//
// 	for msg, tc := range testcases {
// 		assert.Equal(t, tc.expected, DetermineType(tc.s), msg)
// 	}
// }
