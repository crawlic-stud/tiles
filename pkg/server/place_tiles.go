package server

import (
	"encoding/json"
	"net/http"
	"tiles/pkg/db/gen"
	"tiles/pkg/models"
)

type PlaceTilesRequest struct {
	GameID int64               `json:"gameID"`
	Tiles  []models.CustomTile `json:"tiles"`
}

func (h *Handler) PlaceTiles(r *http.Request) Response {
	req, err := ReadBody[PlaceTilesRequest](r)
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "invalid request body: %v", err)
	}

	game, err := h.Store.GetGameWithGrid(r.Context(), req.GameID)
	if err != nil {
		return JSONError(http.StatusBadRequest, err.Error())
	}

	for _, tile := range req.Tiles {
		game.CustomTiles.Set(tile)
	}
	tilesJSON, err := json.Marshal(game.CustomTiles.ToSlice())
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "failed to serialize tiles: %v", err)
	}

	if err = h.Store.UpdateGameTiles(r.Context(), gen.UpdateGameTilesParams{
		CustomTiles: tilesJSON,
		ID:          req.GameID,
	}); err != nil {
		return JSONErrorf(http.StatusInternalServerError, "failed to update game grid: %v", err)
	}

	return JSON(http.StatusOK, nil)
}
