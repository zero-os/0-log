# 0-log
<<<<<<< 3ad08006ab1979f7a8628410d06c8cbee7ea837c
<<<<<<< f41440b429fd639c9da9428dd24dc4f54601f654

A package that prints out log messages in a format the [0-Core][core] log monitor can read.

## Format

The logs printed out by this package follow the specs specified in the [logging documentation][monitorFormat] and are formatted as follows:
=======
A tool that prints out log messages in a format the [0-Orchestrator][orchestrator] log monitor can read.

## Format
The logs printed out by this tool follow the specs specified in the [logging documentation][monitorFormat] and are formatted as follows:
>>>>>>> adds README documentation
=======
A package that prints out log messages in a format the [0-Core][core] log monitor can read.

## Format
The logs printed out by this package follow the specs specified in the [logging documentation][monitorFormat] and are formatted as follows:
>>>>>>> Printing the log now happens to io.Writer

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

<<<<<<< 3ad08006ab1979f7a8628410d06c8cbee7ea837c
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
=======
## Supported log levels
The package currently supports following log levels from the [0-Core log monitor][monitorLevels]:
>>>>>>> Printing the log now happens to io.Writer

* 1: stdout
* 2: stderr
* 20: result message, JSON

## Usage
<<<<<<< 3ad08006ab1979f7a8628410d06c8cbee7ea837c
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
=======
>>>>>>> Printing the log now happens to io.Writer
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
    // Log() detects if a message is multi-lined and apply the multi-line format
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
	testStr := testStruct{
		Message: "Hello world",
	}
	zerolog.Log(zerolog.LevelJSON, testStr)
    // output: 20::{"message":"Hello world"}
}
```

<<<<<<< 3ad08006ab1979f7a8628410d06c8cbee7ea837c
[orchestrator]: https://github.com/zero-os/0-orchestrator
>>>>>>> adds README documentation
=======
[core]: https://github.com/zero-os/0-core
>>>>>>> Printing the log now happens to io.Writer
[monitorFormat]: https://github.com/zero-os/0-core/blob/master/docs/monitoring/logging.md#message-format
[monitorLevels]: https://github.com/zero-os/0-core/blob/master/docs/monitoring/logging.md#log-levels