package server

import (
	"net/http"
	"tiles/pkg/models"
)

type AddCharacterToGameRequest struct {
	GameID    int64            `json:"gameID"`
	Character models.Character `json:"character"`
}

func (h *Handler) AddCharacterToGame(r *http.Request) Response {
	req, err := ReadBody[AddCharacterToGameRequest](r)
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "invalid request body: %v", err)
	}

	if err = h.Store.AddCharacterToGame(r.Context(), req.Character, req.GameID); err != nil {
		return JSONErrorf(http.StatusInternalServerError, "failed to add character to game: %v", err)
	}

	if err = h.hub.BroadcastRerender(req.GameID); err != nil {
		return JSONErrorf(http.StatusInternalServerError, "failed to broadcast rerender: %v", err)
	}

	return JSON(http.StatusOK, nil)
}
