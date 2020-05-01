-- mi: 010.add_image_hash.sql

-- TODO: Add NOT NULL constraint
ALTER TABLE images ADD COLUMN hash VARCHAR(64);

-- TODO: Remove from first migration
ALTER TABLE images DROP COLUMN path;
