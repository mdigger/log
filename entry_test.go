package log_test

import (
	"os"
	"testing"
	"time"

	"github.com/mdigger/log"
)

func TestTrace(t *testing.T) {
	log := log.New(log.NewConsoleHandler(os.Stdout, log.LstdFlags))
	open := func(filename string) (err error) {
		defer log.WithField("filename", filename).Trace("open").Stop(&err)
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
		file.Close()
		return nil
	}
	open("README.md")
	open("entry_test.go")
}
