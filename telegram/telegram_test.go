package telegram

import (
	"errors"
	"testing"

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
	if err := bot.Log(log.ERROR, "test", "message",
		"key1", 1,
		"key2", true,
		"key3", errors.New("error"),
		"key4", []byte{'t', 'e', 's', 't'}); err != nil {
		t.Fatal(err)
	}
}
