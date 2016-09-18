package log_test

import (
	"os"
	"testing"

	"github.com/mdigger/log"
)

func TestTimestamp(t *testing.T) {
	logger := log.New(log.NewJSONHandler(os.Stdout, 0))
	for _, flag := range []int{
		log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds,
		log.LUTC | log.Ldate | log.Lmicroseconds,
		log.LUTC | log.Ldate | log.Ltime,
		log.LUTC | log.Ldate,
		log.LUTC | log.Ltime | log.Lmicroseconds,
		log.LUTC | log.Lmicroseconds,
		log.LUTC | log.Ltime,
		log.LUTC,
		log.Ldate | log.Ltime | log.Lmicroseconds,
		log.Ldate | log.Lmicroseconds,
		log.Ldate | log.Ltime,
		log.Ldate,
		log.Ltime | log.Lmicroseconds,
		log.Lmicroseconds,
		log.Ltime,
		0,
	} {
		logger.SetFlags(flag)
		logger.WithField("flag", flag).Info("")
	}
}
