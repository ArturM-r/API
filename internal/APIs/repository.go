package APIs

import (
	"API/internal/errs"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APRepository struct {
	Db *pgxpool.Pool
}

func NewAPRepository(db *pgxpool.Pool) *APRepository {
	return &APRepository{
		Db: db,
	}
}
func (r *APRepository) Create(ctx context.Context, userID uuid.UUID, title string) error {
	query := `INSERT INTO apis (user_id, title) VALUES ($1, $2)`

	tag, err := r.Db.Exec(ctx, query, userID, title)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil

}

func (r *APRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*API, error) {
	query := `SELECT * FROM apis WHERE id = $1 and user_id = $2`

	var u API

	err := r.Db.QueryRow(ctx, query, id, userID).Scan(
		&u.ID, &u.UserID, &u.UserID, &u.Completed, &u.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("no file %w", err)
	}
	return &u, nil
}

func (r *APRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM apis WHERE id = $1 and user_id = $2`
	tag, err := r.Db.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
func (r *APRepository) Update(ctx context.Context, id, userID uuid.UUID, title *string, completed *bool) error {
	if title == nil && completed == nil {
		return errs.ErrBadRequest
	}

	query := "UPDATE apis SET "
	args := make([]any, 0, 4)
	argPos := 1

	if title != nil {
		query += fmt.Sprintf("title = $%d, ", argPos)
		args = append(args, *title)
		argPos++
	}

	if completed != nil {
		query += fmt.Sprintf("completed = $%d, ", argPos)
		args = append(args, *completed)
		argPos++
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", argPos, argPos+1)

	args = append(args, id, userID)

	tag, err := r.Db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update api: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *APRepository) GetAll(ctx context.Context, userID uuid.UUID, filter QueryFilter) ([]API, error) {
	query := `SELECT id, user_id, title, completed, created_at FROM apis WHERE user_id = $1`

	args := []any{userID}
	argPos := 2

	if filter.completed != nil {
		query += fmt.Sprintf(" AND completed = $%d", argPos)
		args = append(args, *filter.completed)
		argPos++
	}

	if filter.title != nil {
		search := strings.TrimSpace(*filter.title)
		if search != "" {
			query += fmt.Sprintf(" AND title ILIKE '%%' || $%d || '%%'", argPos)
			args = append(args, search)
			argPos++
		}
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filter.limit, filter.offset)

	rows, err := r.Db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get all apis: %w", err)
	}
	defer rows.Close()

	var APIs []API

	for rows.Next() {
		var api API
		if err := rows.Scan(
			&api.ID,
			&api.UserID,
			&api.Title,
			&api.Completed,
			&api.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan api: %w", err)
		}

		APIs = append(APIs, api)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate apis: %w", err)
	}

	return APIs, nil
}
