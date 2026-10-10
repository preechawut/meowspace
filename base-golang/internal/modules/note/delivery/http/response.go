package http

import (
	"time"

	"base-golang/internal/modules/note/domain"
)

// dto
type NoteResponse struct {
	ID        int64     `json:"id,omitempty"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func toNoteResponse(value domain.Note) NoteResponse {
	return NoteResponse{
		ID:        value.ID,
		Title:     value.Title,
		CreatedAt: value.CreatedAt,
	}
}

func toNoteResponses(values []domain.Note) []NoteResponse {
	responses := make([]NoteResponse, len(values))
	for i, value := range values {
		responses[i] = toNoteResponse(value)
	}
	return responses
}

type NotePageResponse struct {
	Items      []NoteResponse `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int64          `json:"total"`
	TotalPages int64          `json:"total_pages"`
}
