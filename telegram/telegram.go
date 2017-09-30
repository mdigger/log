package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mdigger/errors"
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
// (Template). Если шаблон не указан, то используется шаблон по умолчанию.
func New(token string, chatID int64, tmplt *Template) *Telegram {
	if tmplt == nil {
		tmplt = &Template{
			tmpl:   defaultTemplate,
			format: "HTML",
		}
	}
	return &Telegram{token: token, chatID: chatID, Template: tmplt}
}

// Write отсылает лог в Telegram.
func (t *Telegram) Write(lvl log.Level, calldepth int, category, msg string, fields []log.Field) error {
	if lvl < log.INFO {
		return nil
	}
	// изменяем ошибку
	for i, field := range fields {
		if err, ok := field.Value.(error); ok {
			msg := err.Error()
			if err, ok := err.(*errors.Error); ok {
				if err := err.Cause(); err != nil {
					msg += "\n\t<i>cause:</i> " + err.Error()
					msg += fmt.Sprintf("\n\t<code>%#+v</code>", err)
				}
				for _, src := range err.Stacks() {
					msg += "\n\t- " + src.Func + " [" + src.String() + "] "
				}
			} else {
				msg += fmt.Sprintf("\n\t<code>%#+[1]v</code>", err)
			}
			field.Value = template.HTML(msg)
			fields[i] = field
		}
	}
	// формируем текст сообщения на основании шаблона
	var entry = &struct {
		*log.Entry
		Header string
		Footer string
	}{
		Entry:  log.NewEntry(lvl, category, msg, fields),
		Header: t.Header,
		Footer: t.Footer,
	}
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, entry); err != nil {
		return err
	}
	entry.Free()
	// fmt.Println(buf.String())
	// return nil
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
