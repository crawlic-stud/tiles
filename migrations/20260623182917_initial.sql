-- +goose Up
CREATE TABLE games(
    id BIGSERIAL PRIMARY KEY,
    custom_tiles JSONB NOT NULL,
    background TEXT NOT NULL,
    tile_size INT NOT NULL,
    hide_tiles BOOLEAN NOT NULL DEFAULT FALSE,
    width INT NOT NULL,
    height INT NOT NULL
);

CREATE TABLE characters(
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    scale REAL NOT NULL,
    image TEXT NOT NULL
);

CREATE TABLE game_characters(
    id BIGSERIAL PRIMARY KEY,
    game_id BIGINT NOT NULL,
    character_id BIGINT NOT NULL,
    x INT NOT NULL,
    y INT NOT NULL,
    FOREIGN KEY(game_id) REFERENCES games(id),
    FOREIGN KEY(character_id) REFERENCES characters(id)
);

-- +goose Down
DROP TABLE game_characters;
DROP TABLE characters;
DROP TABLE games;
