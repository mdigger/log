package log

import (
	"os"
	"testing"
)

func TestFlags(t *testing.T) {
	hp := NewPlainHandler(os.Stdout, 0)
	log := New(hp)
	hj := NewJSONHandler(os.Stdout, 0)
	log.AddHandler(hj)
	for _, flag := range []int{
		Ldate,
		Ltime,
		Lmicroseconds,
		Llongfile,
		Lshortfile,
		LUTC,
		Lindent,
		Ldate | LUTC,
		Ltime | LUTC,
		Lmicroseconds | LUTC,
		Ldate | Ltime | LUTC,
		Ldate | Lmicroseconds | LUTC,
		Llongfile | Lshortfile,
	} {
		hp.SetFlags(flag)
		hj.SetFlags(flag)
		entry := log.WithField("flag", flag)
		entry.Info("info message")
		entry.Debug("debug message")
		entry.Error("error message")
	}
}
