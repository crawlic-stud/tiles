package server

import (
	"net/http"
	"tiles/pkg/db"
)

type HideTilesRequest struct {
	GameID int64 `json:"gameID"`
	Hide   bool  `json:"hide"`
}

func (h *Handler) HideTiles(r *http.Request) Response {
	req, err := ReadBody[HideTilesRequest](r)
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "invalid request body: %v", err)
	}

	if err = h.Store.HideGameTiles(r.Context(), db.HideGameTilesParams{
		HideTiles: req.Hide,
		ID:        req.GameID,
	}); err != nil {
		return JSONErrorf(http.StatusInternalServerError, "failed to hide tiles: %v", err)
	}

	if err = h.hub.BroadcastRerender(req.GameID); err != nil {
		return JSONErrorf(http.StatusInternalServerError, "failed to broadcast rerender: %v", err)
	}

	return JSON(http.StatusOK, nil)
}
