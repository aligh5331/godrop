CREATE TABLE refresh_tokens
(
    id           VARCHAR PRIMARY KEY,
    session_id   VARCHAR     NOT NULL REFERENCES sessions (id),
    hashed_token VARCHAR     NOT NULL UNIQUE,
    created_at   TIMESTAMPTZ NOT NULL,
    expires_at   TIMESTAMPTZ NOT NULL,
    revoked_at   TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00+00'
);

CREATE UNIQUE INDEX one_active_rt_per_session
    ON refresh_tokens (session_id) WHERE revoked_at = '0001-01-01 00:00:00+00';

CREATE INDEX idx_refresh_tokens_hashed_token ON refresh_tokens (hashed_token);
CREATE INDEX idx_refresh_tokens_session_id ON refresh_tokens (session_id);
