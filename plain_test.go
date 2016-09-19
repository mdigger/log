package log

import (
	"os"
	"testing"
)

func TestPlainHandler(t *testing.T) {
	h := NewPlainHandler(os.Stdout, LstdFlags)
	h.SetFlags(h.Flags() | Lshortfile)
	h.Level()
	h.SetLevel(LDebug)
	h.SetOutput(os.Stdout)

	log := New(h)

	log.Info("info message")
	entry := log.WithField("key", "value")
	entry.Info("info message")
	entry.WithField("key 2", "value 2").Info("info message")
	entry.WithSource(0).Info("message with source")
	log.WithSource(0).Info("message with source")
	log.Error("error message")
	log.Debug("debug message")
}
