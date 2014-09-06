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

func (s Segment) TextPrefix() string {
	switch s.Type {
	case SegmentTypeAdded:
		return "+"
	case SegmentTypeRemoved:
		return "-"
	case SegmentTypeContext:
		return " "
	default:
		return "?"
	}
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
