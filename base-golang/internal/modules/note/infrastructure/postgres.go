package infrastructure

import (
	"context"
	"database/sql"

	"base-golang/internal/modules/note/domain"
	"base-golang/internal/modules/note/service"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List(ctx context.Context, options service.ListOptions) ([]domain.Note, int64, error) {
	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM notes").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	const query = `
		SELECT id, title, created_at
		FROM notes
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	offset := int64(options.Page-1) * int64(options.PageSize)
	rows, err := r.db.QueryContext(ctx, query, options.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notes := []domain.Note{}
	for rows.Next() {
		var note domain.Note
		err = rows.Scan(&note.ID, &note.Title, &note.CreatedAt)
		if err != nil {
			return nil, 0, err
		}

		notes = append(notes, note)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

func (r *Repository) Create(ctx context.Context, note domain.Note) (domain.Note, error) {
	const query = `
		INSERT INTO notes (title)
		VALUES ($1)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		note.Title,
	).Scan(&note.ID, &note.CreatedAt)
	if err != nil {
		return domain.Note{}, err
	}

	return note, nil
}
