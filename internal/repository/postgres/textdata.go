package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type TextDataRepository struct {
	pool *pgxpool.Pool
}

func NewTextDataRepository(pool *pgxpool.Pool) *TextDataRepository {
	return &TextDataRepository{pool: pool}
}

func (r *TextDataRepository) Create(ctx context.Context, textData *model.TextData) error {
	query := `
		INSERT INTO text_data (user_id, title, data_encrypted, meta, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		textData.UserID,
		textData.Title,
		textData.DataEncrypted,
		textData.Meta,
		textData.Version,
		textData.CreatedAt,
		textData.UpdatedAt,
	).Scan(&textData.ID)

	return err
}

func (r *TextDataRepository) GetByID(ctx context.Context, id string) (*model.TextData, error) {
	query := `
		SELECT id, user_id, title, data_encrypted, meta, version, created_at, updated_at
		FROM text_data
		WHERE id = $1
	`

	var textData model.TextData
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&textData.ID,
		&textData.UserID,
		&textData.Title,
		&textData.DataEncrypted,
		&textData.Meta,
		&textData.Version,
		&textData.CreatedAt,
		&textData.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, model.ErrTextDataNotFound
	}
	return &textData, err
}

func (r *TextDataRepository) GetByUserID(ctx context.Context, userID string) ([]*model.TextData, error) {
	query := `
		SELECT id, user_id, title, data_encrypted, meta, version, created_at, updated_at
		FROM text_data
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var textDataList []*model.TextData
	for rows.Next() {
		var textData model.TextData
		if err := rows.Scan(
			&textData.ID,
			&textData.UserID,
			&textData.Title,
			&textData.DataEncrypted,
			&textData.Meta,
			&textData.Version,
			&textData.CreatedAt,
			&textData.UpdatedAt,
		); err != nil {
			return nil, err
		}
		textDataList = append(textDataList, &textData)
	}

	return textDataList, nil
}

func (r *TextDataRepository) Update(ctx context.Context, textData *model.TextData) error {
	query := `
		UPDATE text_data
		SET title = $2, data_encrypted = $3, meta = $4, version = $5, updated_at = $6
		WHERE id = $1 AND version = $7
	`

	result, err := r.pool.Exec(ctx, query,
		textData.ID,
		textData.Title,
		textData.DataEncrypted,
		textData.Meta,
		textData.Version,
		textData.UpdatedAt,
		textData.Version-1,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrVersionConflict
	}

	return nil
}

func (r *TextDataRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM text_data WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrTextDataNotFound
	}

	return nil
}
