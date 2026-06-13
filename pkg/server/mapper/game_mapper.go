package mapper

import (
	"encoding/json"
	"tiles/pkg/db"
	"tiles/pkg/models"
)

func GameFromDB(gameDB db.Game) (game models.Game, err error) {
	game = models.Game{
		ID:              gameDB.ID,
		Width:           int(gameDB.Width),
		Height:          int(gameDB.Height),
		TileSize:        int(gameDB.TileSize),
		BackgroundImage: gameDB.Background,
		HideTiles:       gameDB.HideTiles,
	}
	game.InitGrid()

	var customTiles []models.CustomTile
	if err = json.Unmarshal([]byte(gameDB.CustomTiles), &customTiles); err != nil {
		return
	}
	if err = game.SetCustomTiles(customTiles); err != nil {
		return
	}

	return game, nil
}
