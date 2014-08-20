package godiff

import (
	"bytes"
	"text/template"
)

type Diff struct {
	Truncated bool
	Source    struct {
		Parent   string
		Name     string
		ToString string
	}
	Destination struct {
		Parent   string
		Name     string
		ToString string
	}
	Hunks        []*Hunk
	LineComments []*Comment
}

var diffTpl = template.Must(template.New("diff").Parse(
	"{{range .Hunks}}{{.}}{{end}}"))

func (d Diff) String() string {
	buf := bytes.NewBuffer([]byte{})
	diffTpl.Execute(buf, d)

	return buf.String()
}
