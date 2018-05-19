package testutils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertStructEqual checks if two structs are equal.
func AssertStructEqual(t *testing.T, expected, actual interface{}, msg string) bool {
	t.Helper()

	// NOTE(Jeff): reflect.DeepEqual has issues with time.Time comparisons as of Go 1.9
	if !assert.True(t, reflect.DeepEqual(expected, actual), msg) {
		printExpectedActual(expected, actual)
		return false
	}
	return true
}

// AssertSliceEqual checks if two slices are equal.
func AssertSliceEqual(t *testing.T, expected, actual interface{}, msg string) bool {
	t.Helper()

	// Check types
	expectedType := reflect.TypeOf(expected).Kind()
	actualType := reflect.TypeOf(actual).Kind()
	if expectedType != reflect.Slice || actualType != reflect.Slice {
		panic("must use slice types")
	}
	expectedSlice := reflect.ValueOf(expected)
	actualSlice := reflect.ValueOf(actual)

	// Check slice lengths match
	mismatch := "mismatched lengths: " + msg
	if !assert.Equal(t, expectedSlice.Len(), actualSlice.Len(), mismatch) {
		return false
	}

	// Check if each individual element is equal
	equal := true
	errCount := 0
	for i := 0; i < expectedSlice.Len(); i++ {
		if errCount >= 3 {
			t.Log("----- Too many errors... -----")
			return equal
		}
		if !AssertStructEqual(t, expectedSlice.Index(i).Interface(), actualSlice.Index(i).Interface(), msg) {
			if errCount == 0 {
				t.Logf("----- FIRST ERROR IN POSITION [%v] -----", i)
			}
			equal = false
			errCount++
		}
	}
	return equal
}
