package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type CardRepository struct {
	pool *pgxpool.Pool
}

func NewCardRepository(pool *pgxpool.Pool) *CardRepository {
	return &CardRepository{pool: pool}
}

func (r *CardRepository) Create(ctx context.Context, card *model.Card) error {
	query := `
		INSERT INTO cards (user_id, title, card_number_encrypted, card_holder_encrypted, expiry_encrypted, cvv_encrypted, meta, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		card.UserID,
		card.Title,
		card.CardNumberEncrypted,
		card.CardHolderEncrypted,
		card.ExpiryEncrypted,
		card.CVVEncrypted,
		card.Meta,
		card.Version,
		card.CreatedAt,
		card.UpdatedAt,
	).Scan(&card.ID)

	return err
}

func (r *CardRepository) GetByID(ctx context.Context, id string) (*model.Card, error) {
	query := `
		SELECT id, user_id, title, card_number_encrypted, card_holder_encrypted, expiry_encrypted, cvv_encrypted, meta, version, created_at, updated_at
		FROM cards
		WHERE id = $1
	`

	var card model.Card
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&card.ID,
		&card.UserID,
		&card.Title,
		&card.CardNumberEncrypted,
		&card.CardHolderEncrypted,
		&card.ExpiryEncrypted,
		&card.CVVEncrypted,
		&card.Meta,
		&card.Version,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, model.ErrCardNotFound
	}
	return &card, err
}

func (r *CardRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Card, error) {
	query := `
		SELECT id, user_id, title, card_number_encrypted, card_holder_encrypted, expiry_encrypted, cvv_encrypted, meta, version, created_at, updated_at
		FROM cards
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*model.Card
	for rows.Next() {
		var card model.Card
		if err := rows.Scan(
			&card.ID,
			&card.UserID,
			&card.Title,
			&card.CardNumberEncrypted,
			&card.CardHolderEncrypted,
			&card.ExpiryEncrypted,
			&card.CVVEncrypted,
			&card.Meta,
			&card.Version,
			&card.CreatedAt,
			&card.UpdatedAt,
		); err != nil {
			return nil, err
		}
		cards = append(cards, &card)
	}

	return cards, nil
}

func (r *CardRepository) Update(ctx context.Context, card *model.Card) error {
	query := `
		UPDATE cards
		SET title = $2, card_number_encrypted = $3, card_holder_encrypted = $4, expiry_encrypted = $5, cvv_encrypted = $6, meta = $7, version = $8, updated_at = $9
		WHERE id = $1 AND version = $10
	`

	result, err := r.pool.Exec(ctx, query,
		card.ID,
		card.Title,
		card.CardNumberEncrypted,
		card.CardHolderEncrypted,
		card.ExpiryEncrypted,
		card.CVVEncrypted,
		card.Meta,
		card.Version,
		card.UpdatedAt,
		card.Version-1,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrVersionConflict
	}

	return nil
}

func (r *CardRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM cards WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrCardNotFound
	}

	return nil
}
