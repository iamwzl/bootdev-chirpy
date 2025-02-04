-- +goose Up
ALTER TABLE users ADD COLUMN hashed_password text;
ALTER TABLE users ADD COLUMN password_reset_required BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE users
	SET hashed_password = 'RESETME', password_reset_required = TRUE
	WHERE hashed_password IS NULL;
ALTER TABLE users ALTER COLUMN hashed_password SET NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN password_reset_required;
ALTER TABLE users DROP COLUMN hashed_password;