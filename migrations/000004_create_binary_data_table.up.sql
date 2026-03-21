CREATE TABLE IF NOT EXISTS binary_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    data_encrypted BYTEA NOT NULL,
    meta TEXT,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_binary_data_user_id ON binary_data(user_id);
CREATE INDEX idx_binary_data_title ON binary_data(title);