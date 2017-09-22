# A simple structured logging

[![GoDoc](https://godoc.org/github.com/mdigger/log?status.svg)](https://godoc.org/github.com/mdigger/log)

In general, this is another "bike" for logging with blackjack.

```go
package log

import (
    "os"
    "time"

    "github.com/mdigger/log"
)

func main() {
    log.Info("info message")
    log.WithField("time", time.Now()).Debug("debug")
    log.Warn("warn", "time", time.Now(), "state", true)
}
```