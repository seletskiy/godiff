package godiff

import (
	"io"
	"regexp"
	"text/template"

	"github.com/seletskiy/tplutil"
)

var changesetTpl = template.New("changeset")

func init() {
	commentsTpl := template.New("comment")

	reBeginningOfLine := regexp.MustCompile(`(?m)^`)
	reNewLine := regexp.MustCompile(`^|\n`)
	reDanglingSpace := regexp.MustCompile(`(?m)\s+$`)

	funcs := template.FuncMap{
		"indent": func(input string) string {
			return reBeginningOfLine.ReplaceAllString(input, "    ")
		},
		"writeComments": func(input CommentsTree) string {
			res, _ := tplutil.ExecuteToString(commentsTpl, input)
			return res
		},
		"comment": func(input string) string {
			return reDanglingSpace.ReplaceAllString(
				reNewLine.ReplaceAllString(input, `$0# `),
				``)
		},
	}

	template.Must(
		commentsTpl.Funcs(funcs).Funcs(tplutil.Last).Parse(
			tplutil.Strip(commentsTplText)))
	template.Must(
		changesetTpl.Funcs(funcs).Funcs(tplutil.Last).Parse(
			tplutil.Strip(changesetTplText)))
}

func WriteChangeset(changeset Changeset, to io.Writer) error {
	return changesetTpl.Execute(to, changeset)
}
