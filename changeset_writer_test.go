package godiff

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestWillOmitDanglingWhitespaceOnRender(t *testing.T) {
	changeset := Changeset{
		Diffs: []*Diff{
			{
				DiffComments: CommentsTree{
					{
						Text: "\n\nevery dangling whitespace    \n\n" +
							"   should be removed\n\n\n",
					},
				},
			},
		},
	}

	buf := bytes.Buffer{}
	err := WriteChangeset(changeset, &buf)
	if err != nil {
		t.Fatal(err)
	}

	output, err := ioutil.ReadFile(
		"_test/no_dangling_whitespaces.diff")

	if string(output) != buf.String() {
		t.Log("all dangling whitespaces should be trimmed")
		t.Fatal("\n" + makeDiff(buf.String(), string(output)))
	}
}
