package godiff

import (
	"regexp"
	"time"
)

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
	Comments    CommentsTree
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

type CommentsTree []*Comment

//const replyIndent = "    "

var begOfLineRe = regexp.MustCompile("(?m)^")

//func (c Comment) String() string {
//    comments, _ := commentTpl.Execute(c)

//    for _, reply := range c.Comments {
//        comments += reply.AsReply()
//    }

//    return comments
//}

//func (c Comment) AsReply() string {
//    return begOfLineRe.ReplaceAllString(
//        commentSpacing+c.String(),
//        replyIndent,
//    )
//}

var reWhiteSpace = regexp.MustCompile(`\s+`)

func (c Comment) Short(length int) string {
	sticked := []rune(reWhiteSpace.ReplaceAllString(c.Text, " "))

	if len(sticked) > length {
		return string(sticked[:length]) + "..."
	} else {
		return string(sticked)
	}
}

const ignorePrefix = "###"

var reBeginningOfLine = regexp.MustCompile(`(?m)^`)
var reIgnorePrefixSpace = regexp.MustCompile("(?m)^" + ignorePrefix + " $")

func Note(String string) string {
	return reIgnorePrefixSpace.ReplaceAllString(
		reBeginningOfLine.ReplaceAllString(String, ignorePrefix+" "),
		ignorePrefix)
}

//const commentSpacing = "\n\n"
//const commentPrefix = "# "

//func (comments CommentsTree) String() string {
//    res := ""

//    if len(comments) > 0 {
//        res = "---" + commentSpacing
//    }

//    for i, comment := range comments {
//        res += comment.String()
//        if i < len(comments)-1 {
//            res += commentSpacing
//        }
//    }

//    if len(comments) > 0 {
//        return danglingSpacesRe.ReplaceAllString(
//            begOfLineRe.ReplaceAllString(res, "# "), "")
//    } else {
//        return ""
//    }
//}
