-- mi: 006.offer_approving.sql

ALTER TABLE offers ADD COLUMN is_approved BOOL DEFAULT FALSE;
