CREATE TABLE IF NOT EXISTS cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    card_number_encrypted VARCHAR(255) NOT NULL,
    card_holder_encrypted VARCHAR(255) NOT NULL,
    expiry_encrypted VARCHAR(255) NOT NULL,
    cvv_encrypted VARCHAR(255) NOT NULL,
    meta TEXT,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cards_user_id ON cards(user_id);
CREATE INDEX idx_cards_title ON cards(title);