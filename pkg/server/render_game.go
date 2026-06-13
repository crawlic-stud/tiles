package server

import (
	"fmt"
	"net/http"
	"strconv"
	"tiles/pkg/server/mapper"
	"tiles/templates"
)

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

	game, err := h.Store.GetGameWithGrid(r.Context(), gameID)
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
		game.Grid.SetCharacter(int(character.X), int(character.Y), mapper.GameCharacterFromDB(character))
	}

	return HTML(templates.Main(game))
}
