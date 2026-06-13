package models

type CustomTile struct {
	Type TileType `json:"type"`
	X    int      `json:"x"`
	Y    int      `json:"y"`
}

type TileKey struct {
	X int
	Y int
}

type CustomTiles map[TileKey]TileType

func (t CustomTiles) Set(tile CustomTile) {
	key := TileKey{X: tile.X, Y: tile.Y}

	if tile.Type == TileFloor {
		delete(t, key)
		return
	}

	t[key] = tile.Type
}

func (t CustomTiles) ToSlice() []CustomTile {
	result := make([]CustomTile, 0, len(t))

	for key, tileType := range t {
		result = append(result, CustomTile{
			X:    key.X,
			Y:    key.Y,
			Type: tileType,
		})
	}

	return result
}
