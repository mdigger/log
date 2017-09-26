package log

import (
	"os"
	"testing"
)

func TestWriter(t *testing.T) {
	w := NewWriter(os.Stderr, DEBUG, &Console{
		TimeFormat: "15:04:03",
		WithSource: true,
	})
	log := w.New("test", "id", 4)
	log.With("a", "b").Info("info message")
}

func TestDefault(t *testing.T) {
	SetFormat(&Console{WithSource: true})
	New("test", "a", "a").With("a", "b").Info("info message")
}

func TestStd(t *testing.T) {
	SetFormat(&Console{WithSource: true})
	log := StdLog(INFO, "test", "a", "b")
	log.Print("std message")
	New("aaa", "1", "2").StdLog(DEBUG).Print("test message")
}

func TestJSON(t *testing.T) {
	w := NewWriter(os.Stderr, DEBUG, new(JSON))
	log := w.New("test", "id", 4)
	log.With("a", "b").Warn("info message")
}
