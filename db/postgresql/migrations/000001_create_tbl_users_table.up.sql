CREATE SEQUENCE IF NOT EXISTS seq_tbl_user_id START 1;

CREATE TABLE IF NOT EXISTS tbl_users (
    id              BIGINT PRIMARY KEY DEFAULT nextval('seq_tbl_user_id'),
    user_uuid       UUID        NOT NULL UNIQUE,

    first_name      TEXT        NOT NULL,
    last_name       TEXT        NOT NULL,
    user_name       TEXT        NOT NULL,
    password        TEXT        NOT NULL,
    email           TEXT        NOT NULL,
    role_id         INTEGER     NOT NULL,
    status          BOOLEAN     NOT NULL,
    login_session   TEXT,
    profile_photo   TEXT,
    user_alias      TEXT,
    phone_number    TEXT,
    user_avatar_id  INTEGER,
    commission      NUMERIC     DEFAULT 0.00,

    status_id       INTEGER     NOT NULL DEFAULT 1,
    "order"         INTEGER     DEFAULT 1,
    created_by      INTEGER     NOT NULL,
    created_at      TIMESTAMP   NOT NULL,
    updated_by      INTEGER,
    updated_at      TIMESTAMP,
    deleted_by      INTEGER,
    deleted_at      TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_tbl_users_user_uuid ON tbl_users(user_uuid);
CREATE INDEX IF NOT EXISTS idx_tbl_users_role_id ON tbl_users(role_id);

INSERT INTO tbl_users (
    user_uuid, first_name, last_name, user_name, password, email,
    role_id, status, login_session, profile_photo, user_alias, phone_number,
    commission, status_id, "order", created_by, created_at
) VALUES
(
    'c5b66b62-2cb0-4a2e-b704-1da97d8ed10d',
    'Supper', 'Admin', 'ADMIN', '123', 'admin@gmail.com',
    1, true, 'bdeb581454a4441784be1e355faeab63', 'user1.png', 'KM001', '010123123',
    0.00, 1, 1, 1, NOW()
),
(
    '83751b48-68f3-4805-a7bd-60ab8311936d',
    'IT', 'Developer', 'IT', '12e!!121#', 'it@gmail.com',
    1, true, 'bdeb581454a4441784be1e355faeab57', 'user2.png', 'KM002', '430123123',
    0.00, 1, 1, 1, NOW()
);
