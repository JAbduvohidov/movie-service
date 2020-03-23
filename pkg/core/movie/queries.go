package movie

const getAllMoviesDML = `SELECT id,
       title,
       description,
       image,
       year,
       country,
       actors,
       genres,
       creators,
       studio,
       extLink
FROM movies
WHERE removed = FALSE;`

const searchMoviesDML = `SELECT id,
       title,
       description,
       image,
       year,
       country,
       actors,
       genres,
       creators,
       studio,
       extLink
FROM movies
WHERE removed = FALSE
  AND title ILIKE $1 OR description ILIKE $1 ORDER BY title;`

const getMovieDML = `SELECT
       id,
       title,
       description,
       image,
       year,
       country,
       actors,
       genres,
       creators,
       studio,
       extlink
FROM movies
WHERE id = $1;`

const addMovieDML = `INSERT INTO movies (title, description, image, year, country, actors, genres, creators, studio, extLink)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

const deleteMovieDML = `UPDATE movies SET removed = TRUE WHERE id = $1;`

const updateMovieDML = `UPDATE movies
SET title       = $1,
    description = $2,
    image       = $3,
    year        = $4,
    country     = $5,
    actors      = $6,
    genres      = $7,
    creators    = $8,
    studio      = $9,
    extlink     = $10
WHERE id = $11;`