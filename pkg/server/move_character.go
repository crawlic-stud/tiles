package server

import (
	"encoding/json"
	"log"
	"net/http"
	"tiles/pkg/db"
	"tiles/pkg/game"
	"tiles/pkg/models"
	"tiles/pkg/server/mapper"
)

type MoveCharacterRequest struct {
	ID     int64 `json:"id"`
	GameID int64 `json:"gameID"`
	X      int64 `json:"x"`
	Y      int64 `json:"y"`
}

type MoveCharacterResponse struct {
	ID   int64    `json:"id"`
	Path [][2]int `json:"path"`
}

func (h *Handler) MoveCharacter(r *http.Request) Response {
	var req MoveCharacterRequest
	req, err := ReadBody[MoveCharacterRequest](r)
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "invalid request body: %v", err)
	}

	log.Printf("moving character %d in game %d", req.ID, req.GameID)

	gameDb, err := h.Store.GetGameByID(r.Context(), req.GameID)
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "cant find game with ID=%d: %v", req.GameID, err)
	}

	characterDb, err := h.Store.GetGameCharacterByID(r.Context(), db.GetGameCharacterByIDParams{
		ID:     req.ID,
		GameID: req.GameID,
	})
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "cant find character with ID=%d and gameID=%d: %v", req.ID, req.GameID, err)
	}

	var grid models.Grid
	err = json.Unmarshal([]byte(gameDb.Grid), &grid)
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "failed to parse grid: %v", err)
	}

	character := mapper.GameCharacterFromDb(characterDb)
	path := game.CalculateAndSavePath(grid, character, models.Position{
		X: int(req.X),
		Y: int(req.Y),
	})
	if len(path) == 0 {
		return JSONErrorf(http.StatusBadRequest, "failed to calculate path")
	}

	err = h.Store.UpdateCharacterPosition(r.Context(), db.UpdateCharacterPositionParams{
		X:           req.X,
		Y:           req.Y,
		GameID:      req.GameID,
		CharacterID: characterDb.ID,
	})
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "failed to update character position: %v", err)
	}

	data := MoveCharacterResponse{Path: mapper.MapPath(path), ID: req.ID}
	if err = h.hub.BroadcastJSON(gameDb.ID, MessageTypeMove, data); err != nil {
		return JSONErrorf(http.StatusInternalServerError, "failed to broadcast: %v", err)
	}

	return JSONResponse{
		StatusCode: http.StatusOK,
		Data:       data,
	}
}
