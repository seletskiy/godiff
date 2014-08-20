package godiff

import (
	"bytes"
	"regexp"
	"text/template"
	"time"
)

const commentShortLength = 20

var reWhiteSpace = regexp.MustCompile(`\s+`)

type UnixTimestamp int

func (u UnixTimestamp) String() string {
	return time.Unix(int64(u/1000), 0).Format(time.ANSIC)
}

type Comment struct {
	Id          int64
	Version     int
	Text        string
	CreatedDate UnixTimestamp
	UpdatedDate UnixTimestamp
	Comments    []*Comment
	Author      struct {
		Name         string
		EmailAddress string
		Id           int
		DisplayName  string
		Active       bool
		Slug         string
		Type         string
	}

	Anchor CommentAnchor

	PermittedOperations struct {
		Editable  bool
		Deletable bool
	}

	Indent   int
	Parented bool
}

type CommentAnchor struct {
	FromHash string
	ToHash   string
	Line     int64  `json:"line"`
	LineType string `json:"lineType"`
	Path     string `json:"path"`
	SrcPath  string `json:"srcPath"`
	FileType string `json:"fileType"`
}

const replyIndent = "    "

var commentTpl = template.Must(template.New("comment").Parse(
	"\n\n" +
		"[{{.Id}}] | {{.Author.DisplayName}} | {{.UpdatedDate}}\n" +
		"\n" +
		"{{.Text}}\n" +
		"\n---"))

func (c Comment) String() string {
	buf := bytes.NewBuffer([]byte{})
	commentTpl.Execute(buf, c)

	for _, reply := range c.Comments {
		buf.WriteString(
			begOfLineRe.ReplaceAllString(reply.String(), "\n"+replyIndent))
	}

	return buf.String()
}

func (c Comment) Short(length int) string {
	sticked := []rune(reWhiteSpace.ReplaceAllString(c.Text, " "))

	if len(sticked) > length {
		return string(sticked[:length]) + "..."
	} else {
		return string(sticked)
	}
}
