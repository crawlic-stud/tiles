package server

import (
	"encoding/json"
	"net/http"
	"tiles/pkg/db"
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

	// TODO: when moving to postgres - update json directly
	if err = h.Store.UpdateGameTiles(r.Context(), db.UpdateGameTilesParams{
		CustomTiles: string(tilesJSON),
		ID:          req.GameID,
	}); err != nil {
		return JSONErrorf(http.StatusInternalServerError, "failed to update game grid: %v", err)
	}

	// if err = h.hub.BroadcastRerender(req.GameID); err != nil {
	// 	return JSONErrorf(http.StatusInternalServerError, "failed to broadcast rerender: %v", err)
	// }

	return JSON(http.StatusOK, nil)
}
