-- mi: 007.add_account_name.sql

-- TODO: Add NOT NULL

UPDATE accounts SET name = 'default' WHERE name IS NULL;

ALTER TABLE accounts ADD COLUMN name TEXT;

ALTER TABLE accounts ALTER COLUMN name SET NOT NULL;
