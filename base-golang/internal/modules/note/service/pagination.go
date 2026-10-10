package service

import (
	"base-golang/internal/modules/note/domain"
	"base-golang/internal/shared/apperror"
)

type ListOptions struct {
	Page     int
	PageSize int
}

func (o ListOptions) Validate() error {
	if o.Page < 1 {
		return apperror.New(apperror.Invalid, "page must be a positive integer")
	}
	if o.PageSize < 1 || o.PageSize > 100 {
		return apperror.New(apperror.Invalid, "page_size must be between 1 and 100")
	}
	if o.Page > 100000 {
		return apperror.New(apperror.Invalid, "page is too large")
	}

	return nil
}

type NotePage struct {
	Notes      []domain.Note
	Total      int64
	TotalPages int64
}
