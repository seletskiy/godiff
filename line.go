package godiff

import "regexp"

type Line struct {
	Destination    int64
	Source         int64
	Line           string
	Truncated      bool
	ConflictMarker string
	CommentIds     []int64
	Comments       CommentsTree
}

var danglingSpacesRe = regexp.MustCompile("(?m) +$")

var lineTpl = loadSparseTemplate("line", `
{{.Line}}

{{if .Comments}}
	{{"\n"}}
	{{.Comments}}
{{end}}
`)

func (l Line) String() string {
	return lineTpl.Execute(l)
}
