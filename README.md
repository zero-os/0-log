# 0-log
<<<<<<< f41440b429fd639c9da9428dd24dc4f54601f654

A package that prints out log messages in a format the [0-Core][core] log monitor can read.

## Format

The logs printed out by this package follow the specs specified in the [logging documentation][monitorFormat] and are formatted as follows:
=======
A tool that prints out log messages in a format the [0-Orchestrator][orchestrator] log monitor can read.

## Format
The logs printed out by this tool follow the specs specified in the [logging documentation][monitorFormat] and are formatted as follows:
>>>>>>> adds README documentation

Single-line:
```
<loglevel>::<message>
```

Multi-line:
```
<loglevel>:::
<message line 1>
<message line 2>
:::
```

<<<<<<< f41440b429fd639c9da9428dd24dc4f54601f654
## Supported log levels
The package currently supports following log levels from the [0-Core log monitor][monitorLevels]:
=======
## Log levels
Possible log levels from the [specs][monitorLevels]:

* 1: stdout
* 2: stderr
* 3: message for endusers / public message
* 4: message for operator / internal message
* 5: log msg (unstructured = level5, cat=unknown)
* 6: log msg structured
* 7: warning message
* 8: ops error
* 9: critical error
* 10: statistics/monitoring message(s)
* 20: result message, JSON
* 21: result message, yaml
* 22: result message, toml
* 23: result message, hrd
* 30: job, json (full result of a job)

## Supported log levels
The tool currently supports following log levels:
>>>>>>> adds README documentation

* 1: stdout
* 2: stderr
* 20: result message, JSON

## Usage
<<<<<<< f41440b429fd639c9da9428dd24dc4f54601f654
```go
package main

import zerolog "github.com/zero-os/0-log"

func main() {
    // print to the zero-os stdout (single-line)
    zerolog.Log(zerolog.LevelStdout, "Hello world")
    // output: 1::Hello world 

    // print to the zero-os stderr
    zerolog.Log(zerolog.LevelStderr, "Hello world")
    // output: 2::Hello world 

    // print a multi-line message
    // Log() detects if a message is multi-lined
    // and applies the multi-line format if it is
    zerolog.Log(zerolog.LevelStdout, "Hello\nworld")
    /* output: 
    1:::
    Hello
    world
    :::
    */

    // print a json result message
	type testStruct struct {
		Message string `json:"message"`
	}
	zerolog.Log(zerolog.LevelJSON, testStruct{
        Message: "Hello world",
    })
    // output: 20::{"message":"Hello world"}
}
```

[core]: https://github.com/zero-os/0-core
=======
Import:
```go
import zerolog "github.com/zero-os/0-log"
```
To use the tool, simply call the `zerolog.Log(<loglevel>, <message>)` function and provide the following paramters:
* loglevel:
    this can either be the integer specified by the spec
    or provided abstractions. e.g:
    * stdout(1): `LoglevelStdout`
    * stderr(2): `LoglevelStderr`
    * JSON result message (20): `LoglevelJSON`
* message:
    For each log level different types could be expected for the message.
    * string: (e.g.: stdout(1), stderr(2))
        * a string or string alias
        * type that implements `fmt.Stringer` (`String() string`)
        * type that implements `encoding.TextMarshaler` (`MarshalText() ([]byte, error)`)
    * JSON: (e.g.: result message JSON(20))
        * any struct that can be marshalled to JSON

The tool will return an error when:
* using not [supported log levels](#Supported-log-levels)
* message is nil
* message conditions for a specific log level are not met or invalid

### Examples:
```go
zerolog.Log(zerolog.LoglevelStdout, "Hello world")
```
output:
```
1::Hello world 
```

The tool will detect if a provided message is multi-line:
 ```go
zerolog.Log(zerolog.LoglevelStderr, "Hello\nworld")
```
output:
```
2:::
Hello
world
:::
```

[orchestrator]: https://github.com/zero-os/0-orchestrator
>>>>>>> adds README documentation
[monitorFormat]: https://github.com/zero-os/0-core/blob/master/docs/monitoring/logging.md#message-format
[monitorLevels]: https://github.com/zero-os/0-core/blob/master/docs/monitoring/logging.md#log-levels