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
// 	Message   string
// 	Fields    log.Fields
// 	CallStack []*log.SourceInfo
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
	"{{if .Header}}{{.Header}}\n{{end}}" +
		"{{if .Category}}*[{{.Category}}]*: {{end}}{{.Message}}\n\n" +
		"{{if .Fields}}{{range $name, $value := .Fields}}_{{$name}}_: " +
		"{{$value}}\n{{end}}\n{{end}}{{range $value := .CallStack}}" +
		"- `{{$value}}`\n{{end}}" +
		"{{if .Footer}}\n{{.Footer}}\n{{end}}"))
