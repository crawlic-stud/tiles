package models

import (
	"fmt"
)

type Game struct {
	ID              int64
	Grid            Grid
	Width           int
	Height          int
	TileSize        int
	BackgroundImage string
	HideTiles       bool
	CustomTiles     CustomTiles
}

type GridOption func(Grid)

func (g *Game) InitGrid(opts ...GridOption) Grid {
	g.Grid = make(Grid, g.Height)
	for y := 0; y < g.Height; y++ {
		g.Grid[y] = make([]Tile, g.Width)
		for x := 0; x < g.Width; x++ {
			g.Grid.SetTile(x, y, DefaultTile())
		}
	}
	return g.Grid
}

func (g *Game) SetCustomTiles(tiles []CustomTile) error {
	g.CustomTiles = make(CustomTiles, len(tiles))
	for _, tile := range tiles {
		if tile.Y < 0 || tile.Y >= g.Height {
			return fmt.Errorf("tile.Y out of range: %d", tile.Y)
		}
		if tile.X < 0 || tile.X >= g.Width {
			return fmt.Errorf("tile.X out of range: %d", tile.X)
		}
		g.CustomTiles.Set(tile)
		g.Grid.SetTile(tile.X, tile.Y, Tile{Type: tile.Type})
	}
	return nil
}
