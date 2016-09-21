package log_test

import (
	"os"
	"time"

	"github.com/mdigger/log"
)

func Example_console() {
	// new log output to the console
	logger := log.New(log.NewPlainHandler(os.Stdout, 0))

	logger.Info("info test")
	logger.Infof("info %v", "test")

	logger.Error("error test")
	logger.Errorf("error %v", "test")

	logger.Debug("debug test")
	logger.Debugf("debug %v", "test")

	// an informational message with additional parameters
	logger.WithField("key", "value").Info("test")
	logger.WithFields(log.Fields{
		"key":  "value",
		"key2": "value2",
	}).Info("test")

	// to be added to the file name and line number of the source code
	logger.WithField("key", "value").WithSource(0).Info("test")

	// Output:
	// info test
	// info test
	// error: error test
	// error: error test
	// test                         key=value
	// test                         key=value key2=value2
	// test                         key=value source="example_test.go:31"
}

func Example_json() {
	logger := log.New(log.NewJSONHandler(os.Stdout, 0))

	logger.Info("info test")
	logger.Infof("info %v", "test")

	logger.Error("error test")
	logger.Errorf("error %v", "test")

	logger.Debug("debug test")
	logger.Debugf("debug %v", "test")

	// an informational message with additional parameters
	logger.WithField("key", "value").Info("test")
	logger.WithFields(log.Fields{
		"key":  "value",
		"key2": "value2",
	}).Info("test")

	// to be added to the file name and line number of the source code
	logger.WithField("key", "value").WithSource(0).Info("test")

	// Output:
	// {"level":"info","message":"info test"}
	// {"level":"info","message":"info test"}
	// {"level":"error","message":"error test"}
	// {"level":"error","message":"error test"}
	// {"level":"info","message":"test","fields":{"key":"value"}}
	// {"level":"info","message":"test","fields":{"key":"value","key2":"value2"}}
	// {"level":"info","message":"test","fields":{"key":"value","source":"example_test.go:63"}}
}

func Example_mixed() {
	// new log output to the console
	clog := log.NewPlainHandler(os.Stdout, 0)
	// new log output to the console in JSON format
	json := log.NewJSONHandler(os.Stdout, 0)
	json.SetFlags(0)
	// output to multiple logs in different formats
	logger := log.New(clog, json)
	logger.Info("info")

	// Output:
	// info
	// {"level":"info","message":"info"}
}

func Example() {
	log.Info("info message")
	log.WithField("time", time.Now()).Debug("debug")

	var err error
	filename := "README.md"
	// to form the log at the beginning of the open file and at the end add in
	// the description of the error if it happens
	defer log.WithField("file", filename).Trace("open").Stop(&err)
	_, err = os.Open(filename)
}
