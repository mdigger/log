package log

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestWriter(t *testing.T) {
	w := NewWriter(os.Stderr, DEBUG, &Console{
		TimeFormat: "15:04:05",
	})
	log := w.New("test", "id", 4)
	log.With("a", "b").Info("info message")
	log.With(Fields{
		"name":  "name",
		"bool":  true,
		"int":   5,
		"array": []string{"1", "2", "3"},
		"":      "aga",
	},
		"aaa", "bbb",
		complex(3, 15), time.Now(),
	).Info("test message")
}

func TestJSON(t *testing.T) {
	w := NewWriter(os.Stderr, DEBUG, new(JSON))
	log := w.New("test", "id", 4)
	log.With("a", "b").Warn("info message")
}

func TestWriterColor(t *testing.T) {
	w := NewWriter(os.Stderr, DEBUG, &Color{KeyIndent: 8})
	log := w.New("test", "id", 4)
	log.With("a", "b").Info("info message")
	log.With(Fields{
		"name":  "name",
		"bool":  true,
		"int":   5,
		"array": []string{"1", "2", "3"},
		"":      "aga",
	},
		"aaa", "bbb",
		complex(3, 15), time.Now(),
	).Error("test message")
	w.Warn("warn")
	w.Trace("trace")
	w.Debug("debug")
}

func TestWriterErrors(t *testing.T) {
	w := NewWriter(os.Stderr, DEBUG, &Color{KeyIndent: 8})
	w.Info("info", errors.New("simple error"))
	err := fmt.Errorf("fmt error")
	w.Error("error", err)
	w.Error("error 2", errors.New(err.Error()))
}

// func TestDefault(t *testing.T) {
// 	SetFormat(&Console{WithSource: true})
// 	New("test", "a", "a").With("a", "b").Info("info message")
// }

// func TestStd(t *testing.T) {
// 	SetFormat(&Console{WithSource: true})
// 	log := StdLog(INFO, "test", "a", "b")
// 	log.Print("std message")
// 	New("aaa", "1", "2").StdLog(DEBUG).Print("test message")
// }
