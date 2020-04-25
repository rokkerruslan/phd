-- mi: 008.account_deleting.sql

-- TODO: Add NOT NULL constraint
ALTER TABLE accounts ADD COLUMN is_deleted BOOL NOT NULL DEFAULT FALSE;
