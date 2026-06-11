package server

import (
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

	game, err := h.getGameWithGrid(r.Context(), req.GameID)
	if err != nil {
		return JSONError(http.StatusBadRequest, err.Error())
	}

	if err = game.AddCustomTiles(req.Tiles); err != nil {
		return JSONError(http.StatusBadRequest, err.Error())
	}

	gridJSON, err := game.GetGridJSON()
	if err != nil {
		return JSONErrorf(http.StatusInternalServerError, "corrupted JSON: %v", err)
	}

	// TODO: when moving to postgres - update json directly
	if err = h.Store.UpdateGameGrid(r.Context(), db.UpdateGameGridParams{
		Grid: gridJSON,
		ID:   req.GameID,
	}); err != nil {
		return JSONErrorf(http.StatusInternalServerError, "failed to update game grid: %v", err)
	}

	// if err = h.hub.BroadcastRerender(req.GameID); err != nil {
	// 	return JSONErrorf(http.StatusInternalServerError, "failed to broadcast rerender: %v", err)
	// }

	return JSON(http.StatusOK, nil)
}
