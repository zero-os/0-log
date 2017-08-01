package zerolog

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevelSwitch(t *testing.T) {
	// stdout
	err := Log(LoglevelStdout, "")
	assert.NoError(t, err)

	// stderr
	err = Log(LoglevelStderr, "")
	assert.NoError(t, err)

	// json
	err = Log(LoglevelJSON, "")
	assert.NoError(t, err)

	// invalid
	err = Log(255, "")
	assert.Error(t, err)
	assert.Equal(t, ErrLevelNotValid.Error(), err.Error())

	// nil message
	err = Log(1, nil)
	assert.Error(t, err)
	assert.Equal(t, "message was nil", err.Error())

}

func TestStringInput(t *testing.T) {
	// check valid strings
	//normal string
	err := Log(LoglevelStdout, "hello\nworld")
	assert.NoError(t, err)

	//string alias
	var sa stringAlias
	sa = "hello world"
	err = Log(LoglevelStdout, sa)
	assert.NoError(t, err)

	//implements stringer
	st := strigger{
		s: "lorem ipsum",
	}
	err = Log(LoglevelStdout, st)
	assert.NoError(t, err)

	//implements TextMarshaler
	tm := textMarchal{
		"dolor sit amet",
	}
	err = Log(LoglevelStdout, tm)
	assert.NoError(t, err)

	// check invalid strings
	//empty struct
	err = Log(LoglevelStdout, struct{}{})
	assert.Error(t, err)
	assert.Equal(t, "could not turn message into string", err.Error())

	//alias
	var ia intAlias
	ia = 1
	err = Log(LoglevelStdout, ia)
	assert.Error(t, err)
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
	if !assert.NoError(t, err) {
		return
	}

	// check output is as expected
	jsonStr, err := json.Marshal(tstruct)
	if !assert.NoError(t, err) {
		return
	}
	out := formatLog(LoglevelJSON, string(jsonStr))

	assert.Equal(t, tstructExpected, out)

	// write a value json can't marshal
	err = Log(LoglevelJSON, math.Inf(1))
	assert.Error(t, err)
}

func TestFormatLog(t *testing.T) {
	input1 := "stdout single line test"
	expectResult1 := "1::stdout single line test"
	out1 := formatLog(LoglevelStdout, input1)

	assert.Equal(t, expectResult1, out1)

	input2 := "stderr\nmultiline test"
	expectResult2 := "2:::\nstderr\nmultiline test\n:::"
	out2 := formatLog(LoglevelStderr, input2)

	assert.Equal(t, expectResult2, out2)
}

func TestMultiline(t *testing.T) {
	str1 := `
This 
is a
multilined string	
`
	str2 := "This one\nis too"
	str3 := "this one is not"

	assert.True(t, isMultiline(str1))
	assert.True(t, isMultiline(str2))
	assert.False(t, isMultiline(str3))
}

type testStruct struct {
	TestField      string
	OtherTestfield int
}
