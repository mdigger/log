package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Console represents a console logger handler with configurable Level.
type Console struct {
	handler
	tty bool
}

// NewConsole returns a new initialized the console log.
func NewConsole(w io.Writer, flag int) *Console {
	var console = new(Console)
	console.SetOutput(w)
	console.SetFlags(flag)
	return console
}

// SetOutput sets the output destination for the console logger.
func (h *Console) SetOutput(w io.Writer) {
	var tty bool
	if out, ok := w.(*os.File); ok {
		if fi, err := out.Stat(); err == nil {
			tty = fi.Mode()&(os.ModeDevice|os.ModeCharDevice) != 0
		}
	}
	h.mu.Lock()
	h.w = w
	h.tty = tty
	h.mu.Unlock()
}

// Context returns a new Context for a console log.
func (h *Console) Context() *Context {
	return NewContext(h, nil)
}

// Handle implements the Handler interface.
func (h *Console) Handle(e *Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if e.Level < h.level || h.w == nil {
		return nil
	}
	flag := h.flag

	buf := buffers.Get().(*bytes.Buffer)
	buf.Reset()

	timestamp := e.Timestamp
	if flag&LUTC != 0 {
		timestamp = timestamp.UTC()
	}
	if flag&Ldate != 0 {
		buf.WriteString(timestamp.Format("2006-01-02 "))
	}
	if flag&(Ltime|Lmicroseconds) != 0 {
		buf.WriteString(timestamp.Format("15:04:05"))
		if flag&Lmicroseconds != 0 {
			buf.WriteString(timestamp.Format(".000000"))
		}
		buf.WriteRune(' ')
	}

	if flag&(Llongfile|Lshortfile) != 0 {
		if e.Source == nil {
			e.Source = NewSource(3)
		}
		var file = e.Source.File
		if flag&Lshortfile != 0 {
			file = filepath.Base(file)
		}
		fmt.Fprintf(buf, "%s:%d ", file, e.Source.Line)
	}

	level := Strings[e.Level]
	if h.tty {
		fmt.Fprintf(buf, "\033[%dm%s\033[0m", Colors[e.Level], level)
	} else {
		buf.WriteString(level)
	}
	buf.WriteRune(' ')
	buf.WriteString(e.Message)

	if len(e.Fields) > 0 {
		if ml := Padding - len(e.Message); ml > 0 && flag&Lindent != 0 {
			buf.WriteString(strings.Repeat(" ", ml))
		}

		names := make([]string, 0, len(e.Fields))
		for name := range e.Fields {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			buf.WriteRune(' ')
			if h.tty {
				fmt.Fprintf(buf, "\033[%dm%s\033[0m", green, quote(name))
			} else {
				buf.WriteString(quote(name))
			}
			buf.WriteRune('=')
			buf.WriteString(quote(fmt.Sprint(e.Fields[name])))
		}
	}

	buf.WriteRune('\n')
	_, err := buf.WriteTo(h.w)
	buffers.Put(buf)
	return err
}

func quote(str string) string {
	if str == "" || strings.ContainsAny(str, QuoteWithChars) {
		return strconv.Quote(str)
	}
	return str
}

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

// Auxiliary setting to output a console log.
var (
	// Padding defines message padding for console log output.
	Padding = 28
	// Colors mapping for console logger.
	Colors = map[Level]int{
		DebugLevel:   gray,
		InfoLevel:    blue,
		WarningLevel: yellow,
		ErrorLevel:   red,
	}
	// Strings mapping for console logger.
	Strings = map[Level]string{
		DebugLevel:   "▸",
		InfoLevel:    "•",
		WarningLevel: "⚡︎",
		ErrorLevel:   "⨯︎",
	}
	// QuoteWithChars contains a list of characters that require to wrap the
	// string in the output in quotation marks.
	QuoteWithChars = " \t\r\n\"="
)

// buffers pool
var buffers = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
