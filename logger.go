/*Package zerolog is a package that prints messages in a format the 0-Core log monitor can read and use for logging, statistics, selfhealing and other features.

Usage:
	zerolog.Log(zerolog.LevelStdout, "Hello world")
Output:
	1::Hello world

Message

Accepted message types may very on provided log level:

String message (e.g.: LevelStdout, LevelStderr) takes strings, string aliases, types that implement fmt.Stringer, types that implement encoding.TextMarshaler.

Stats message (e.g: LevelStats) takes a MsgStat to have fields and validation for data required by the 0-Core statistics monitor
https://github.com/zero-os/0-core/blob/master/docs/monitoring/stats.md

JSON messages (e.g.: LevelJSON) takes any type that can be marshalled to JSON.

Additional info

package docs: https://github.com/zero-os/0-log/blob/master/README.md

Information about the 0-Core monitoring: https://github.com/zero-os/0-core/blob/master/docs/monitoring/logging.md
*/
package zerolog

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// Level represents the level that will be logged at
type Level uint8

const (
	// LevelStdout stdout
	LevelStdout Level = 1
	// LevelStderr stderr
	LevelStderr Level = 2
	// LevelStats statistics/monitoring message
	LevelStats Level = 10
	// LevelJSON JSON result message
	LevelJSON Level = 20
)

var (
	// ErrLevelNotValid defines an error where the loggin level is not supported/valid
	ErrLevelNotValid = errors.New("logging level not valid")
	// ErrNilMessage represents an error where the supplied message was nil
	ErrNilMessage = errors.New("message was nil")
)

// Log prints a message in the 0-Core logging format
func Log(lvl Level, message interface{}) error {
	// check if message is nil/empty
	if message == nil {
		return ErrNilMessage
	}

	var msgStr string
	var err error

	// check if message is nil/empty
	if message == nil || message == "" {
		return ErrNilMessage
	}

	switch lvl {
	// string messages
	case LevelStdout, LevelStderr:
		msgStr, err = msgString(message)
		if err != nil {
			return err
		}
	// stats messages
	case LevelStats:
		msgStr, err = msgStat(message)
		if err != nil {
			return err
		}
	// json messages
	case LevelJSON:
		msgBs, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("could not marshal provided message into JSON: %s", err)
		}
		msgStr = string(msgBs)
	default:
		return ErrLevelNotValid
	}

	// print messages
	printLog(os.Stdout, lvl, msgStr)

	return nil
}

// msgString checks if the interface can be turned into a string and returns it as such
func msgString(msg interface{}) (string, error) {
	// check if string
	if str, ok := msg.(string); ok {
		if str == "" {
			return "", ErrNilMessage
		}

		return str, nil
	}

	// check if implements fmt.Stringer
	if m, ok := msg.(fmt.Stringer); ok {
		return m.String(), nil
	}

	// check if implements encoding.TextMarshaler
	if m, ok := msg.(encoding.TextMarshaler); ok {
		str, err := m.MarshalText()
		if err != nil {
			return "", fmt.Errorf("could not MarshalText provided message: %s", err)
		}
		return string(str), nil
	}

	// check if msg reflects string
	if reflect.TypeOf(msg).Kind() == reflect.String {
		return reflect.ValueOf(msg).String(), nil
	}

	return "", fmt.Errorf("could not turn message into string")
}

// msgStat validates and formats a statistics log message
// validates if message conforms to 0-core statistics spec:
// https://github.com/zero-os/0-core/blob/master/docs/monitoring/stats.md
func msgStat(msg interface{}) (string, error) {
	// check if msg is type MsgStat
	statMsg, ok := msg.(MsgStat)
	if !ok {
		return "", fmt.Errorf("statistics log message was not of type MsgStat")
	}

	err := statMsg.Validate()
	if err != nil {
		return "", err
	}

	str := fmt.Sprintf("%s:%f|%s", statMsg.Key, statMsg.Value, statMsg.OP)

	if statMsg.Tags != "" {
		str = str + "|" + statMsg.Tags
	}

	return str, nil
}

// MsgStat represents the data needed for a statistics message
type MsgStat struct {
	Key   string
	Value float64
	OP    string
	Tags  string
}

// Validate validates the MsgStat according to spec:
// https://github.com/zero-os/0-core/blob/master/docs/monitoring/stats.md
func (msg *MsgStat) Validate() error {
	if msg.Key == "" {
		return fmt.Errorf("stats message does not contain a key")
	}
	if msg.OP != "A" && msg.OP != "D" {
		return fmt.Errorf("stats message contains invalid OP")
	}
	return nil
}

// printLog formats the log and writes it to the io.Writer
func printLog(w io.Writer, lvl Level, msg string) {
	ml := isMultiline(msg)
	if ml {
		fmt.Fprintf(w, "%d:::\n%s\n:::\n", lvl, msg)
		return
	}

	fmt.Fprintf(w, "%d::%s\n", lvl, msg)
}

// isMultiline return true when a string constains \n
func isMultiline(str string) bool {
	return strings.Contains(str, "\n")
}
