package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Source describes information about the file name and line number of the
// source code.
type Source struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

var root, _ = os.Getwd()

// NewSource returns information about the file name and line number of the
// source code. Calldepth is the count of the number of frames to skip when
// computing the file name and line number. A value of 0 will print the details
// for the caller.
func NewSource(calldepth int) *Source {
	_, file, line, _ := runtime.Caller(1 + calldepth)
	return &Source{
		File: file,
		Line: line,
	}
}

// String implements io.Stringer.
func (s Source) String() string {
	file := s.File
	if file == "" {
		return "<unknown>"
	} else if relfile, err := filepath.Rel(root, file); err == nil {
		file = relfile
	}
	return fmt.Sprintf("%s:%d", file, s.Line)
}

// MarshalText returns the Source string.
func (s Source) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}
