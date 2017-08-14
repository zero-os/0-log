package zerolog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestLogLevelSwitch(t *testing.T) {
	// stdout
	err := Log(LevelStdout, "Hello world")
	assertNoError(t, err)

	// stderr
	err = Log(LevelStderr, "Hello world")
	assertNoError(t, err)

	// json
	var ts testStruct
	err = Log(LevelJSON, ts)
	assertNoError(t, err)

	// invalid
	err = Log(255, "Hello world")
	assertError(t, err)
	assertEqual(t, ErrLevelNotValid.Error(), err.Error())

	// nil message
	err = Log(1, nil)
	assertError(t, err)
	assertEqual(t, ErrNilMessage.Error(), err.Error())

	// empty message
	err = Log(1, "")
	assertError(t, err)
	assertEqual(t, ErrNilMessage.Error(), err.Error())
}

func TestStringInput(t *testing.T) {
	// check valid strings
	//normal string
	err := Log(LevelStdout, "Hello\nworld")
	assertNoError(t, err)

	//string alias
	var sa stringAlias
	sa = "Hello world"
	err = Log(LevelStdout, sa)
	assertNoError(t, err)

	//implements stringer
	st := stringer{
		s: "lorem ipsum",
	}
	err = Log(LevelStdout, st)
	assertNoError(t, err)

	//implements TextMarshaler
	tm := textMarchal{
		"dolor sit amet",
	}
	err = Log(LevelStdout, tm)
	assertNoError(t, err)

	// check invalid strings
	//empty struct
	err = Log(LevelStdout, struct{}{})
	assertError(t, err)
	assertEqual(t, "could not turn message into string", err.Error())

	//alias
	var ia intAlias
	ia = 1
	err = Log(LevelStdout, ia)
	assertError(t, err)

	// TextMarshaler error
	var tme textMarchalError
	err = Log(LevelStdout, tme)
	assertError(t, err)

	//empty string in msgString
	_, err = msgString("")
	assertError(t, err)
	assertEqual(t, ErrNilMessage, err)
}

// test types for TestStringInput
type stringAlias string
type intAlias int

// test type that implements fmt.Stringer
type stringer struct {
	s string
}

func (s stringer) String() string {
	return s.s
}

// test type that implements encoding.TextMarshaler
type textMarchal struct {
	s string
}

func (tm textMarchal) MarshalText() ([]byte, error) {
	return []byte(tm.s), nil
}

// test type that implements encoding.TextMarshaler
// returns an error
type textMarchalError struct {
}

func (tm textMarchalError) MarshalText() ([]byte, error) {
	return nil, fmt.Errorf("An expected error")
}

func TestStatsInput(t *testing.T) {
	valFullStatMsg := MsgStat{
		Key:   "somekey",
		Value: 123.456,
		OP:    "A",
		Tags:  "foo=bar,hello=world",
	}
	// test message formatting
	str, err := msgStat(valFullStatMsg)
	if !assertNoError(t, err) || !assertNotEqual(t, "", str) {
		return
	}
	if !assertEqual(t, "somekey:123.456000|A|foo=bar,hello=world", str) {
		return
	}

	// test invalid  message
	_, err = msgStat("")
	if !assertError(t, err) {
		return
	}

	invalKey := MsgStat{
		Key:   "",
		Value: 123.456,
		OP:    "A",
		Tags:  "foo=bar,hello=world",
	}
	_, err = msgStat(invalKey)
	if !assertError(t, err) {
		return
	}

	invalOP := MsgStat{
		Key:   "someKey",
		Value: 123.456,
		OP:    "B",
		Tags:  "foo=bar,hello=world",
	}
	_, err = msgStat(invalOP)
	if !assertError(t, err) {
		return
	}

	// test logging valid Stats messages
	err = Log(LevelStats, valFullStatMsg)
	if !assertNoError(t, err) {
		return
	}

	// test logging invalid Stats messages
	err = Log(LevelStats, invalKey)
	assertError(t, err)
}

func TestJSONInput(t *testing.T) {
	// marshal test structure and check output
	tstruct := testStruct{
		TestField:      "Hello world",
		OtherTestfield: 1,
	}
	tstructExpected := "20::{\"TestField\":\"Hello world\",\"OtherTestfield\":1}\n"

	// check no error if logged
	err := Log(LevelJSON, tstruct)
	if !assertNoError(t, err) {
		return
	}

	// check output is as expected
	jsonStr, err := json.Marshal(tstruct)
	if !assertNoError(t, err) {
		return
	}

	var tw testWriter
	printLog(&tw, LevelJSON, string(jsonStr))

	assertEqual(t, tstructExpected, tw.Val)

	// write a value json can't marshal
	err = Log(LevelJSON, math.Inf(1))
	assertError(t, err)

}

func TestFormatLog(t *testing.T) {
	input1 := "stdout single line test"
	expectResult1 := "1::stdout single line test\n"
	var tw testWriter
	printLog(&tw, LevelStdout, input1)

	assertEqual(t, expectResult1, tw.Val)

	input2 := "stderr\nmultiline test"
	expectResult2 := "2:::\nstderr\nmultiline test\n:::\n"
	printLog(&tw, LevelStderr, input2)

	assertEqual(t, expectResult2, tw.Val)
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

type testWriter struct {
	Val string
}

func (tw *testWriter) Write(p []byte) (int, error) {
	tw.Val = string(p)
	return len(p), nil
}

// --------------------------------------------------------
// test assertions
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
		t.Errorf("Values were not equal\nExpected: %v\nActual: %v\n", expected, actual)
		return false
	}

	return true
}

func assertNotEqual(t *testing.T, expected, actual interface{}) bool {
	if ObjectsAreEqual(expected, actual) {
		t.Errorf("Values were equal\nExpected: %v\nActual: %v\n", expected, actual)
		return false
	}

	return true
}

// equal check
// source: github.com/stretchr/testify/assert/assertions.go
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
