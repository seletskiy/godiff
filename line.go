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

//var lineTpl = tplutil.SparseTemplate("line", `
//{{.Line}}

//{{if .Comments}}
//    {{"\n"}}
//    {{.Comments}}
//{{end}}
//`)

//func (l Line) String() string {
//    result, _ := lineTpl.Execute(l)
//    return result
//}
