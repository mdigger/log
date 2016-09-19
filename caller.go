package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var root, _ = os.Getwd()

// Caller describes a source file name and line number from where the call
// originated.
type Caller struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

// String returns the file name and line number in a text view, separated by a
// colon.
func (c *Caller) String() string {
	file := c.File
	if file == "" {
		return "<unknown>"
	} else if relfile, err := filepath.Rel(root, file); err == nil {
		file = relfile
	}
	return fmt.Sprintf("%s:%d", file, c.Line)
}

// MarshalText used when serializing Caller in text format.
func (c *Caller) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// MakeCaller calls the function get the name of the file and line of source
// code that generated this feature. As the parameter indicates the depth in
// the stack.
func MakeCaller(calldepth int) *Caller {
	_, file, line, _ := runtime.Caller(1 + calldepth)
	return &Caller{
		File: file,
		Line: line,
	}
}
