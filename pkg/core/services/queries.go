package services

const moviesDDL = `CREATE TABLE IF NOT EXISTS movies
(
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    image       TEXT NOT NULL,
    year        TEXT    DEFAULT '',
    country     TEXT    DEFAULT '',
    actors      TEXT[]  DEFAULT '{}',
    genres      TEXT[]  DEFAULT '{}',
    creators    TEXT[]  DEFAULT '{}',
    studio      TEXT    DEFAULT '',
    extLink     TEXT    DEFAULT '',
    removed     BOOLEAN DEFAULT FALSE
);`