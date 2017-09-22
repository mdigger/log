package log

import (
	"errors"
	"os"
	"testing"
	"time"
)

var errTest = errors.New("test error")

func TestConsole(t *testing.T) {
	for _, h := range []StreamFormatter{new(Console), Color, JSON("  ")} {
		err := h.Log(os.Stderr, -65, "category", "message",
			"key1", "s",
			"key2", 2,
			"key3", true,
			"key4", errTest,
			"key5", time.Now(),
			"key6", []string{"s1", "s2"},
			"key7", []int{1, 2},
			"key8", complex(12, 53),
			"key9", []byte{'t', 'e', 's', 't'},
			"key10", map[string]interface{}{
				"k1": 1,
				"k2": "2",
				"k3": true,
			},
		)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestColor(t *testing.T) {
	SetFormat(Color)
	SetLevel(TRACE)
	fields := Fields{"test": "string", "bool": true}
	log := New("category")
	log.Trace("trace", fields)
	log.Debug("debug", fields)
	log.Info("info", fields)
	log.Error("error", fields)
	log.Fatal("fatal", fields)
}
