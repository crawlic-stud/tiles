package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tiles/pkg/db"
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

	fmt.Println(req)

	game := models.Game{
		Width:    req.Width,
		Height:   req.Height,
		TileSize: req.TileSize,
	}
	grid := game.InitGrid()

	err = game.AddCustomTiles(req.CustomTiles)
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "failed to add custom tiles: %v", err)
	}

	gridJSON, err := json.Marshal(grid)
	if err != nil {
		return JSONErrorf(http.StatusBadRequest, "failed to serialize grid: %v", err)
	}

	var gameDb db.Game
	if err = h.Store.WithTx(r.Context(), func(q *db.Queries) error {
		gameDb, err = q.CreateGame(r.Context(), db.CreateGameParams{
			Background: req.BackgroundImage,
			Grid:       string(gridJSON),
			TileSize:   int64(req.TileSize),
		})
		if err != nil {
			return err
		}

		for _, character := range req.Characters {
			// if character id is not passed, then its new
			if character.ID == 0 {
				characterDB, err := q.CreateCharacter(r.Context(), db.CreateCharacterParams{
					Name:  character.Name,
					Type:  string(character.Type),
					Scale: character.Scale,
					Image: character.Image,
				})
				if err != nil {
					return err
				}
				character.ID = int(characterDB.ID)
			}

			// assign character id to the game
			_, err = q.CreateGameCharacter(r.Context(), db.CreateGameCharacterParams{
				GameID:      gameDb.ID,
				CharacterID: int64(character.ID),
				X:           int64(character.Position.X),
				Y:           int64(character.Position.Y),
			})
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return JSONErrorf(http.StatusInternalServerError, "failed to create game: %v", err.Error())
	}

	return JSON(http.StatusCreated, CreateGameResponse{ID: gameDb.ID})
}
