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
    year        = $3,
    country     = $4,
    actors      = $5,
    genres      = $6,
    creators    = $7,
    studio      = $8,
    extlink     = $9
WHERE id = $10;`

const updateMovieImageDML = `UPDATE movies
SET image = $1
WHERE id = $2;`
