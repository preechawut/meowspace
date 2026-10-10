package service

import (
	"context"
	"errors"
	"testing"

	"base-golang/internal/modules/note/domain"
)

type fakeRepository struct {
	calls    int
	received domain.Note
	err      error
	notes    []domain.Note
	options  ListOptions
	total    int64
}

func (r *fakeRepository) Create(ctx context.Context, note domain.Note) (domain.Note, error) {
	r.calls++
	r.received = note
	note.ID = 42
	return note, r.err
}

func (r *fakeRepository) List(ctx context.Context, options ListOptions) ([]domain.Note, int64, error) {
	r.options = options
	return r.notes, r.total, r.err
}

func TestCreateRejectsInvalidInputBeforeRepository(t *testing.T) {
	repository := &fakeRepository{}
	service := New(repository)
	for _, title := range []string{" ", "a\x00b"} {
		_, err := service.Create(context.Background(), domain.Note{Title: title})
		if err == nil || repository.calls != 0 {
			t.Fatal("invalid input reached repository")
		}
	}
}

func TestCreateNormalizesInputAndReturnsSavedNote(t *testing.T) {
	repository := &fakeRepository{}
	note, err := New(repository).Create(context.Background(), domain.Note{Title: "  hello  "})
	if err != nil || note.ID != 42 || repository.received.Title != "hello" {
		t.Fatalf("unexpected result: %+v, %v", note, err)
	}
}

func TestRepositoryErrorsArePreserved(t *testing.T) {
	failure := errors.New("database unavailable")
	service := New(&fakeRepository{err: failure})
	_, err := service.Create(context.Background(), domain.Note{Title: "hello"})
	if !errors.Is(err, failure) {
		t.Fatal("create lost repository error")
	}
	_, err = service.List(context.Background(), ListOptions{Page: 1, PageSize: 20})
	if !errors.Is(err, failure) {
		t.Fatal("list lost repository error")
	}
}

func TestListPagination(t *testing.T) {
	repo := &fakeRepository{notes: []domain.Note{{ID: 9}, {ID: 8}}, total: 5}
	page, err := New(repo).List(context.Background(), ListOptions{Page: 2, PageSize: 2})
	if err != nil || len(page.Notes) != 2 || page.Total != 5 || page.TotalPages != 3 || repo.options.Page != 2 {
		t.Fatalf("page: %+v %v", page, err)
	}
	repo.notes = nil
	repo.total = 0
	page, err = New(repo).List(context.Background(), ListOptions{Page: 1, PageSize: 20})
	if err != nil || page.TotalPages != 0 {
		t.Fatalf("empty page: %+v %v", page, err)
	}
	for _, options := range []ListOptions{{Page: 0, PageSize: 20}, {Page: 1, PageSize: 101}, {Page: 1, PageSize: 0}} {
		if _, err := New(repo).List(context.Background(), options); err == nil {
			t.Fatalf("accepted %+v", options)
		}
	}
}
