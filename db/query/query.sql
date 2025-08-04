-- name: GetAnime :one
SELECT * FROM anime WHERE id = $1;

-- name: GetAnimes :many
SELECT * FROM anime ORDER BY id DESC;

-- name: CreateAnime :one
INSERT INTO anime (title, description, image, type) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: DeleteAnime :exec
DELETE FROM anime WHERE id = $1;

-- name: UpdateAnime :one
UPDATE anime SET title = $2, description = $3, image = $4, type = $5 WHERE id = $1 RETURNING *;

-- name: CreateGenre :one
INSERT INTO genre (name) VALUES ($1) RETURNING *;

-- name: GetGenre :one
SELECT * FROM genre WHERE id = $1;

-- name: GetGenres :many
SELECT * FROM genre;

-- name: DeleteGenre :exec
DELETE FROM genre WHERE id = $1;

-- name: UpdateGenre :one
UPDATE genre SET name = $2 WHERE id = $1 RETURNING *;

-- name: GetGenresByAnimeId :many
SELECT * FROM genre WHERE id IN (SELECT genre_id FROM anime_genre WHERE anime_id = $1);

-- name: GetAnimesByGenreId :many
SELECT * FROM anime WHERE id IN (SELECT anime_id FROM anime_genre WHERE genre_id = $1);

-- name: CreateAnimeGenre :one
INSERT INTO anime_genre (anime_id, genre_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteAnimeGenre :exec
DELETE FROM anime_genre WHERE anime_id = $1 AND genre_id = $2;

-- name: CreateEpisode :one
INSERT INTO episode (episode_number, episode_url, anime_id) VALUES ($1, $2, $3) RETURNING *;

-- name: GetEpisodesByAnimeId :many
SELECT * FROM episode WHERE anime_id = $1;

-- name: GetEpisode :one
SELECT * FROM episode WHERE id = $1;

-- name: DeleteEpisode :exec
DELETE FROM episode WHERE id = $1;

-- name: UpdateEpisode :one
UPDATE episode SET episode_number = $2, episode_url = $3 WHERE id = $1 RETURNING *;



