package zerolog

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
)

func TestLogLevelSwitch(t *testing.T) {
	// stdout
	err := Log(LevelStdout, "Hello world")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// stderr
	err = Log(LevelStderr, "Hello world")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// json
	var ts testStruct
	err = Log(LevelJSON, ts)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// invalid
	err = Log(255, "Hello world")
	if err == nil {
		t.Error("expected an error")
	}
	if err.Error() != ErrLevelNotValid.Error() {
		t.Error("unexpected error message")
	}

	// nil message
	err = Log(1, nil)
	if err == nil {
		t.Error("expected an error")
	}
	if err.Error() != ErrNilMessage.Error() {
		t.Error("unexpected error message")
	}

	// empty message
	err = Log(1, "")
	if err == nil {
		t.Error("expected an error")
	}
	if err.Error() != ErrNilMessage.Error() {
		t.Errorf("unexpected error message %s", err)
	}

	// test nil messages
	err = Log(LevelStderr, nil)
	if err == nil {
		t.Error("expected an error")
	}
	err = Log(LevelStdout, nil)
	if err == nil {
		t.Error("expected an error")
	}
	err = Log(LevelJSON, nil)
	if err == nil {
		t.Error("expected an error")
	}
	err = Log(LevelStatistics, nil)
	if err == nil {
		t.Error("expected an error")
	}
}

func TestStringInput(t *testing.T) {
	// check valid strings
	//normal string
	err := Log(LevelStdout, "Hello\nworld")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	//string alias
	var sa stringAlias
	sa = "Hello world"
	err = Log(LevelStdout, sa)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	//implements stringer
	st := stringer{
		s: "lorem ipsum",
	}
	err = Log(LevelStdout, st)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	//implements TextMarshaler
	tm := textMarshal{
		"dolor sit amet",
	}
	err = Log(LevelStdout, tm)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// check invalid strings
	//empty struct
	err = Log(LevelStdout, struct{}{})
	if err == nil {
		t.Error("expected an error")
	}
	if err.Error() != "could not turn message into string" {
		t.Errorf("unexpected error message: %s", err)
	}

	//alias
	var ia intAlias
	ia = 1
	err = Log(LevelStdout, ia)
	if err == nil {
		t.Error("expected an error")
	}

	// TextMarshaler error
	var tme textMarchalError
	err = Log(LevelStdout, tme)
	if err == nil {
		t.Error("expected an error")
	}

	//empty string in msgString
	_, err = msgString("")
	if err == nil {
		t.Error("expected an error")
	}
	if err.Error() != ErrNilMessage.Error() {
		t.Errorf("unexpected error: %s", err)
	}
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
	valFullStatMsg := MsgStatistics{
		Key:   "somekey",
		Value: 123.456,
		OP:    AggregationAverages,
		Tags: map[string]interface{}{
			"foo": "bar",
		},
	}
	// test message formatting
	str, err := msgStatistics(valFullStatMsg)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if str == "" {
		t.Error("unexpected empty result")
	}
	if str != "somekey:123.456000|A|foo=bar" {
		t.Errorf("unexpected result: %s", str)
	}

	// test invalid  message
	_, err = msgStatistics("")
	if err == nil {
		t.Error("expected an error")
	}

	invalKey := MsgStatistics{
		Key:   "",
		Value: 123.456,
		OP:    AggregationDifferentiates,
		Tags: map[string]interface{}{
			"foo": "bar",
		},
	}
	_, err = msgStatistics(invalKey)
	if err == nil {
		t.Error("expected an error")
	}

	invalOP := MsgStatistics{
		Key:   "someKey",
		Value: 123.456,
		OP:    "B",
		Tags: map[string]interface{}{
			"foo": "bar",
		},
	}
	_, err = msgStatistics(invalOP)
	if err == nil {
		t.Error("expected an error")
	}

	// test logging valid Stats messages
	err = Log(LevelStatistics, valFullStatMsg)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// test logging invalid Stats messages
	err = Log(LevelStatistics, invalKey)
	if err == nil {
		t.Error("expected an error")
	}
}

func TestJSONInput(t *testing.T) {
	// marshal test structure and check output
	tstruct := testStruct{
		TestField:      "Hello world",
		OtherTestfield: 1,
	}
	tstructExpected := "20::{\"TestField\":\"Hello world\",\"OtherTestfield\":1}\n"
	tstructExpectedNoLogPrefix := "{\"TestField\":\"Hello world\",\"OtherTestfield\":1}"

	// check no error if logged
	err := Log(LevelJSON, tstruct)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// check output is as expected
	jsonStr, err := json.Marshal(tstruct)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	var tw testWriter
	printLog(&tw, LevelJSON, string(jsonStr))

	if tw.Val != tstructExpected {
		t.Errorf("unexpected result: %s", tw.Val)
	}

	// write a value json can't marshal
	err = Log(LevelJSON, math.Inf(1))
	if err == nil {
		t.Error("expected an error")
	}

	// check formatted string output
	str, err := msgJSON(tstruct)
	if str != tstructExpectedNoLogPrefix {
		t.Errorf("unexpected result: %s", str)
	}
}

func TestFormatLog(t *testing.T) {
	input1 := "stdout single line test"
	expectResult1 := "1::stdout single line test\n"
	var tw testWriter
	printLog(&tw, LevelStdout, input1)

	if tw.Val != expectResult1 {
		t.Errorf("unexpected result: %s", tw.Val)
	}

	input2 := "stderr\nmultiline test"
	expectResult2 := "2:::\nstderr\nmultiline test\n:::\n"
	printLog(&tw, LevelStderr, input2)

	if tw.Val != expectResult2 {
		t.Errorf("unexpected result: %s", tw.Val)
	}
}

func TestMultiline(t *testing.T) {
	str1 := `
This 
is a
multilined string	
`
	str2 := "This one\nis too"
	str3 := "this one is not"

	if !isMultiline(str1) {
		t.Error("should be true")
	}
	if !isMultiline(str2) {
		t.Error("should be true")
	}
	if isMultiline(str3) {
		t.Error("should be false")
	}
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

	if err := a.Validate(); err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if err := d.Validate(); err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if err := z.Validate(); err == nil {
		t.Error("expected an error")
	}
}

func TestMetricTags(t *testing.T) {
	emptyTags := MetricTags{}
	if emptyTags.String() != "" {
		t.Errorf("unexpected result: %s", emptyTags.String())
	}

	mockTagsString := MetricTags{
		"foo":   "bar",
		"hello": "world",
	}
	if err := checkMetricsResult(mockTagsString); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	mockTagsByteSlice := MetricTags{
		"foo":   []byte("bar"),
		"hello": []byte("world"),
	}
	if err := checkMetricsResult(mockTagsByteSlice); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	barS := metricResultStringer{s: "bar"}
	worldS := metricResultStringer{s: "world"}
	mockTagsStringer := MetricTags{
		"foo":   barS,
		"hello": worldS,
	}
	if err := checkMetricsResult(mockTagsStringer); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	barTM := metricResultTextMarshall{s: "bar"}
	worldTM := metricResultTextMarshall{s: "world"}
	mockTagsTextMarshal := MetricTags{
		"foo":   barTM,
		"hello": worldTM,
	}
	if err := checkMetricsResult(mockTagsTextMarshal); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// test last resort
	barAnon := metricResultAnon{s: "bar"}
	worldAnon := metricResultAnon{s: "world"}
	mockTagsAnon := MetricTags{
		"foo":   barAnon,
		"hello": worldAnon,
	}
	if err := checkMetricsResult(mockTagsAnon); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
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
