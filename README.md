# 0-log
A package that prints out log messages in a format the [0-Core][core] log monitor can read.

## Format
The logs printed out by this package follow the specs specified in the [logging documentation][monitorFormat] and are formatted as follows:

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

## Supported log levels
The package currently supports following log levels from the [0-Core log monitor][monitorLevels]:

* 1: stdout
* 2: stderr
* 20: result message, JSON

## Usage
```go
package main

import zerolog "github.com/zero-os/0-log"

func main() {
    // print to the zero-os stdout
    zerolog.Log(zerolog.LevelStdout, "Hello world")
    // output: 1::Hello world 

    // print to the zero-os stderr
    zerolog.Log(zerolog.LevelStderr, "Hello world")
    // output: 2::Hello world 

    // print multi-line
    zerolog.Log(zerolog.LevelStdout, "Hello\nworld")
    /* output: 
    1:::
    Hello
    world
    :::
    */
}
```

[core]: https://github.com/zero-os/0-core
[monitorFormat]: https://github.com/zero-os/0-core/blob/master/docs/monitoring/logging.md#message-format
[monitorLevels]: https://github.com/zero-os/0-core/blob/master/docs/monitoring/logging.md#log-levels