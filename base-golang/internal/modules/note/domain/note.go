package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"base-golang/internal/shared/apperror"
)

type Note struct {
	ID        int64
	Title     string
	CreatedAt time.Time
}

func (note *Note) Validate() error {
	note.Title = strings.TrimSpace(note.Title)

	if note.Title == "" || utf8.RuneCountInString(note.Title) > 200 || strings.ContainsRune(note.Title, 0) {
		return apperror.New(apperror.Invalid, "title must contain 1 to 200 characters")
	}

	return nil
}
