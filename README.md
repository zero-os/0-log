# 0-log
A tool that prints out log messages in a format the [0-Orchestrator][orchestrator] log monitor can read.

## Format
The logs printed out by this tool follow the specs specified in the [logging documentation][monitorFormat] and are formatted as follows:

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

* 1: stdout
* 2: stderr
* 20: result message, JSON

## Usage
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
[monitorFormat]: https://github.com/zero-os/0-core/blob/master/docs/monitoring/logging.md#message-format
[monitorLevels]: https://github.com/zero-os/0-core/blob/master/docs/monitoring/logging.md#log-levels