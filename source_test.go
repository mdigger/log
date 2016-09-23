package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestSource(t *testing.T) {
	source := NewSource(0)
	if source.String() != "source_test.go:13" {
		t.Error("bad source")
	}
	if err := json.NewEncoder(ioutil.Discard).Encode(source); err != nil {
		t.Error(err)
	}

	if err := checkSource(WithSource(0), 21); err != nil {
		t.Error(err)
	}
	if err := checkSource(WithField("name", "value").WithSource(0), 24); err != nil {
		t.Error(err)
	}
}

func checkSource(c *Context, line int) error {
	src, ok := c.Fields["source"]
	if !ok || src == nil {
		return errors.New("empty source")
	}
	source, ok := src.(*Source)
	if !ok || source == nil {
		return errors.New("bad source")
	}
	var filename = "source_test.go"
	if file := filepath.Base(source.File); file != filename {
		return fmt.Errorf("bad source file name: %v", file)
	}
	if source.Line != line {
		return fmt.Errorf("bad source file line: %v", source.Line)
	}
	return nil
}
