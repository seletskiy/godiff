package godiff

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	stateStartOfFile   = "stateStartOfFile"
	stateDiffHeader    = "stateDiffHeader"
	stateHunkHeader    = "stateHunkHeader"
	stateHunkBody      = "stateHunkBody"
	stateComment       = "stateComment"
	stateCommentDelim  = "stateCommentDelim"
	stateCommentHeader = "stateCommentHeader"
)

var (
	reFromFile      = regexp.MustCompile(`^--- ([^ ]+)\t(.*)`)
	reToFile        = regexp.MustCompile(`^\+\+\+ ([^ ]+)\t(.*)`)
	reHunk          = regexp.MustCompile(`^@@ -(\d+),(\d+) \+(\d+),(\d+) @@`)
	reCommentDelim  = regexp.MustCompile(`^#\s+---`)
	reCommentHeader = regexp.MustCompile(`^#\s+\[(\d+)\]\s+\|([^|]+)\|(.*)`)
	reCommentText   = regexp.MustCompile(`^#\s*(.*)`)
	reIndent        = regexp.MustCompile(`^#(\s+)`)
)

type parser struct {
	state     string
	changeset Changeset
	diff      *Diff
	hunk      *Hunk
	segment   *Segment
	comment   *Comment
	line      *Line

	segmentType  string
	commentsList []*Comment
}

func ParseDiff(r io.Reader) (Changeset, error) {
	buffer := bufio.NewReader(r)

	current := parser{}
	current.state = stateStartOfFile

	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}

		current.switchState(line)
		current.createNodes(line)
		current.locateNodes(line)
		current.parseLine(line)
	}

	for _, comment := range current.commentsList {
		comment.Text = strings.TrimSpace(comment.Text)
	}

	return current.changeset, nil
}

func (current *parser) switchState(line string) error {
	inComment := false
	switch current.state {
	case stateStartOfFile:
		switch line[0] {
		case '-':
			current.state = stateDiffHeader
		}
	case stateDiffHeader:
		switch line[0] {
		case '@':
			current.state = stateHunkHeader
		}
	case stateHunkHeader:
		current.state = stateHunkBody
		fallthrough
	case stateHunkBody, stateComment, stateCommentDelim, stateCommentHeader:
		switch line[0] {
		case ' ':
			current.state = stateHunkBody
			current.segmentType = SegmentTypeContext
		case '-':
			current.state = stateHunkBody
			current.segmentType = SegmentTypeRemoved
		case '+':
			current.state = stateHunkBody
			current.segmentType = SegmentTypeAdded
		case '@':
			current.state = stateHunkHeader
		case '#':
			inComment = true
			switch {
			case reCommentDelim.MatchString(line):
				current.state = stateCommentDelim
			case reCommentHeader.MatchString(line):
				current.state = stateCommentHeader
			case reCommentText.MatchString(line):
				current.state = stateComment
			}
		}
	}

	if !inComment {
		current.comment = nil
	}

	return nil
}

func (current *parser) createNodes(line string) error {
	switch current.state {
	case stateDiffHeader, stateHunkHeader:
		switch line[0] {
		case '-':
			current.diff = &Diff{}
			current.changeset.Diffs = append(current.changeset.Diffs, current.diff)
		case '@':
			current.hunk = &Hunk{}
			current.segment = &Segment{}
		}
	case stateCommentDelim:
		current.comment = &Comment{}
	case stateComment:
		switch {
		case reCommentDelim.MatchString(line):
			// noop
		case reCommentText.MatchString(line):
			if current.comment == nil {
				current.comment = &Comment{}
			}
		}
	case stateHunkBody:
		if current.segment.Type != current.segmentType {
			current.segment = &Segment{Type: current.segmentType}
			current.hunk.Segments = append(current.hunk.Segments, current.segment)
		}

		current.line = &Line{}
		current.segment.Lines = append(current.segment.Lines, current.line)
	}

	return nil
}

func (current *parser) locateNodes(line string) error {
	switch current.state {
	case stateComment:
		current.locateComment(line)
	case stateHunkBody:
		current.locateLine(line)
	}

	return nil
}

