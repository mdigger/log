package log_test

import (
	"errors"
	"os"
	"time"

	"github.com/mdigger/log"
)

func Example() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFlags(log.Lindent | log.Lshortfile)
	log.Padding = 18

	logger := log.WithField("service", "service name")
	logger.Debug("service started")
	defer logger.Debug("service stoped")

	ctx := logger.WithFields(log.Fields{
		"file":     "README.md",
		"type":     "text/markdown",
		"temp":     true,
		"duration": time.Second * 3,
	})
	ctx.Infof("info %v", "message")
	logger.Warning("warning message")
	logger.WithError(errors.New("error")).Error("error message")
	// Output:
	// example_test.go:18 ▸ service started    service="service name"
	// example_test.go:27 • info message       duration=3s file=README.md service="service name" temp=true type="text/markdown"
	// example_test.go:28 ⚡︎ warning message    service="service name"
	// example_test.go:29 ⨯︎ error message      error=error service="service name"
	// example_test.go:36 ▸ service stoped     service="service name"
}

// Errors are passed to WithError(), populating the "error" field.
func ExampleWithError() {
	err := errors.New("boom")
	log.WithError(err).Error("upload failed")
}

// Multiple fields can be set, via chaining, or WithFields().
func ExampleWithFields() {
	log.WithFields(log.Fields{
		"user": "Tobi",
		"file": "sloth.png",
		"type": "image/png",
	}).Info("upload")
}

// Structured logging is supported with fields, and is recommended over the
// formatted message variants.
func ExampleWithField() {
	log.WithField("user", "Tobo").Info("logged in")
}

// Unstructured logging is supported, but not recommended since it is hard to
// query.
func ExampleInfof() {
	log.Infof("%s logged in", "Tobi")
}

const filename = "README.md"

// Trace can be used to simplify logging of start and completion events, for
// example an upload or open which may fail.
func ExampleContext_Trace() (err error) {
	defer log.WithField("file", filename).Trace("open").Stop(&err)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	file.Close()
	return nil
}

// Creating a log for output in several formats.
func ExampleNew() {
	log := log.New(
		log.NewConsole(os.Stderr, log.LstdFlags|log.Lindent),
		log.NewJSON(os.Stderr, log.LstdFlags),
	)
	log.Info("multiple output log started")
	// Output:
}

// Creating a new console log.
func ExampleConsole() {
	handler := log.NewConsole(os.Stderr, log.LstdFlags|log.Lindent)
	handler.SetLevel(log.DebugLevel)
	log := handler.Context()
	log.Info("log started")
}

// Creating a new JSON log.
func ExampleJSON() {
	handler := log.NewJSON(os.Stderr, log.LstdFlags|log.Lindent)
	handler.SetLevel(log.DebugLevel)
	log := handler.Context()
	log.Info("log started")
}
