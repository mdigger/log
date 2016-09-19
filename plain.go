package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

// Plain logger settings for console output.
var (
	// Padding defines message padding for plain log output.
	Padding = 28
	// Colors mapping for plain logger.
	Colors = map[Level]int{
		LevelDebug: gray,
		LevelInfo:  blue,
		LevelError: red,
	}
	// Strings mapping for plain logger
	Strings = map[Level]string{
		LevelDebug: "▸",
		LevelInfo:  "•",
		LevelError: "⨯",
	}
)

// PlainHandler implements a plain text log handler.
type PlainHandler struct {
	w     io.Writer
	flag  int
	isTTY bool
	level Level // log entry's severity level.
	mu    sync.Mutex
}

// NewPlainHandler creates a new plain logger Handler.
func NewPlainHandler(w io.Writer, flag int) *PlainHandler {
	handler := &PlainHandler{flag: flag}
	handler.SetOutput(w)
	return handler
}

// Level returns the minimum event level that is supported by the logger.
func (h *PlainHandler) Level() Level {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.level
}

// SetLevel sets the minimum event level that is supported by the logger.
func (h *PlainHandler) SetLevel(level Level) {
	h.mu.Lock()
	h.level = level
	h.mu.Unlock()
}

// Flags returns the output flags for the logger.
func (h *PlainHandler) Flags() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.flag
}

// SetFlags sets the output flags for the logger.
func (h *PlainHandler) SetFlags(flag int) {
	h.mu.Lock()
	h.flag = flag
	h.mu.Unlock()
}

// SetOutput sets the output destination for the logger.
func (h *PlainHandler) SetOutput(w io.Writer) {
	h.mu.Lock()
	h.w = w
	h.isTTY = false
	if out, ok := w.(*os.File); ok {
		if fi, err := out.Stat(); err == nil {
			h.isTTY = fi.Mode()&(os.ModeDevice|os.ModeCharDevice) != 0
		}
	}
	h.mu.Unlock()
}

// Handle implements Handler.
func (h *PlainHandler) Handle(e *Entry) error {
	h.mu.Lock()
	if e.Level < h.level {
		h.mu.Unlock()
		return nil
	}
	buf := buffers.Get().(*bytes.Buffer)
	buf.Reset()
	if h.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		timestamp := e.Timestamp
		if h.flag&LUTC != 0 {
			timestamp = timestamp.UTC()
		}
		if h.flag&Ldate != 0 {
			buf.WriteString(timestamp.Format("2006-01-02 "))
		}
		if h.flag&(Ltime|Lmicroseconds) != 0 {
			buf.WriteString(timestamp.Format("15:04:05"))
			if h.flag&Lmicroseconds != 0 {
				buf.WriteString(timestamp.Format(".000000000"))
			}
			buf.WriteRune(' ')
		}
	}
	if h.flag&(Llongfile|Lshortfile) != 0 {
		if e.Source == nil {
			h.mu.Unlock()
			e.Source = getCaller(4)
			h.mu.Lock()
		}
		var file = e.Source.File
		if h.flag&Lshortfile != 0 {
			file = filepath.Base(file)
		}
		fmt.Fprintf(buf, "%s:%d ", file, e.Source.Line)
	}
	var prefix string
	if h.isTTY {
		fmt.Fprintf(buf, "\033[%dm%s\033[0m ", Colors[e.Level], Strings[e.Level])
	} else if e.Level != LevelInfo {
		prefix = fmt.Sprintf("%s: ", e.Level.String())
		buf.WriteString(prefix)
	}
	buf.WriteString(e.Message)
	if len(e.Fields) > 0 {
		if ml := Padding - len(e.Message) - len(prefix); ml > 0 &&
			h.flag&(Llongfile|Lshortfile) == 0 {
			buf.WriteString(strings.Repeat(" ", ml))
		}
		quote := func(str string) {
			if strings.ContainsAny(str, " \t\r\n\"=:") {
				fmt.Fprintf(buf, "%q", str)
			} else {
				buf.WriteString(str)
			}
		}
		names := make([]string, 0, len(e.Fields))
		for name := range e.Fields {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			buf.WriteRune(' ')
			if h.isTTY {
				fmt.Fprintf(buf, "\033[%dm", green)
			}
			quote(name)
			if h.isTTY {
				buf.WriteString("\033[0m")
			}
			buf.WriteRune('=')
			quote(fmt.Sprint(e.Fields[name]))
		}
	}
	buf.WriteRune('\n')
	buf.WriteTo(h.w)
	h.mu.Unlock()
	buffers.Put(buf)
	return nil
}

// buffers pool
var buffers = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
