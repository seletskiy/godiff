package godiff

import (
	"bytes"
	"regexp"
	"text/template"
)

type Template struct {
	*template.Template
}

var reInsignificantWhitespace = regexp.MustCompile(`(?m)\n?^\s*`)

func loadSparseTemplate(name, text string) *Template {
	stripped := reInsignificantWhitespace.ReplaceAllString(text, ``)
	return &Template{
		template.Must(template.New("comment").Parse(stripped)),
	}
}

func (t *Template) Execute(v interface{}) string {
	buf := &bytes.Buffer{}
	t.Template.Execute(buf, v)

	return buf.String()
}
