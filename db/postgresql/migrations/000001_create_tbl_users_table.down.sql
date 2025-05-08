-- 001_create_tbl_users_table.down.sql

-- Drop the users table (which also drops related indexes)
DROP TABLE IF EXISTS tbl_users;

-- Drop the sequence explicitly
DROP SEQUENCE IF EXISTS seq_tbl_user_id;
