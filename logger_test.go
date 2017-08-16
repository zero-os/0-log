package zerolog

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/zero-os/0-log/assert"
)

func TestLogLevelSwitch(t *testing.T) {
	// stdout
	err := Log(LevelStdout, "Hello world")
	assert.NoError(t, err)

	// stderr
	err = Log(LevelStderr, "Hello world")
	assert.NoError(t, err)

	// json
	var ts testStruct
	err = Log(LevelJSON, ts)
	assert.NoError(t, err)

	// invalid
	err = Log(255, "Hello world")
	assert.Error(t, err)
	assert.Equal(t, ErrLevelNotValid.Error(), err.Error())

	// nil message
	err = Log(1, nil)
	assert.Error(t, err)
	assert.Equal(t, ErrNilMessage.Error(), err.Error())

	// empty message
	err = Log(1, "")
	assert.Error(t, err)
	assert.Equal(t, ErrNilMessage.Error(), err.Error())

	// test nil messages
	err = Log(LevelStderr, nil)
	assert.Error(t, err)
	err = Log(LevelStdout, nil)
	assert.Error(t, err)
	err = Log(LevelJSON, nil)
	assert.Error(t, err)
	err = Log(LevelStatistics, nil)
	assert.Error(t, err)
}

func TestStringInput(t *testing.T) {
	// check valid strings
	//normal string
	err := Log(LevelStdout, "Hello\nworld")
	assert.NoError(t, err)

	//string alias
	var sa stringAlias
	sa = "Hello world"
	err = Log(LevelStdout, sa)
	assert.NoError(t, err)

	//implements stringer
	st := stringer{
		s: "lorem ipsum",
	}
	err = Log(LevelStdout, st)
	assert.NoError(t, err)

	//implements TextMarshaler
	tm := textMarshal{
		"dolor sit amet",
	}
	err = Log(LevelStdout, tm)
	assert.NoError(t, err)

	// check invalid strings
	//empty struct
	err = Log(LevelStdout, struct{}{})
	assert.Error(t, err)
	assert.Equal(t, "could not turn message into string", err.Error())

	//alias
	var ia intAlias
	ia = 1
	err = Log(LevelStdout, ia)
	assert.Error(t, err)

	// TextMarshaler error
	var tme textMarchalError
	err = Log(LevelStdout, tme)
	assert.Error(t, err)

	//empty string in msgString
	_, err = msgString("")
	assert.Error(t, err)
	assert.Equal(t, ErrNilMessage, err)
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
type textMarshal struct {
	s string
}

func (tm textMarshal) MarshalText() ([]byte, error) {
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
		OP:    AggregationAverages,
		Tags: map[string]interface{}{
			"foo": "bar",
		},
	}
	// test message formatting
	str, err := msgStat(valFullStatMsg)
	if !assert.NoError(t, err) || !assert.NotEqual(t, "", str) {
		return
	}
	if !assert.Equal(t, "somekey:123.456000|A|foo=bar", str) {
		return
	}

	// test invalid  message
	_, err = msgStat("")
	if !assert.Error(t, err) {
		return
	}

	invalKey := MsgStat{
		Key:   "",
		Value: 123.456,
		OP:    AggregationDifferentiates,
		Tags: map[string]interface{}{
			"foo": "bar",
		},
	}
	_, err = msgStat(invalKey)
	if !assert.Error(t, err) {
		return
	}

	invalOP := MsgStat{
		Key:   "someKey",
		Value: 123.456,
		OP:    "B",
		Tags: map[string]interface{}{
			"foo": "bar",
		},
	}
	_, err = msgStat(invalOP)
	if !assert.Error(t, err) {
		return
	}

	// test logging valid Stats messages
	err = Log(LevelStatistics, valFullStatMsg)
	if !assert.NoError(t, err) {
		return
	}

	// test logging invalid Stats messages
	err = Log(LevelStatistics, invalKey)
	assert.Error(t, err)
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
	if !assert.NoError(t, err) {
		return
	}

	// check output is as expected
	jsonStr, err := json.Marshal(tstruct)
	if !assert.NoError(t, err) {
		return
	}

	var tw testWriter
	printLog(&tw, LevelJSON, string(jsonStr))

	assert.Equal(t, tstructExpected, tw.Val)

	// write a value json can't marshal
	err = Log(LevelJSON, math.Inf(1))
	assert.Error(t, err)

}

func TestFormatLog(t *testing.T) {
	input1 := "stdout single line test"
	expectResult1 := "1::stdout single line test\n"
	var tw testWriter
	printLog(&tw, LevelStdout, input1)

	assert.Equal(t, expectResult1, tw.Val)

	input2 := "stderr\nmultiline test"
	expectResult2 := "2:::\nstderr\nmultiline test\n:::\n"
	printLog(&tw, LevelStderr, input2)

	assert.Equal(t, expectResult2, tw.Val)
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

type testWriter struct {
	Val string
}

func (tw *testWriter) Write(p []byte) (int, error) {
	tw.Val = string(p)
	return len(p), nil
}

func TestAggregationType(t *testing.T) {
	a := AggregationAverages
	d := AggregationDifferentiates
	var z AggregationType
	z = "z"

	assert.NoError(t, a.Validate())
	assert.NoError(t, d.Validate())
	assert.Error(t, z.Validate())
}

func TestMetricTags(t *testing.T) {
	emptyTags := MetricTags{}
	assert.Equal(t, "", emptyTags.String())

	mockTagsString := MetricTags{
		"foo":   "bar",
		"hello": "world",
	}
	assert.NoError(t, checkMetricsResult(mockTagsString))

	mockTagsByteSlice := MetricTags{
		"foo":   []byte("bar"),
		"hello": []byte("world"),
	}
	assert.NoError(t, checkMetricsResult(mockTagsByteSlice))

	barS := metricResultStringer{s: "bar"}
	worldS := metricResultStringer{s: "world"}
	mockTagsStringer := MetricTags{
		"foo":   barS,
		"hello": worldS,
	}
	assert.NoError(t, checkMetricsResult(mockTagsStringer))

	barTM := metricResultTextMarshall{s: "bar"}
	worldTM := metricResultTextMarshall{s: "world"}
	mockTagsTextMarshal := MetricTags{
		"foo":   barTM,
		"hello": worldTM,
	}
	assert.NoError(t, checkMetricsResult(mockTagsTextMarshal))

	// test last resort
	barAnon := metricResultAnon{s: "bar"}
	worldAnon := metricResultAnon{s: "world"}
	mockTagsAnon := MetricTags{
		"foo":   barAnon,
		"hello": worldAnon,
	}
	assert.NoError(t, checkMetricsResult(mockTagsAnon))

}

func checkMetricsResult(mt MetricTags) error {
	// order is not guaranteed
	switch mt.String() {
	case "foo=bar,hello=world", "hello=world,foo=bar":
		return nil
	// for last resort check
	case "foo={bar},hello={world}", "hello={world},foo={bar}":
		return nil
	default:
		return fmt.Errorf("unexpected MetricTags result: %s", mt)
	}
}

type metricResultStringer struct {
	s string
}

func (mrs metricResultStringer) String() string {
	return mrs.s
}

type metricResultTextMarshall struct {
	s string
}

func (mrtm metricResultTextMarshall) MarshalText() ([]byte, error) {
	return []byte(mrtm.s), nil
}

type metricResultAnon struct {
	s string
}
