package godiff

type Hunk struct {
	SourceLine      int64
	SourceSpan      int64
	DestinationLine int64
	DestinationSpan int64
	Truncated       bool
	Segments        []*Segment
}
