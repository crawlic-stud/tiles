package server

import (
	"encoding/json"
	"net/http"
	"tiles/pkg/db/gen"
	"tiles/pkg/models"
)

type CreateGameRequest struct {
	Characters      []models.Character  `json:"characters"`
	CustomTiles     []models.CustomTile `json:"customTiles"`
	BackgroundImage string              `json:"backgroundImage"`
	TileSize        int                 `json:"tileSize"`
	Width           int                 `json:"width"`
	Height          int                 `json:"height"`
}

type CreateGameResponse struct {
	ID int64 `json:"id"`
}

func (h *Handler) CreateGame(r *http.Request) Response {
	req, err := ReadBody[CreateGameRequest](r)
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "invalid request body: %v", err)
	}

	tilesJSON, err := json.Marshal(req.CustomTiles)
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "failed to serialize tiles: %v", err)
	}

	var gameDb gen.Game
	if err = h.Store.WithTx(r.Context(), func(q *gen.Queries) error {
		gameDb, err = q.CreateGame(r.Context(), gen.CreateGameParams{
			Background:  req.BackgroundImage,
			Width:       int32(req.Width),
			Height:      int32(req.Height),
			TileSize:    int32(req.TileSize),
			CustomTiles: tilesJSON,
		})
		if err != nil {
			return err
		}

		for _, character := range req.Characters {
			if err = h.Store.AddCharacterToGameTx(r.Context(), q, character, gameDb.ID); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return JSONErrorf(http.StatusInternalServerError, "failed to create game: %v", err.Error())
	}

	return JSON(http.StatusCreated, CreateGameResponse{ID: gameDb.ID})
}
