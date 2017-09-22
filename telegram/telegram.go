package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mdigger/log"
)

const (
	botAPIURL = "https://api.telegram.org/"
	agent     = "mdigger-telegram-bot-log/1.0"
)

var httpClient = &http.Client{
	Timeout: time.Second * 30,
}

// Telegram описывает бота для отправки сообщений об ошибках в чат Telegram.
type Telegram struct {
	token     string
	chatID    int64
	*Template        // шаблон для формирования сообщения
	Header    string // заголовок
	Footer    string // подвал сообщения
}

// New создает бота для отправки уведомлений об ошибках в чат Telegram.
//
// Для отправки сообщений в Telegram необходимо указать токен бота (token),
// идентификатор чата (chatID) и шаблон для формирования сообщения
// (TelegramTemplate). Если шаблон не указан, то используется шаблон по
// умолчанию.
func New(token string, chatID int64, tmplt *Template) *Telegram {
	if tmplt == nil {
		tmplt = &Template{
			tmpl:   defaultTemplate,
			format: "Markdown",
		}
	}
	return &Telegram{token: token, chatID: chatID, Template: tmplt}
}

// Log реагирует ТОЛЬКО на сообщения с уровнем ERROR и выше. Все остальные
// записи игнорируются.
func (t *Telegram) Log(lvl log.Level, category, msg string, fields ...interface{}) error {
	if lvl < log.ERROR {
		return nil // игнорируем все, кроме ошибок
	}
	// формируем текст сообщения на основании шаблона
	var entry = &struct {
		Category  string
		Message   string
		Fields    map[string]string
		CallStack []*log.SourceInfo
		Header    string
		Footer    string
	}{
		Category:  category,
		Message:   msg,
		Fields:    nil,
		CallStack: log.CallStack(1),
		Header:    t.Header,
		Footer:    t.Footer,
	}

	switch len(fields) {
	case 0: // нет дополнительных полей
		break
	case 1: // дополнительные поля представлены одним элементом
		if list, ok := fields[0].(map[string]interface{}); ok {
			entry.Fields = make(map[string]string, len(list))
			for name, value := range list {
				entry.Fields[name] = fmt.Sprint(value)
			}
		}
	default:
		entry.Fields = make(map[string]string, len(fields)>>1)
		var name string
		for i, field := range fields {
			if i%2 == 0 {
				name = fmt.Sprint(field)
			} else {
				entry.Fields[name] = fmt.Sprint(field)
			}
		}
	}
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, entry); err != nil {
		return err
	}
	// отправляем на Telegram
	return t.Send(buf.String(), t.format)
}

// Send позволяет отправить в Telegram произвольный текст. Параметр format
// задает форматирование текста сообщения и может быть либо "Markdown", либо
// "HTML".
func (t *Telegram) Send(text, format string) error {
	var params = url.Values{}
	params.Set("chat_id", strconv.FormatInt(t.chatID, 10))
	params.Set("text", text)
	if format == "Markdown" || format == "HTML" {
		params.Set("parse_mode", format)
	}
	var apiURL = botAPIURL + "bot" + t.token + "/sendMessage"
	req, err := http.NewRequest("POST", apiURL,
		strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", agent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return nil
	}
	var telegramError = new(struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	})
	if err = json.NewDecoder(resp.Body).Decode(telegramError); err == nil {
		return errors.New(telegramError.Description)
	}
	return errors.New(resp.Status)
}
