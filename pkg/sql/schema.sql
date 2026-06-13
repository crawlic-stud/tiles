CREATE TABLE games(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    custom_tiles TEXT NOT NULL,
    background TEXT NOT NULL,
    tile_size INTEGER NOT NULL,
    hide_tiles BOOLEAN NOT NULL DEFAULT FALSE,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL
);

CREATE TABLE characters(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    scale REAL NOT NULL,
    image TEXT NOT NULL
);

CREATE TABLE game_characters(
    game_id INTEGER NOT NULL,
    character_id INTEGER NOT NULL,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    FOREIGN KEY(game_id) REFERENCES games(id),
    FOREIGN KEY(character_id) REFERENCES characters(id),
    PRIMARY KEY(game_id, character_id)
);