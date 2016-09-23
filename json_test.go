package log

import (
	"errors"
	"os"
	"testing"
)

func TestJSON(t *testing.T) {
	json := NewJSON(os.Stdout,
		LstdFlags|LUTC|Lmicroseconds|Lshortfile|Lindent)
	log := json.Context()
	log.WithError(errors.New("error")).Info("message")
	log.WithSource(0).Error("error")
	log.Info("")
	log.Debug("debug")
}
