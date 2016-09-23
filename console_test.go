package log

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	if GetLevel() > DebugLevel {
		SetLevel(DebugLevel)
	}
	SetOutput(os.Stdout)
	SetFlags(Flags() | Lshortfile | LUTC | Lmicroseconds)

	Debug("debug")
	Info("info")
	Warning("warning")
	Error("error")

	Debugf("%v", "debug")
	Infof("%v", "info")
	Warningf("%v", "warning")
	Errorf("%v", "error")

	entry := WithField("name", "value")
	entry.Debug("debug")
	entry.Info("info")
	entry.Warning("warning")
	entry.Error("error")

	entry.Debugf("%v", "debug")
	entry.Infof("%v", "info")
	entry.Warningf("%v", "warning")
	entry.Errorf("%v", "error")

	entry = WithFields(Fields{"name": ""})
	entry.Debug("debug")
	entry.Info("info")
	entry.Warning("warning")
	entry.Error("error")

	entry = WithSource(0)
	entry.Debug("debug")
	entry.Info("info")
	entry.Warning("warning")
	entry.Error("error")

	entry = WithError(errors.New("error"))
	entry.Debug("debug")
	entry.Info("info")
	entry.Warning("warning")
	entry.Error("error")

	SetLevel(InfoLevel)
	SetFlags(0)
	WithError(nil).Debug("debug")
	WithSource(999).Error("error")
	WithFields(Fields{
		"name": true,
		"time": time.Now().Round(time.Second),
	}).WithField("name", "value").Warning("warning")

	SetOutput(nil)
	Info("info")
}

func TestConsoleTTY(t *testing.T) {
	h := NewConsole(ioutil.Discard, 0)
	h.tty = true
	h.Context().WithField("name", "value").Info("message")
	h.tty = false
	h.Context().WithField("name", "value").Info("message")
}
