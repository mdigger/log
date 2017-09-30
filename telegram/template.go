package telegram

import "html/template"

// Template описывает шаблон для формирования сообщения для Telegram.
type Template struct {
	format string
	tmpl   *template.Template
}

// NewTemplate разбирает шаблон и возвращает его в разобранном виде.
// Формат текст задается в параметре format и может быть либо "Markdown", либо
// "HTML". В шаблоне доступны следующие поля:
// 	Category  string
// 	Level     string
// 	Message   string
// 	Fields    log.Fields
// 	Header    string
// 	Footer    string
func NewTemplate(text, format string) (*Template, error) {
	t, err := template.New("").Parse(text)
	if err != nil {
		return nil, err
	}
	return &Template{format: format, tmpl: t}, nil
}

var defaultTemplate = template.Must(template.New("").Parse(
	`{{if .Header}}{{.Header}}
{{end}}[<b>{{.Level}}</b>] {{if .Category}}{{.Category}}: {{end}}{{.Message}}
{{- range .Fields}}
<i>{{.Name}}:</i>	{{.Value}}{{end}}
{{- if .Footer}}
{{.Footer}}{{end}}`))
