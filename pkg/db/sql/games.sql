-- name: GetGameByID :one
SELECT * FROM games WHERE id = $1;

-- name: CreateGame :one
INSERT INTO games (width, height, custom_tiles, background, tile_size) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateGameTiles :exec
UPDATE games SET custom_tiles = $1 WHERE id = $2;

-- name: HideGameTiles :exec
UPDATE games SET hide_tiles = $1 WHERE id = $2;
