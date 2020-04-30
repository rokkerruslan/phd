-- mi: 010.add_image_hash.sql

-- TODO: Add NOT NULL constraint
ALTER TABLE images ADD COLUMN hash VARCHAR(256);

ALTER TABLE images DROP COLUMN path;
