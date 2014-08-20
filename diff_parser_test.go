package godiff

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
)

type diffTest struct {
	in  []byte
	out []byte
	err []byte
}

func TestParseDiff(t *testing.T) {
	for name, testCase := range getParseTests("_test") {
		expected := string(testCase.out)

		review, err := ParseDiff(bytes.NewBuffer(testCase.in))
		if err != nil {
			t.Fatal(err)
		}

		actual := review.String()
		if actual != expected {
			t.Logf("while testing on `%s`\n", name)
			t.Logf("expected:\n%v", expected)
			t.Logf("actual:\n%v", actual)
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
		//{"_test/new_comment_to_del.diff", 3, SegmentTypeRemoved},
		//{"_test/new_comment_to_ctx_with_add.diff", 3, SegmentTypeContext},
		//{"_test/new_comment_to_ctx_with_add_del.diff", 5, SegmentTypeContext},
	}

	for _, test := range tests {
		reviewText, err := ioutil.ReadFile(test.file)
		if err != nil {
			t.Fatal(err)
		}

		review, err := ParseDiff(bytes.NewBuffer(reviewText))
		if err != nil {
			t.Fatal(err)
		}

		if len(review.Diffs) != 1 {
			t.Fatal("unexpected number of review")
			t.FailNow()
		}

		comment := review.Diffs[0].LineComments[0]

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
	reviewText, err := ioutil.ReadFile("_test/with_one_nested_comment.diff")
	if err != nil {
		t.Fatal(err)
	}

	review, err := ParseDiff(bytes.NewBuffer(reviewText))
	if err != nil {
		t.Fatal(err)
	}

	if len(review.Diffs[0].LineComments) != 1 {
		t.Fatal("expected only one line comment")
	}
}

func getParseTests(dir string) map[string]*diffTest {
	diffTests := make(map[string]*diffTest)
	reDiffTest := regexp.MustCompile(`/([^/]+)\.(in|out|err)\.diff$`)

	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !reDiffTest.MatchString(path) {
			return nil
		}

		matches := reDiffTest.FindStringSubmatch(path)

		caseName := matches[1]

		if _, exist := diffTests[caseName]; !exist {
			diffTests[caseName] = &diffTest{}
		}

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

		return nil
	})

	for _, val := range diffTests {
		if val.out == nil {
			val.out = val.in
		}
	}

	return diffTests
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
