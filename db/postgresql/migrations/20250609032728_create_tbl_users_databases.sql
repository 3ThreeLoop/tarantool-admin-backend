-- +goose Up
-- USERS DATABASE CONNECTIONS TABLE
CREATE TABLE tbl_users_databases (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES tbl_users(id) ON DELETE CASCADE,
    db_uuid UUID NOT NULL UNIQUE,
    db_name VARCHAR NOT NULL,
    host VARCHAR NOT NULL,
    port INTEGER NOT NULL,
    username VARCHAR,
    password VARCHAR,
    is_active BOOLEAN DEFAULT TRUE,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_by INTEGER,
    updated_at TIMESTAMP,
    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS tbl_users_databases;
