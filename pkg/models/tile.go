package models

import "math"

type TileType string

const (
	TileFloor    TileType = "" // this is default type, so its empty to save on grid json size
	TileDoor     TileType = "DOOR"
	TileObstacle TileType = "OBSTACLE"
	TileWater    TileType = "WATER"
	TileWall     TileType = "WALL"
)

type Tile struct {
	Type      TileType   `json:"type,omitempty"`
	Character *Character `json:"character,omitempty"`
}

func (t Tile) MovementCost() int {
	switch t.Type {
	case TileFloor, TileDoor:
		return 1
	case TileWater:
		return 2
	case TileObstacle:
		return 3
	default:
		return math.MaxInt
	}
}

func (t Tile) IsWalkable() bool {
	return t.MovementCost() != math.MaxInt
}

func DefaultTile() Tile {
	return Tile{
		Type: TileFloor,
	}
}
