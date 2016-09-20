package log

import (
	"os"
	"testing"
	"time"
)

func TestTrace(t *testing.T) {
	h := NewPlainHandler(os.Stdout, Lshortfile)
	h.SetLevel(DebugLevel)
	log := New(h)

	open := func(filename string) (err error) {
		defer log.WithField("file", filename).Trace("open").Stop(&err)
		file, err := os.Open(filename)
		if err == nil {
			time.Sleep(time.Second)
			file.Close()
		}
		return err
	}

	for _, filename := range []string{"~README.md", "trace_test.go"} {
		open(filename)
	}
	log.Tracef("%v", "tracef").
		AddFields(Fields{"int": 8}).
		AddField("bool", true).
		Stop(nil)
}
