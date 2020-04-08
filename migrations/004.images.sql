-- mi: 004.images.sql

CREATE TABLE images (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    path  TEXT NOT NULL,
    title TEXT NOT NULL,

    author_id INTEGER NOT NULL,
    event_id  INTEGER NOT NULL,

    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,

    FOREIGN KEY (author_id) REFERENCES accounts(id),
    FOREIGN KEY (event_id)  REFERENCES events(id)
)
