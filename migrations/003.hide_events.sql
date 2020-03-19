-- mi: 003.hide_events.sql

ALTER TABLE events ADD COLUMN is_hidden BOOL DEFAULT FALSE;
