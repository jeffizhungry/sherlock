package testutils

import (
	"fmt"

	"github.com/jeffizhungry/sherlock/pkg/debug"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// printExpectedActual pretty prints expected and actual values.
func printExpectedActual(expected interface{}, actual interface{}) {
	expectedString := fmt.Sprintf("type: %T\ndata: %v", expected, debug.Prettify(expected))
	actualString := fmt.Sprintf("type: %T\ndata: %v", actual, debug.Prettify(actual))

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedString, actualString, false)
	fmt.Println("----- expected vs. actual -----")
	fmt.Println(dmp.DiffPrettyText(diffs))
}

// printExpectedActual pretty prints expected and actual values.
func printExpectedActualString(expected string, actual string) {
	expectedString := fmt.Sprintf("type: %T\ndata: %v", expected, expected)
	actualString := fmt.Sprintf("type: %T\ndata: %v", actual, actual)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedString, actualString, false)
	fmt.Println("----- expected vs. actual -----")
	fmt.Println(dmp.DiffPrettyText(diffs))
}
