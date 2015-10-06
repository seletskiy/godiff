package godiff

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

type diffTest struct {
	name string
	in   []byte
	out  []byte
	err  []byte
}

func TestReadChangeset(t *testing.T) {
	for _, testCase := range getParseTests("_test") {
		expected := string(testCase.out)

		changeset, err := ReadChangeset(bytes.NewBuffer(testCase.in))
		if err != nil {
			if err.Error() != strings.TrimSpace(string(testCase.err)) {
				t.Logf("while testing on `%s`\n", testCase.name)
				t.Fatal(err)
			} else {
				continue
			}
		}

		buf := &bytes.Buffer{}
		err = WriteChangeset(changeset, buf)
		if err != nil {
			panic(err)
		}

		actual := buf.String()
		if actual != expected {
			t.Logf("while testing on `%s`\n", testCase.name)
			t.Logf("expected:\n%#v", expected)
			t.Logf("actual:\n%#v", actual)
			t.Logf("diff:\n%v", makeDiff(actual, expected))
			t.FailNow()
		}
	}
}

func TestNewCommentForLine(t *testing.T) {
	tests := []struct {
		file     string
		expected CommentAnchor
	}{
		{"_test/new_comment_to_add.diff",
			CommentAnchor{
				Line: 3, LineType: SegmentTypeAdded, Path: "/tmp/b"}},
		{"_test/new_comment_to_del.diff",
			CommentAnchor{
				Line:     3,
				LineType: SegmentTypeRemoved,
				Path:     "/tmp/a",
			}},
		{"_test/new_comment_to_ctx_with_add.diff",
			CommentAnchor{
				Line:     3,
				LineType: SegmentTypeContext,
				Path:     "/tmp/a",
			}},
		{"_test/new_comment_to_ctx_with_add_del.diff",
			CommentAnchor{
				Line:     5,
				LineType: SegmentTypeContext,
				Path:     "/tmp/a",
			}},
	}

	for _, test := range tests {
		changesetText, err := ioutil.ReadFile(test.file)
		if err != nil {
			t.Fatal(err)
		}

		changeset, err := ReadChangeset(bytes.NewBuffer(changesetText))
		if err != nil {
			t.Fatal(err)
		}

		if len(changeset.Diffs) != 1 {
			t.Fatal("unexpected number of changeset")
			t.FailNow()
		}

		comment := changeset.Diffs[0].LineComments[0]

		if comment.Anchor.Line != test.expected.Line {
			t.Fatal("comment binded to incorrect line")
			t.FailNow()
		}

		if comment.Anchor.LineType != test.expected.LineType {
			t.Fatal("unexpected line type")
			t.FailNow()
		}

		if comment.Anchor.Path != test.expected.Path {
			t.Fatal("unexpected path", comment.Anchor.Path)
			t.FailNow()
		}
	}
}

func TestDoNotAddNestedCommentsToLineComments(t *testing.T) {
	changesetText, err := ioutil.ReadFile("_test/with_one_nested_comment.diff")
	if err != nil {
		t.Fatal(err)
	}

	changeset, err := ReadChangeset(bytes.NewBuffer(changesetText))
	if err != nil {
		t.Fatal(err)
	}

	if len(changeset.Diffs[0].LineComments) != 1 {
		t.Fatal("expected only one line comment")
	}
}

func getParseTests(dir string) []*diffTest {
	diffTests := make(map[string]*diffTest)
	diffTestsList := make([]*diffTest, 0)
	reChangesetTest := regexp.MustCompile(`/([^/]+)\.(in|out|err)\.diff$`)

	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !reChangesetTest.MatchString(path) {
			return nil
		}

		matches := reChangesetTest.FindStringSubmatch(path)

		caseName := matches[1]

		if _, exist := diffTests[caseName]; !exist {
			diffTests[caseName] = &diffTest{}
		}

		diffTests[caseName].name = caseName

		var target *[]byte

		switch matches[2] {
		case "in":
			target = &diffTests[caseName].in
		case "out":
			target = &diffTests[caseName].out
		case "err":
			target = &diffTests[caseName].err
		}

		*target, _ = ioutil.ReadFile(path)

		diffTestsList = append(diffTestsList, diffTests[caseName])

		return nil
	})

	for _, val := range diffTests {
		if val.out == nil {
			val.out = val.in
		}
	}

	return diffTestsList
}

func makeDiff(actual, expected string) string {
	a, _ := ioutil.TempFile(os.TempDir(), "actual")
	defer func() {
		os.Remove(a.Name())
	}()
	b, _ := ioutil.TempFile(os.TempDir(), "expected")
	defer func() {
		os.Remove(b.Name())
	}()

	a.WriteString(actual)
	b.WriteString(expected)
	cmd := exec.Command("diff", "-u", b.Name(), a.Name())
	buf := bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	cmd.Run()

	return buf.String()
}
