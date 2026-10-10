package service

import (
	"context"

	"base-golang/internal/modules/note/domain"
)

type Repository interface {
	List(ctx context.Context, options ListOptions) ([]domain.Note, int64, error)
	Create(ctx context.Context, note domain.Note) (domain.Note, error)
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) List(ctx context.Context, options ListOptions) (NotePage, error) {
	err := options.Validate()
	if err != nil {
		return NotePage{}, err
	}

	notes, total, err := s.repository.List(ctx, options)
	if err != nil {
		return NotePage{}, err
	}

	totalPages := total / int64(options.PageSize)
	if total%int64(options.PageSize) != 0 {
		totalPages++
	}

	return NotePage{
		Notes:      notes,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) Create(ctx context.Context, note domain.Note) (domain.Note, error) {
	err := note.Validate()
	if err != nil {
		return domain.Note{}, err
	}

	return s.repository.Create(ctx, note)
}
