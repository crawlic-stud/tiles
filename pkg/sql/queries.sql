-- name: GetGameByID :one
SELECT * FROM games WHERE id = ?;

-- name: CreateGame :one
INSERT INTO games (width, height, custom_tiles, background, tile_size) VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: CreateCharacter :one
INSERT INTO characters (name, type, scale, image) VALUES (?, ?, ?, ?) RETURNING *;

-- name: CreateGameCharacter :one
INSERT INTO game_characters (game_id, character_id, x, y) VALUES (?, ?, ?, ?) RETURNING *;

-- name: GetGameCharacterByID :one
SELECT c.*, gc.x, gc.y FROM characters c
JOIN game_characters gc ON c.id = gc.character_id
WHERE gc.game_id = ? AND c.id = ?;

-- name: GetGameCharacters :many
SELECT c.*, gc.x, gc.y FROM characters c
JOIN game_characters gc ON c.id = gc.character_id
WHERE gc.game_id = ?;

-- name: UpdateCharacterPosition :exec
UPDATE game_characters SET x = ?, y = ? WHERE game_id = ? AND character_id = ?;

-- name: UpdateGameTiles :exec
UPDATE games SET custom_tiles = ? WHERE id = ?;

-- name: HideGameTiles :exec
UPDATE games SET hide_tiles = ? WHERE id = ?;