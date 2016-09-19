package log

import (
	"os"
	"testing"
)

func TestJSONHandler(t *testing.T) {
	h := NewJSONHandler(os.Stderr, LstdFlags)
	h.SetFlags(h.Flags() | Lshortfile | Lindent | LUTC)
	h.Level()
	h.SetLevel(LDebug)
	h.SetOutput(os.Stdout)
	log := New(h)

	log.Info("info message")
	entry := log.WithField("key", "value")
	entry.Info("info message")
	entry.WithField("key2", "value2").Info("info message")
	entry.WithSource(0).Info("message with source")
	log.WithSource(0).Info("message with source")
}
