//Package assert provides testing tools
package assert

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"testing"
)

// True checks if values is true
func True(t *testing.T, val bool) bool {
	if !val {
		return failf(t, "Expected True, got False")
	}

	return true
}

// False checks if value is false
func False(t *testing.T, val bool) bool {
	if val {
		return failf(t, "Expected False, got True")
	}

	return true
}

// Error checks if error is not nil
func Error(t *testing.T, err error) bool {
	if err == nil {
		return failf(t, "Expected an error, got nil")
	}

	return true
}

// NoError checks if Error is nil
func NoError(t *testing.T, err error) bool {
	if err != nil {
		return failf(t, "Unexpected error, got error: %s", err)
	}

	return true
}

// Equal checks if expected value equals the actual value
func Equal(t *testing.T, expected, actual interface{}) bool {
	if !objectsAreEqual(expected, actual) {
		return failf(t, "Values were unexpectedly not equal\n\t\tExpected: %v\n\t\tActual: %v\n", expected, actual)
	}

	return true
}

// NotEqual checks if expected value does not equals the actual value
func NotEqual(t *testing.T, expected, actual interface{}) bool {
	if objectsAreEqual(expected, actual) {
		return failf(t, "Values were unexpectedly equal\n\t\tExpected: %v\n\t\tActual: %v\n", expected, actual)
	}

	return true
}

// objectsAreEqual equal check logic
func objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	// strings are pretty common
	if exp, ok := expected.(string); ok {
		act, ok := actual.(string)
		if !ok {
			return false
		}
		return exp == act
	}

	return reflect.DeepEqual(expected, actual)
}

// failf reports a failure
// needs to be called from public function for accurate error location pointer
func failf(t *testing.T, failureMessage string, msgArgs ...interface{}) bool {
	_, file, line, _ := runtime.Caller(2)
	file = path.Base(file)

	errormsg := fmt.Sprintf(failureMessage, msgArgs...)

	t.Errorf("\n\n\tError at: %s:%d\n\tError: %s\n\n", file, line, errormsg)
	return false
}
