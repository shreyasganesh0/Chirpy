-- +goose Up
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT user_fk
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ
);

CREATE TRIGGER set_updated_at_refresh_tokens
BEFORE UPDATE ON refresh_tokens 
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_col();

-- +goose Down
DROP TRIGGER IF EXISTS set_updated_at_refresh_tokens ON refresh_tokens;
DROP TABLE refresh_tokens;
