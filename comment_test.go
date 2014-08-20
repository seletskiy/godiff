package godiff

import "testing"

func TestShortComment(t *testing.T) {
	tests := []struct {
		comment Comment
		short   string
	}{
		{
			Comment{Text: "Привет!\n   Как дела? Нормально?"},
			"Привет! Как дела? Но...",
		},
		{
			Comment{Text: "fixed"},
			"fixed",
		},
	}

	for _, test := range tests {
		actual := test.comment.Short(commentShortLength)
		if actual != test.short {
			t.Fatalf("unexpected %#v", actual)
		}
	}
}
