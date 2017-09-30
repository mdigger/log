package log

import (
	"strconv"
	"sync"
	"unicode/utf8"
)

// Use simple []byte instead of bytes.Buffer to avoid large dependency.
type buffer []byte

func (bp *buffer) Write(p []byte) (int, error) {
	*bp = append(*bp, p...)
	return len(p), nil
}
func (bp *buffer) WriteString(s string) {
	*bp = append(*bp, s...)
}
func (bp *buffer) WriteQuote(s string) {
	*bp = strconv.AppendQuote(*bp, s)
}
func (bp *buffer) WriteByte(c byte) {
	*bp = append(*bp, c)
}
func (bp *buffer) WriteRune(r rune) {
	if r < utf8.RuneSelf {
		*bp = append(*bp, byte(r))
		return
	}

	b := *bp
	n := len(b)
	for n+utf8.UTFMax > cap(b) {
		b = append(b, 0)
	}
	w := utf8.EncodeRune(b[n:n+utf8.UTFMax], r)
	*bp = b[:n+w]
}
func (bp *buffer) Free() {
	buffers.Put(bp)
}
func (bp *buffer) Reset() {
	*bp = (*bp)[:0]
}

var buffers = sync.Pool{New: func() interface{} { return make([]byte, 0, 1<<10) }}
