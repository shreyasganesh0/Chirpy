-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    email TEXT NOT NULL
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_col()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON users 
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_col();

-- +goose Down
DROP TRIGGER IF EXISTS set_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_col;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";

