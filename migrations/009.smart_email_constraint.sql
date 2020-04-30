-- mi: 009.smart_email_constraint.sql

ALTER TABLE accounts DROP CONSTRAINT accounts_email_key;

CREATE UNIQUE INDEX unique_emails ON accounts (email) WHERE is_deleted = FALSE;
