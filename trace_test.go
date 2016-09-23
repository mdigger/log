package log

import (
	"os"
	"testing"
)

func TestTrace(t *testing.T) {
	open := func(filename string) (err error) {
		defer Tracef("trace: %v", filename).Stop(&err)
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}
	open2 := func(filename string) (err error) {
		defer Trace("open file").WithFields(Fields{"test": true}).Stop(&err)
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}
	open3 := func(filename string) (err error) {
		defer WithField("file", filename).Trace("open").Stop(&err)
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}
	open4 := func(filename string) (err error) {
		defer WithField("file", filename).Tracef("open %v", filename).Stop(&err)
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}

	for _, name := range []string{"trace_test.go", "bad_file.txt"} {
		_ = open(name)
		_ = open2(name)
		_ = open3(name)
		_ = open4(name)
	}
}
