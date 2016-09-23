package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestHandler(t *testing.T) {
	var h handler
	h.SetLevel(DebugLevel)
	if h.Level() != DebugLevel {
		t.Error("set level")
	}

	h.SetFlags(LstdFlags)
	if h.Flags() != LstdFlags {
		t.Error("set flags")
	}

	h.SetOutput(ioutil.Discard)
	if h.w != ioutil.Discard {
		t.Error("set output")
	}
}

func TestErrorHandler(t *testing.T) {
	tw := ErrorWriter(ioutil.Discard, 200)
	ch := NewConsole(tw, LstdFlags)
	h := New(ch)
	for i := 0; i < 100; i++ {
		err := h.Infof("message %v", i)
		if err != nil {
			if err == io.ErrClosedPipe {
				fmt.Printf("writen %d logs\n", i)
				return
			}
			t.Error(err)
		}
	}
}

func ErrorWriter(w io.Writer, n int64) io.Writer {
	return &errorWriter{w, n}
}

type errorWriter struct {
	w io.Writer
	n int64
}

func (t *errorWriter) Write(p []byte) (n int, err error) {
	if t.n <= 0 {
		return len(p), nil
	}
	// real write
	n = len(p)
	if int64(n) > t.n {
		n = int(t.n)
	}
	n, err = t.w.Write(p[0:n])
	t.n -= int64(n)
	if err == nil {
		n = len(p)
	}
	if t.n <= 0 {
		err = io.ErrClosedPipe
	}
	return
}
