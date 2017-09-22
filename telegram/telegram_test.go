package telegram

import (
	"errors"
	"testing"

	"github.com/mdigger/log"
)

var (
)

func TestTelegram(t *testing.T) {
	bot := New(token, chatID, nil)
	// bot.Header = "header"
	// bot.Footer = "footer"
	if err := bot.Log(log.ERROR, "test", "message",
		"key1", 1,
		"key2", true,
		"key3", errors.New("error"),
		"key4", []byte{'t', 'e', 's', 't'}); err != nil {
		t.Fatal(err)
	}
}
