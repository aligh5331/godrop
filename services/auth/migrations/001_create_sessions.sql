CREATE TABLE sessions
(
    id           VARCHAR PRIMARY KEY,
    user_id      VARCHAR     NOT NULL,
    user_agent   TEXT,
    ip           VARCHAR     NOT NULL,
    hashed_token VARCHAR     NOT NULL UNIQUE,
    is_revoked   BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL,
    expires_at   TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_hashed_token ON sessions (hashed_token);
