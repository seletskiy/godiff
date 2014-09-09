package godiff

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
	Hunks []*Hunk

	LineComments CommentsTree
	DiffComments CommentsTree

	Note string

	// Lists made only for Stash API compatibility.
	Attributes struct {
		FromHash []string
		ToHash   []string
	}
}

func (d Diff) GetHashFrom() string {
	if len(d.Attributes.FromHash) > 0 {
		return d.Attributes.FromHash[0]
	} else {
		return "???"
	}
}

func (d Diff) GetHashTo() string {
	if len(d.Attributes.ToHash) > 0 {
		return d.Attributes.ToHash[0]
	} else {
		return "???"
	}
}

func (d Diff) ForEachLine(callback func(*Diff, *Hunk, *Segment, *Line)) {
	for _, hunk := range d.Hunks {
		for _, segment := range hunk.Segments {
			for _, line := range segment.Lines {
				callback(&d, hunk, segment, line)
			}
		}
	}
}
