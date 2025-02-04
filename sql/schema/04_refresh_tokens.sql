-- +goose Up
CREATE TABLE refresh_tokens (
	token text PRIMARY KEY,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
	user_id UUID NOT NULL,
	expires_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP + (60*INTERVAL '1 days') NOT NULL,
	revoked_at TIMESTAMPTZ,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON refresh_tokens
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS set_updated_at ON refresh_tokens;
DROP TABLE refresh_tokens;

