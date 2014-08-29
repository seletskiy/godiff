package godiff

import (
	"bytes"
	"text/template"
)

type Changeset struct {
	FromHash   string
	ToHash     string
	Path       string
	Whitespace string
	Diffs      []*Diff
	Errors     []OperationError
}

type OperationError struct {
	Message string
}

func (o OperationError) Error() string {
	return o.Message
}

var changesetTpl = template.Must(template.New("file").Parse(
	"{{with $parent := .}}" +
		"{{range .Diffs}}" +
		"--- {{$parent.Path}}\t{{$parent.FromHash}}\n" +
		"+++ {{$parent.Path}}\t{{$parent.ToHash}}\n" +
		"{{.}}" +
		"{{else}}" +
		"{{end}}" +
		"{{end}}"))

func (r Changeset) String() string {
	buf := bytes.NewBuffer([]byte{})
	changesetTpl.Execute(buf, r)

	return buf.String()
}

func (r Changeset) ForEachComment(callback func(*Diff, *Comment, *Comment)) {
	for _, diff := range r.Diffs {
		stack := make([]*Comment, 0)
		parents := make(map[*Comment]*Comment)
		stack = append(stack, diff.LineComments...)
		pos := 0

		for pos < len(stack) {
			comment := stack[pos]

			if comment.Comments != nil {
				stack = append(stack, comment.Comments...)
				for _, c := range comment.Comments {
					parents[c] = comment
				}
			}

			callback(diff, comment, parents[comment])

			pos++
		}
	}
}

func (r Changeset) ForEachLine(callback func(*Diff, *Line)) {
	for _, diff := range r.Diffs {
		for _, hunk := range diff.Hunks {
			for _, segment := range hunk.Segments {
				for _, line := range segment.Lines {
					callback(diff, line)
				}
			}
		}
	}
}
