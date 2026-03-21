package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type BinaryDataRepository struct {
	pool *pgxpool.Pool
}

func NewBinaryDataRepository(pool *pgxpool.Pool) *BinaryDataRepository {
	return &BinaryDataRepository{pool: pool}
}

func (r *BinaryDataRepository) Create(ctx context.Context, binaryData *model.BinaryData) error {
	query := `
		INSERT INTO binary_data (user_id, title, data_encrypted, meta, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		binaryData.UserID,
		binaryData.Title,
		binaryData.DataEncrypted,
		binaryData.Meta,
		binaryData.Version,
		binaryData.CreatedAt,
		binaryData.UpdatedAt,
	).Scan(&binaryData.ID)

	return err
}

func (r *BinaryDataRepository) GetByID(ctx context.Context, id string) (*model.BinaryData, error) {
	query := `
		SELECT id, user_id, title, data_encrypted, meta, version, created_at, updated_at
		FROM binary_data
		WHERE id = $1
	`

	var binaryData model.BinaryData
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&binaryData.ID,
		&binaryData.UserID,
		&binaryData.Title,
		&binaryData.DataEncrypted,
		&binaryData.Meta,
		&binaryData.Version,
		&binaryData.CreatedAt,
		&binaryData.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, model.ErrBinaryDataNotFound
	}
	return &binaryData, err
}

func (r *BinaryDataRepository) GetByUserID(ctx context.Context, userID string) ([]*model.BinaryData, error) {
	query := `
		SELECT id, user_id, title, data_encrypted, meta, version, created_at, updated_at
		FROM binary_data
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var binaryDataList []*model.BinaryData
	for rows.Next() {
		var binaryData model.BinaryData
		if err := rows.Scan(
			&binaryData.ID,
			&binaryData.UserID,
			&binaryData.Title,
			&binaryData.DataEncrypted,
			&binaryData.Meta,
			&binaryData.Version,
			&binaryData.CreatedAt,
			&binaryData.UpdatedAt,
		); err != nil {
			return nil, err
		}
		binaryDataList = append(binaryDataList, &binaryData)
	}

	return binaryDataList, nil
}

func (r *BinaryDataRepository) Update(ctx context.Context, binaryData *model.BinaryData) error {
	query := `
		UPDATE binary_data
		SET title = $2, data_encrypted = $3, meta = $4, version = $5, updated_at = $6
		WHERE id = $1 AND version = $7
	`

	result, err := r.pool.Exec(ctx, query,
		binaryData.ID,
		binaryData.Title,
		binaryData.DataEncrypted,
		binaryData.Meta,
		binaryData.Version,
		binaryData.UpdatedAt,
		binaryData.Version-1,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrVersionConflict
	}

	return nil
}

func (r *BinaryDataRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM binary_data WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrBinaryDataNotFound
	}

	return nil
}
