package mapper

import "tiles/pkg/models"

func MapPath(path []models.Position) [][2]int {
	mapped := make([][2]int, len(path))
	for i, pos := range path {
		mapped[i] = [2]int{pos.X, pos.Y}
	}
	return mapped
}
