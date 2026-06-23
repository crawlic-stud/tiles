-- name: CreateCharacter :one
INSERT INTO characters (name, type, scale, image) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateGameCharacter :one
INSERT INTO game_characters (game_id, character_id, x, y) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetGameCharacterByID :one
SELECT c.*, gc.x, gc.y FROM characters c
JOIN game_characters gc ON c.id = gc.character_id
WHERE gc.game_id = $1 AND c.id = $2;

-- name: GetGameCharacters :many
SELECT c.*, gc.x, gc.y FROM characters c
JOIN game_characters gc ON c.id = gc.character_id
WHERE gc.game_id = $1;

-- name: UpdateCharacterPosition :exec
UPDATE game_characters SET x = $1, y = $2 WHERE game_id = $3 AND character_id = $4;
