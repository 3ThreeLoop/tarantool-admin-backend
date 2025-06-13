-- +goose Up
-- USERS TABLE
CREATE TABLE tbl_users (
    id SERIAL PRIMARY KEY,
    user_uuid UUID NOT NULL UNIQUE,
    first_name VARCHAR  NULL,
    last_name VARCHAR NULL,
    user_name VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    login_session VARCHAR,
    profile_photo VARCHAR,
    status_id INTEGER NOT NULL DEFAULT 1,
    "order" INTEGER DEFAULT 1,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_by INTEGER,
    updated_at TIMESTAMP,
    deleted_by INTEGER,
    deleted_at TIMESTAMP
);

-- +goose StatementBegin
INSERT INTO tbl_users (
    user_uuid, first_name, last_name, user_name, password, email,
    login_session, profile_photo, status_id, "order",
    created_by, created_at, updated_by, updated_at, deleted_by, deleted_at
) VALUES 
(
    'c5b66b62-2cb0-4a2e-b704-1da97d8ed10d', 'Supper', 'Admin', 'ADMIN', '123',
    'admin@gmail.com', 'bdeb581454a4441784be1e355faeab63',
    'user1.png', 1, 1, 1, NOW(), NULL, NULL, NULL, NULL
),
(
    '83751b48-68f3-4805-a7bd-60ab8311936d', 'IT', 'Developer', 'IT', '12e!!121#',
    'it@gmail.com', 'bdeb581454a4441784be1e355faeab57',
    'user2.png', 1, 1, 1, NOW(), NULL, NULL, NULL, NULL
);

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS tbl_users;