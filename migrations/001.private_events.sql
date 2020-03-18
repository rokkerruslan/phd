-- mi: 001.private_events.sql

ALTER TABLE events ADD COLUMN is_private BOOL NOT NULL DEFAULT FALSE;
