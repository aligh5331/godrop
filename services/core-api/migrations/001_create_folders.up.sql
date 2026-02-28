CREATE TABLE folders
(
    id         VARCHAR     NOT NULL PRIMARY KEY,
    user_id    VARCHAR     NOT NULL REFERENCES users (id),
    parent_id  VARCHAR REFERENCES folders (id),
    name       VARCHAR     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
)

-- Index for soft deletes to keep queries efficient
CREATE INDEX idx_folders_deleted_at ON folders (deleted_at);

CREATE UNIQUE INDEX idx_unique_folder_name_active
    ON folders (user_id, parent_id, name) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_root_name
    ON folders (user_id, name) WHERE parent_id IS NULL AND deleted_at IS NULL;
CREATE INDEX idx_folders_parent_id
    ON folders (parent_id) WHERE deleted_at IS NULL;