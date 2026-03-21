package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type CredentialRepository struct {
	pool *pgxpool.Pool
}

func NewCredentialRepository(pool *pgxpool.Pool) *CredentialRepository {
	return &CredentialRepository{pool: pool}
}

func (r *CredentialRepository) Create(ctx context.Context, credential *model.Credential) error {
	query := `
		INSERT INTO credentials (user_id, title, login, password_encrypted, meta, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		credential.UserID,
		credential.Title,
		credential.Login,
		credential.PasswordEncrypted,
		credential.Meta,
		credential.Version,
		credential.CreatedAt,
		credential.UpdatedAt,
	).Scan(&credential.ID)

	return err
}

func (r *CredentialRepository) GetByID(ctx context.Context, id string) (*model.Credential, error) {
	query := `
		SELECT id, user_id, title, login, password_encrypted, meta, version, created_at, updated_at
		FROM credentials
		WHERE id = $1
	`

	var credential model.Credential
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&credential.ID,
		&credential.UserID,
		&credential.Title,
		&credential.Login,
		&credential.PasswordEncrypted,
		&credential.Meta,
		&credential.Version,
		&credential.CreatedAt,
		&credential.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, model.ErrCredentialNotFound
	}
	return &credential, err
}

func (r *CredentialRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Credential, error) {
	query := `
		SELECT id, user_id, title, login, password_encrypted, meta, version, created_at, updated_at
		FROM credentials
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credentials []*model.Credential
	for rows.Next() {
		var credential model.Credential
		if err := rows.Scan(
			&credential.ID,
			&credential.UserID,
			&credential.Title,
			&credential.Login,
			&credential.PasswordEncrypted,
			&credential.Meta,
			&credential.Version,
			&credential.CreatedAt,
			&credential.UpdatedAt,
		); err != nil {
			return nil, err
		}
		credentials = append(credentials, &credential)
	}

	return credentials, nil
}

func (r *CredentialRepository) Update(ctx context.Context, credential *model.Credential) error {
	query := `
		UPDATE credentials
		SET title = $2, login = $3, password_encrypted = $4, meta = $5, version = $6, updated_at = $7
		WHERE id = $1 AND version = $8
	`

	result, err := r.pool.Exec(ctx, query,
		credential.ID,
		credential.Title,
		credential.Login,
		credential.PasswordEncrypted,
		credential.Meta,
		credential.Version,
		credential.UpdatedAt,
		credential.Version-1,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrVersionConflict
	}

	return nil
}

func (r *CredentialRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM credentials WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrCredentialNotFound
	}

	return nil
}
