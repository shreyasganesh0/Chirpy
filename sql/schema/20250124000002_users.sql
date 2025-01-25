-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE chirps (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    body VARCHAR(140) NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT user_fk
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE
);


CREATE TRIGGER set_updated_at_chirps
BEFORE UPDATE ON chirps 
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_col();

-- +goose Down
DROP TRIGGER IF EXISTS set_updated_at_chirps ON chirps;
DROP TABLE IF EXISTS chirps;

