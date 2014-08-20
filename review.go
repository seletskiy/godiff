package godiff

import (
	"bytes"
	"text/template"
)

type Review struct {
	FromHash   string
	ToHash     string
	Path       string
	Whitespace string
	Diffs      []*Diff
}

func (r Review) ForEachLines(callback func(*Diff, *Line)) {
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

var reviewTpl = template.Must(template.New("file").Parse(
	"{{with $parent := .}}" +
		"{{range .Diffs}}" +
		"--- {{$parent.Path}}\t{{$parent.FromHash}}\n" +
		"+++ {{$parent.Path}}\t{{$parent.ToHash}}\n" +
		"{{.}}" +
		"{{else}}" +
		"{{end}}" +
		"{{end}}"))

func (r Review) String() string {
	buf := bytes.NewBuffer([]byte{})
	reviewTpl.Execute(buf, r)

	return buf.String()
}

func (r Review) ForEachComment(callback func(*Diff, *Comment, *Comment)) {
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

func (current Review) Compare(another Review) []map[string]interface{} {
	existComments := make([]*Comment, 0)

	current.ForEachComment(func(_ *Diff, comment, _ *Comment) {
		existComments = append(existComments, comment)
	})

	changeset := make([]map[string]interface{}, 0)

	another.ForEachComment(func(_ *Diff, comment, parent *Comment) {
		if comment.Id == 0 {
			if parent != nil {
				changeset = append(changeset, map[string]interface{}{
					"text": comment.Text,
					"parent": map[string]interface{}{
						"id": parent.Id,
					},
				})
			} else {
				changeset = append(changeset,
					map[string]interface{}{
						"text": comment.Text,
						"anchor": map[string]interface{}{
							"line":     comment.Anchor.Line,
							"lineType": comment.Anchor.LineType,
							"fromFile": comment.Anchor.SrcPath,
						},
					})
			}
		} else {
			for i, c := range existComments {
				if c != nil && c.Id == comment.Id {
					existComments[i] = nil
					if c.Text != comment.Text {
						changeset = append(changeset,
							map[string]interface{}{
								"text": comment.Text,
								"id":   comment.Id,
							})
					}
				}
			}
		}
	})

	for _, deleted := range existComments {
		if deleted != nil {
			changeset = append(changeset, map[string]interface{}{"id": deleted.Id})
		}
	}

	return changeset
}
