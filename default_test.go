package log

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	SetFlags(Lshortfile)
	SetLevel(Level() - 1)
	Info("info message")
	Infof("%v", "info message")
	Error("error message")
	Errorf("%v", "error message")
	Debug("debug message")
	Debugf("%v", "debug message")

	entry := WithSource(0)
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

	Trace("trace message").Stop(nil)
	Tracef("trace %v", "message").Stop(nil)

	entry3 := WithField("time", time.Now())
	entry3.Info("info message")
	entry3.Infof("%v", "info message")
	entry3.Error("error message")
	entry3.Errorf("%v", "error message")
	entry3.Debug("debug message")
	entry3.Debugf("%v", "debug message")

	entry4 := WithFields(Fields{"time": time.Now()})
	entry4.Info("info message")
	entry4.Infof("%v", "info message")
	entry4.Error("error message")
	entry4.Errorf("%v", "error message")
	entry4.Debug("debug message")
	entry4.Debugf("%v", "debug message")

	WithError(errors.New("error")).Error("error message")

	SetLevel(Level())
	SetFlags(Flags())
	SetOutput(ioutil.Discard)
}
