package zerolog

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func TestLogLevelSwitch(t *testing.T) {
	// stdout
	err := Log(LoglevelStdout, "")
	assertNoError(t, err)

	// stderr
	err = Log(LoglevelStderr, "")
	assertNoError(t, err)

	// json
	err = Log(LoglevelJSON, "")
	assertNoError(t, err)

	// invalid
	err = Log(255, "")
	assertError(t, err)
	assertEqual(t, ErrLevelNotValid.Error(), err.Error())

	// nil message
	err = Log(1, nil)
	assertError(t, err)
	assertEqual(t, "message was nil", err.Error())
}

func TestStringInput(t *testing.T) {
	// check valid strings
	//normal string
	err := Log(LoglevelStdout, "hello\nworld")
	assertNoError(t, err)

	//string alias
	var sa stringAlias
	sa = "hello world"
	err = Log(LoglevelStdout, sa)
	assertNoError(t, err)

	//implements stringer
	st := strigger{
		s: "lorem ipsum",
	}
	err = Log(LoglevelStdout, st)
	assertNoError(t, err)

	//implements TextMarshaler
	tm := textMarchal{
		"dolor sit amet",
	}
	err = Log(LoglevelStdout, tm)
	assertNoError(t, err)

	// check invalid strings
	//empty struct
	err = Log(LoglevelStdout, struct{}{})
	assertError(t, err)
	assertEqual(t, "could not turn message into string", err.Error())

	//alias
	var ia intAlias
	ia = 1
	err = Log(LoglevelStdout, ia)
	assertError(t, err)
}

// test types for TestStringInput
type stringAlias string
type intAlias int

// test type that implements fmt.Stringger
type strigger struct {
	s string
}

func (s strigger) String() string {
	return s.s
}

// test type that implements encoding.TextMarshaler
type textMarchal struct {
	s string
}

func (tm textMarchal) MarshalText() ([]byte, error) {
	return []byte(tm.s), nil
}

func TestJSONInput(t *testing.T) {
	// marshal test structure and check output
	tstruct := testStruct{
		TestField:      "Hello world",
		OtherTestfield: 1,
	}
	tstructExpected := "20::{\"TestField\":\"Hello world\",\"OtherTestfield\":1}"

	// check no error if logged
	err := Log(LoglevelJSON, tstruct)
	if !assertNoError(t, err) {
		return
	}

	// check output is as expected
	jsonStr, err := json.Marshal(tstruct)
	if !assertNoError(t, err) {
		return
	}
	out := formatLog(LoglevelJSON, string(jsonStr))

	assertEqual(t, tstructExpected, out)

	// write a value json can't marshal
	err = Log(LoglevelJSON, math.Inf(1))
	assertError(t, err)
}

func TestFormatLog(t *testing.T) {
	input1 := "stdout single line test"
	expectResult1 := "1::stdout single line test"
	out1 := formatLog(LoglevelStdout, input1)

	assertEqual(t, expectResult1, out1)

	input2 := "stderr\nmultiline test"
	expectResult2 := "2:::\nstderr\nmultiline test\n:::"
	out2 := formatLog(LoglevelStderr, input2)

	assertEqual(t, expectResult2, out2)
}

func TestMultiline(t *testing.T) {
	str1 := `
This 
is a
multilined string	
`
	str2 := "This one\nis too"
	str3 := "this one is not"

	assertTrue(t, isMultiline(str1))
	assertTrue(t, isMultiline(str2))
	assertFalse(t, isMultiline(str3))
}

type testStruct struct {
	TestField      string
	OtherTestfield int
}

// test assertion
func assertTrue(t *testing.T, val bool) bool {
	if !val {
		t.Error("Expected True, got False")
		return false
	}
	return true
}

func assertFalse(t *testing.T, val bool) bool {
	if val {
		t.Error("Expected False, got True")
		return false
	}
	return true
}
func assertNoError(t *testing.T, err error) bool {
	if err != nil {
		t.Errorf("Did not expect error, got: %s", err)
		return false
	}

	return true
}

func assertError(t *testing.T, err error) bool {
	if err == nil {
		t.Error("Expected an error, got nil")
		return false
	}

	return true
}

func assertEqual(t *testing.T, expected, actual interface{}) bool {
	if !ObjectsAreEqual(expected, actual) {
		t.Error("Values were not equal")
		return false
	}

	return true
}

// equal check copied from:
// github.com/stretchr/testify/assert/assertions.go
func ObjectsAreEqual(expected, actual interface{}) bool {

	if expected == nil || actual == nil {
		return expected == actual
	}
	if exp, ok := expected.([]byte); ok {
		act, ok := actual.([]byte)
		if !ok {
			return false
		} else if exp == nil || act == nil {
			return exp == nil && act == nil
		}
		return bytes.Equal(exp, act)
	}
	return reflect.DeepEqual(expected, actual)
}
