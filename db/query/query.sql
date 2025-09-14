-- name: GetAnime :one
SELECT * FROM anime WHERE id = $1;

-- name: GetAnimes :many
SELECT * FROM anime ORDER BY id DESC;

-- name: GetAnimeBySlug :one
SELECT * FROM anime WHERE slug = $1 LIMIT 1;

-- name: CreateAnime :one
INSERT INTO anime (title, description, image, type, slug) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: DeleteAnime :exec
DELETE FROM anime WHERE id = $1;

-- name: UpdateAnime :one
UPDATE anime SET title = $2, description = $3, image = $4, type = $5, slug = $6 WHERE id = $1 RETURNING *;

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
SELECT e.* FROM episode e JOIN anime a ON e.anime_id = a.id WHERE a.slug = $1 ORDER BY e.episode_number ASC;

-- name: CheckEpisodeExists :one
SELECT * FROM episode WHERE episode_number = $1 AND anime_id = $2;

-- name: GetEpisodeBySlugAndNumber :one
SELECT 
    e.id AS episode_id,
    e.episode_number,
    e.episode_url,
    e.created_at AS episode_created_at,
    e.updated_at AS episode_updated_at,
    a.id AS anime_id,
    a.title AS anime_title,
    a.slug AS anime_slug,
    a.description AS anime_description,
    a.image AS anime_image,
    a.type AS anime_type
FROM episode e
JOIN anime a ON e.anime_id = a.id
WHERE a.slug = $1 AND e.episode_number = $2
LIMIT 1;

-- name: DeleteEpisode :exec
DELETE FROM episode WHERE id = $1;

-- name: UpdateEpisode :one
UPDATE episode SET episode_url = $2 WHERE id = $1 RETURNING *;

-- name: CreateUser :one
INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users SET username = $2, email = $3, password = $4, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetCommentsByEpisodeId :many
SELECT * FROM comments WHERE episode_id = $1;

-- name: CreateComment :one
INSERT INTO comments (content, user_id, episode_id) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1;

-- name: UpdateComment :one
UPDATE comments SET content = $2 WHERE id = $1 RETURNING *;

-- name: GetFavoritesByUserId :many
SELECT * FROM favorites WHERE user_id = $1;

-- name: CreateFavorite :one
INSERT INTO favorites (user_id, anime_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteFavorite :exec
DELETE FROM favorites WHERE id = $1;