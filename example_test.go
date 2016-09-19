package log_test

import (
	"os"
	"time"

	"github.com/mdigger/log"
)

func Example_console() {
	// new log output to the console
	logger := log.New(log.NewPlainHandler(os.Stdout, log.Lshortfile))

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
	logger.WithSource(0).Info("info test")
	logger.WithField("key", "value").WithSource(0).Info("test")

	// Output:
	// example_test.go:14 info test
	// example_test.go:15 info test
	// example_test.go:17 error: error test
	// example_test.go:18 error: error test
	// example_test.go:24 test key=value
	// example_test.go:28 test key=value key2=value2
	// example_test.go:31 info test source="example_test.go:31"
	// example_test.go:32 test key=value source="example_test.go:32"
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
	logger.WithSource(0).Info("info test")
	logger.WithField("key", "value").WithSource(0).Info("test")

	// Output:
	// {"level":"info","message":"info test"}
	// {"level":"info","message":"info test"}
	// {"level":"error","message":"error test"}
	// {"level":"error","message":"error test"}
	// {"level":"info","message":"test","fields":{"key":"value"}}
	// {"level":"info","message":"test","fields":{"key":"value","key2":"value2"}}
	// {"level":"info","message":"info test","fields":{"source":"example_test.go:65"}}
	// {"level":"info","message":"test","fields":{"key":"value","source":"example_test.go:66"}}
}

func Example_mixed() {
	// new log output to the console
	clog := log.NewPlainHandler(os.Stdout, log.Lshortfile)
	// new log output to the console in JSON format
	json := log.NewJSONHandler(os.Stdout, 0)
	json.SetFlags(0)
	// output to multiple logs in different formats
	logger := log.New(clog, json)
	logger.Info("info")

	// Output:
	// example_test.go:87 info
	// {"level":"info","message":"info"}
}

func ExampleTrace() {
	// create a handler for a console log
	clog := log.NewPlainHandler(os.Stdout, log.Lshortfile)
	logger := log.New(clog) // инициализируем лог

	var err error
	filename := "README.md"
	// to form the log at the beginning of the open file and at the end add in
	// the description of the error if it happens
	logger.WithField("file", filename).Trace("open").Stop(&err)
	_, err = os.Open(filename)
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
