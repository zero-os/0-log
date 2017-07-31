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

	// yaml
	err = Log(LoglevelYAML, "")
	assert.NoError(t, err)

	// toml
	err = Log(LoglevelTOML, "")
	assert.Error(t, err)
	assert.Equal(t, ErrLevelNotImplemented.Error(), err.Error())

	// invalid
	err = Log(255, "")
	assert.Error(t, err)
	assert.Equal(t, ErrLevelNotValid.Error(), err.Error())

}

func TestStringInput(t *testing.T) {
	err := Log(LoglevelStdout, "hello\nworld")
	assert.NoError(t, err)

	err = Log(LoglevelStdout, struct{}{})
	assert.Error(t, err)
	assert.Equal(t, "message was not a string", err.Error())

	err = Log(LoglevelStdout, nil)
	assert.Error(t, err)
	assert.Equal(t, "message was not a string", err.Error())

	err = Log(LoglevelStderr, "hello\nworld")
	assert.NoError(t, err)

	err = Log(LoglevelStderr, struct{}{})
	assert.Error(t, err)
	assert.Equal(t, "message was not a string", err.Error())

	err = Log(LoglevelStderr, nil)
	assert.Error(t, err)
	assert.Equal(t, "message was not a string", err.Error())
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

func TestYAMLInput(t *testing.T) {
	// marshal test structure and check output
	tstruct := testStruct{
		TestField:      "Hello world",
		OtherTestfield: 1,
	}
	tstructExpected := "21:::\ntestfield: Hello world\nothertestfield: 1\n\n:::"

	// check no error if logged
	err := Log(LoglevelYAML, tstruct)
	if !assert.NoError(t, err) {
		return
	}

	// check output is as expected
	yamlStr, err := marshalYaml(tstruct)
	if !assert.NoError(t, err) {
		return
	}
	out := formatLog(LoglevelYAML, yamlStr)

	assert.Equal(t, tstructExpected, out)

	// write a value yaml can't marshal
	val := make(chan struct{})
	err = Log(LoglevelYAML, val)
	assert.Error(t, err)
	assert.Equal(t, "could not marshal provided message into YAML: cannot marshal type: chan struct {}", err.Error())
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
