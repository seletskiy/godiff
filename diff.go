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
	Hunks        []*Hunk
	LineComments []*Comment

	// Lists made only for Stash API compatibility.
	Attributes struct {
		FromHash []string
		ToHash   []string
	}
}

var diffTpl = loadSparseTemplate(`diff`, `
--- {{.Source.ToString}}{{"\t"}}{{index .Attributes.FromHash 0}}{{"\n"}}
+++ {{.Destination.ToString}}{{"\t"}}{{index .Attributes.ToHash 0}}{{"\n"}}
{{range .Hunks}}
	{{.}}
{{end}}
`)

func (d Diff) String() string {
	return diffTpl.Execute(d)
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