func (current *parser) locateComment(line string) error {
	if current.comment.Parented || strings.TrimSpace(line) == "#" {
		return nil
	}

	current.commentsList = append(current.commentsList, current.comment)

	current.comment.Anchor.LineType = current.segment.Type
	switch current.segment.Type {
	case SegmentTypeContext:
		fallthrough
	case SegmentTypeRemoved:
		current.comment.Anchor.Line = current.line.Source
	case SegmentTypeAdded:
		current.comment.Anchor.Line = current.line.Destination
	}
	current.comment.Anchor.Path = current.diff.Destination.ToString
	current.comment.Anchor.SrcPath = current.diff.Source.ToString
	current.comment.Indent = getIndentSize(line)
	current.comment.Parented = true
	parent := current.findParentComment(current.comment)
	if parent != nil {
		parent.Comments = append(parent.Comments, current.comment)
	} else {
		current.line.Comments = append(current.line.Comments, current.comment)
		current.diff.LineComments = append(current.diff.LineComments, current.comment)
	}

	return nil
}

func (current *parser) locateLine(line string) error {
	sourceOffset := current.hunk.SourceLine - 1
	destinationOffset := current.hunk.DestinationLine - 1
	if len(current.hunk.Segments) > 1 {
		prevSegment := current.hunk.Segments[len(current.hunk.Segments)-2]
		lastLine := prevSegment.Lines[len(prevSegment.Lines)-1]
		sourceOffset = lastLine.Source
		destinationOffset = lastLine.Destination
	}
	hunkLength := int64(len(current.segment.Lines))
	switch current.segment.Type {
	case SegmentTypeContext:
		current.line.Source = sourceOffset + hunkLength
		current.line.Destination = destinationOffset + hunkLength
	case SegmentTypeAdded:
		current.line.Source = sourceOffset
		current.line.Destination = destinationOffset + hunkLength
	case SegmentTypeRemoved:
		current.line.Source = sourceOffset + hunkLength
		current.line.Destination = destinationOffset
	}

	return nil
}

func (current *parser) parseLine(line string) error {
	switch current.state {
	case stateDiffHeader:
		current.parseDiffHeader(line)
	case stateHunkHeader:
		current.parseHunkHeader(line)
	case stateHunkBody:
		current.parseHunkBody(line)
	case stateComment:
		current.parseComment(line)
	case stateCommentHeader:
		current.parseCommentHeader(line)
	}

	return nil
}

func (current *parser) parseDiffHeader(line string) error {
	switch {
	case reFromFile.MatchString(line):
		matches := reFromFile.FindStringSubmatch(line)
		current.changeset.Path = matches[1]
		current.diff.Source.ToString = matches[1]
		current.changeset.FromHash = matches[2]
	case reToFile.MatchString(line):
		matches := reToFile.FindStringSubmatch(line)
		current.diff.Destination.ToString = matches[1]
		current.changeset.ToHash = matches[2]
	default:
		return errors.New("expected diff header, but not found")
	}
	return nil
}

func (current *parser) parseHunkHeader(line string) error {
	matches := reHunk.FindStringSubmatch(line)
	current.hunk.SourceLine, _ = strconv.ParseInt(matches[1], 10, 16)
	current.hunk.SourceSpan, _ = strconv.ParseInt(matches[2], 10, 16)
	current.hunk.DestinationLine, _ = strconv.ParseInt(matches[3], 10, 16)
	current.hunk.DestinationSpan, _ = strconv.ParseInt(matches[4], 10, 16)
	current.diff.Hunks = append(current.diff.Hunks, current.hunk)

	return nil
}

func (current *parser) parseHunkBody(line string) error {
	current.line.Line = line[1 : len(line)-1]
	return nil
}

func (current *parser) parseCommentHeader(line string) error {
	matches := reCommentHeader.FindStringSubmatch(line)
	current.comment.Author.DisplayName = strings.TrimSpace(matches[2])
	current.comment.Id, _ = strconv.ParseInt(matches[1], 10, 16)
	updatedDate, _ := time.ParseInLocation(time.ANSIC,
		strings.TrimSpace(matches[3]),
		time.Local)
	current.comment.UpdatedDate = UnixTimestamp(updatedDate.Unix() * 1000)

	return nil
}

func (current *parser) parseComment(line string) error {
	matches := reCommentText.FindStringSubmatch(line)
	current.comment.Text += "\n" + strings.Trim(matches[1], " \t")

	return nil
}

func (current *parser) findParentComment(comment *Comment) *Comment {
	for i := len(current.commentsList) - 1; i >= 0; i-- {
		c := current.commentsList[i]
		if comment.Indent > c.Indent {
			return c
		}
	}

	return nil
}

func getIndentSize(line string) int {
	matches := reIndent.FindStringSubmatch(line)
	if len(matches) == 0 {
		return 0
	}

	return len(matches[1])
}
