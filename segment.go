package godiff

const (
	SegmentTypeContext = "CONTEXT"
	SegmentTypeRemoved = "REMOVED"
	SegmentTypeAdded   = "ADDED"
)

type Segment struct {
	Type      string
	Truncated bool
	Lines     []*Line
}

func (s Segment) String() string {
	result := ""
	for _, line := range s.Lines {
		operation := "?"
		switch s.Type {
		case SegmentTypeAdded:
			operation = "+"
		case SegmentTypeRemoved:
			operation = "-"
		case SegmentTypeContext:
			operation = " "
		}

		result += operation + line.String() + "\n"
	}

	return result
}

func (s Segment) GetLineNum(l *Line) int64 {
	switch s.Type {
	case SegmentTypeContext:
		fallthrough
	case SegmentTypeRemoved:
		return l.Source
	case SegmentTypeAdded:
		return l.Destination
	}

	return 0
}
