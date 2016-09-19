package log

import (
	"os"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	h := NewPlainHandler(os.Stdout, Lshortfile)
	h.SetLevel(LevelDebug)
	log := New(h)

	log.Info("info message")
	log.Infof("%v", "info message")
	log.Error("error message")
	log.Errorf("%v", "error message")
	log.Debug("debug message")
	log.Debugf("%v", "debug message")

	entry := log.WithSource(0)
	entry.Info("info message")
	entry.Infof("%v", "info message")
	entry.Error("error message")
	entry.Errorf("%v", "error message")
	entry.Debug("debug message")
	entry.Debugf("%v", "debug message")

	entry2 := entry.WithField("time", time.Now())
	entry2.Info("info message")
	entry2.Infof("%v", "info message")
	entry2.Error("error message")
	entry2.Errorf("%v", "error message")
	entry2.Debug("debug message")
	entry2.Debugf("%v", "debug message")
}
