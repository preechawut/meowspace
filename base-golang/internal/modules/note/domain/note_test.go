package domain

import (
	"strings"
	"testing"
)

func TestTitleRules(t *testing.T) {
	for _, title := range []string{"", "   ", "a\x00b", strings.Repeat("ก", 201)} {
		if _, err := validatedNote(title); err == nil {
			t.Fatalf("accepted invalid title %q", title)
		}
	}
	note, err := validatedNote("  hello  ")
	if err != nil || note.Title != "hello" {
		t.Fatalf("unexpected note: %+v, %v", note, err)
	}
	if _, err := validatedNote(strings.Repeat("ก", 200)); err != nil {
		t.Fatal(err)
	}
}

func validatedNote(title string) (Note, error) {
	note := Note{Title: title}
	err := note.Validate()
	return note, err
}
