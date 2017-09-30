package telegram

import (
	"testing"

	"github.com/mdigger/errors"
	"github.com/mdigger/log"
)

var (
	token        = "422160011:AAFz-BJhIFQLrdXI2L8BtxgvivDKeY5s2Ig"
	chatID int64 = -1001068031302
)

func TestTelegram(t *testing.T) {
	bot := New(token, chatID, nil)
	bot.Header = "header"
	bot.Footer = "footer"
	if err := bot.Write(log.ERROR, 1, "test", "message", []log.Field{
		{"key1", 1},
		{"key2", true},
		{"key3", errors.New("error")},
		{"key4", []byte{'t', 'e', 's', 't'}},
		{"key5", errors.New(errors.New("error"))},
	}); err != nil {
		t.Fatal(err)
	}
}
