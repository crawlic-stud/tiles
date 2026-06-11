package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tiles/pkg/models"
	"tiles/pkg/server/mapper"
	"tiles/templates"
)

func (h *Handler) getGameWithGrid(ctx context.Context, gameID int64) (game models.Game, err error) {
	// get game by id
	gameDB, err := h.Store.GetGameByID(ctx, gameID)
	if err != nil {
		return game, fmt.Errorf("cant find game with ID=%d: %v", gameID, err)
	}

	// marshal grid into structure
	var grid models.Grid
	err = json.Unmarshal([]byte(gameDB.Grid), &grid)
	if err != nil {
		return game, fmt.Errorf("failed to parse grid: %v", err)
	}

	game = models.Game{
		Grid:            grid,
		TileSize:        int(gameDB.TileSize),
		Width:           grid.Width(),
		Height:          grid.Height(),
		BackgroundImage: gameDB.Background,
		HideTiles:       gameDB.HideTiles,
	}
	return game, nil
}

func getGameIDFromPath(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	gameID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid game id: %v", err)
	}
	return gameID, nil
}

func (h *Handler) RenderGame(r *http.Request) Response {
	// get id from request path
	gameID, err := getGameIDFromPath(r)
	if err != nil {
		return HTMLError(http.StatusBadRequest, err.Error())
	}

	game, err := h.getGameWithGrid(r.Context(), gameID)
	if err != nil {
		return HTMLError(http.StatusBadRequest, err.Error())
	}
	game.ID = gameID

	// get characters that are in game
	gameCharacters, err := h.Store.GetGameCharacters(r.Context(), gameID)
	if err != nil {
		return HTMLErrorf(http.StatusNotFound, "cant find game characters for game=%d: %v", gameID, err)
	}
	for _, character := range gameCharacters {
		game.Grid.SetCharacter(int(character.X), int(character.Y), mapper.GameCharacterFromDb(character))
	}

	return HTML(templates.Main(game))
}
