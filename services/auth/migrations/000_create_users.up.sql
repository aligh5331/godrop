CREATE TABLE users
(
    id         VARCHAR PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL,
    updated_at TIMESTAMPTZ  NOT NULL,
    deleted_at TIMESTAMPTZ
);

-- Unique index for email as specified in the GORM tags
CREATE UNIQUE INDEX idx_users_email ON users (email);

-- Index for soft deletes to keep queries efficient
CREATE INDEX idx_users_deleted_at ON users (deleted_at);