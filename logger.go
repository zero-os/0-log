package zerolog

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/siddontang/go/log"
)

var (
	// ErrLevelNotImplemented defines en error where the logging level is not yet implented
	ErrLevelNotImplemented = errors.New("logging level not yet implemented")
	// ErrLevelNotValid defines an error where the loggin level is not supported/valid
	ErrLevelNotValid = errors.New("logging level not valid")
)

// Loglevel represents the level that will be logged at
type Loglevel uint8

const (
	// LoglevelStdout stdout log level
	LoglevelStdout Loglevel = 1
	// LoglevelStderr stderr log level
	LoglevelStderr Loglevel = 2
	// LoglevelJSON json log level
	LoglevelJSON Loglevel = 20
	// LoglevelYAML yaml log level
	LoglevelYAML Loglevel = 21
	// LoglevelTOML toml log level
	LoglevelTOML Loglevel = 22
)

// Log prints a message in the Orchestrator logging format
func Log(lvl Loglevel, message interface{}) error {
	var msgStr string

	switch lvl {
	// string messages
	case LoglevelStdout, LoglevelStderr:
		var ok bool
		msgStr, ok = message.(string)
		if !ok {
			return errors.New("message was not a string")
		}
	// json messages
	case LoglevelJSON:
		msgBs, err := json.Marshal(&message)
		if err != nil {
			return fmt.Errorf("could not marshal provided message into JSON: %s", err)
		}
		msgStr = string(msgBs)
	// yaml messages
	case LoglevelYAML:
		yamlStr, err := marshalYaml(message)
		if err != nil {
			return fmt.Errorf("could not marshal provided message into YAML: %s", err)
		}
		msgStr = yamlStr
	//toml messages
	case LoglevelTOML:
		return ErrLevelNotImplemented
	default:
		return ErrLevelNotValid
	}

	// print messages
	fmt.Println(formatLog(lvl, msgStr))

	return nil
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
	buf := strings.Split(str, "\n")
	if len(buf) > 1 {
		return true
	}

	return false
}

func marshalYaml(msg interface{}) (res string, err error) {
	// catch yaml panics
	defer func() {
		if r := recover(); r != nil {

			log.Errorf("Recovered YAML panic: %s", r)
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()
	msgBs, yerr := yaml.Marshal(&msg)
	if yerr != nil {
		err = fmt.Errorf("could not marshal provided message into YAML: %s", yerr)
		return
	}
	res = string(msgBs)
	return
}
