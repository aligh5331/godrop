CREATE TABLE users
(
    id         VARCHAR PRIMARY KEY,
    created_at TIMESTAMPTZ  NOT NULL,
    deleted_at TIMESTAMPTZ
);

-- Index for soft deletes to keep queries efficient
CREATE INDEX idx_users_deleted_at ON users (deleted_at);