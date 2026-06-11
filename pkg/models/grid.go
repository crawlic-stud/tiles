package models

type Grid [][]Tile

func (g Grid) GetTile(x, y int) Tile {
	return g[y][x]
}

func (g Grid) SetTile(x, y int, tile Tile) {
	g[y][x] = tile
}

func (g Grid) SetCharacter(x, y int, c *Character) {
	g[y][x].Character = c
}

func (g Grid) Width() int {
	return len(g[0])
}

func (g Grid) Height() int {
	return len(g)
}
