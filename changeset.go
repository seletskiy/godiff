package godiff

type Changeset struct {
	FromHash   string
	ToHash     string
	Path       string
	Whitespace string
	Diffs      []*Diff
}

func (r Changeset) ForEachComment(callback func(*Diff, *Comment, *Comment)) {
	for _, diff := range r.Diffs {
		stack := make([]*Comment, 0)
		parents := make(map[*Comment]*Comment)
		stack = append(stack, diff.FileComments...)
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

func (r Changeset) ForEachLine(
	callback func(*Diff, *Hunk, *Segment, *Line) error,
) error {
	for _, diff := range r.Diffs {
		err := diff.ForEachLine(callback)
		if err != nil {
			return err
		}
	}

	return nil
}
