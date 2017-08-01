//Package zerolog is a tool that prints messages in a format the 0-Ochestrator can read and use for selfhealing and other features
package zerolog

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Loglevel represents the level that will be logged at
type Loglevel uint8

const (
	// LoglevelStdout stdout loglevel
	LoglevelStdout Loglevel = 1
	// LoglevelStderr stderr loglevel
	LoglevelStderr Loglevel = 2
	// LoglevelJSON json loglevel
	LoglevelJSON Loglevel = 20
)

// ErrLevelNotValid defines an error where the loggin level is not supported/valid
var ErrLevelNotValid = errors.New("logging level not valid")

// Log prints a message in the Orchestrator logging format
func Log(lvl Loglevel, message interface{}) error {
	var msgStr string

	// check if message is nil
	if message == nil {
		return fmt.Errorf("message was nil")
	}

	switch lvl {
	// string messages
	case LoglevelStdout, LoglevelStderr:
		var err error
		msgStr, err = msgString(message)
		if err != nil {
			return err
		}
	// json messages
	case LoglevelJSON:
		msgBs, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("could not marshal provided message into JSON: %s", err)
		}
		msgStr = string(msgBs)
	default:
		return ErrLevelNotValid
	}

	// print messages
	fmt.Println(formatLog(lvl, msgStr))

	return nil
}

// checks if the interface can be turned into a string and returns it as such
// boolean returned determines if the interface could be turned into a string (ok)
func msgString(msg interface{}) (string, error) {

	// check if msg reflects string
	if reflect.TypeOf(msg).Kind() == reflect.String {
		return reflect.ValueOf(msg).String(), nil
	}

	// check if implements Stringer
	if m, ok := msg.(fmt.Stringer); ok {
		return m.String(), nil
	}

	// check if implements TextMarshaler
	if m, ok := msg.(encoding.TextMarshaler); ok {
		str, err := m.MarshalText()
		if err != nil {
			return "", fmt.Errorf("could not TextMarshal provided message: %s", err)
		}
		return string(str), nil
	}

	return "", fmt.Errorf("could not turn message into string")
}

// formatLog formats the log output
func formatLog(lvl Loglevel, msg string) string {
	ml := isMultiline(msg)
	if ml {
		return fmt.Sprintf("%d:::\n%s\n:::", lvl, msg)
	}

	return fmt.Sprintf("%d::%s", lvl, msg)
}

// isMultiline return true when a string constains \n
func isMultiline(str string) bool {
	return strings.Contains(str, "\n")
}
