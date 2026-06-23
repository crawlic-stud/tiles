package db

import (
	"context"
	"errors"
	"tiles/pkg/db/gen"
	"tiles/pkg/models"
)

func (s *Store) AddCharacterToGame(ctx context.Context, character models.Character, gameID int64) error {
	return s.AddCharacterToGameTx(ctx, s.Queries, character, gameID)
}

func (s *Store) AddCharacterToGameTx(ctx context.Context, tx *gen.Queries, character models.Character, gameID int64) error {
	// get game
	game, err := tx.GetGameByID(ctx, gameID)
	if err != nil {
		return err
	}

	// verify position, find a new one if needed
	allCharacters, err := tx.GetGameCharacters(ctx, gameID)
	if err != nil {
		return err
	}
	allPositions := make(map[models.Position]struct{}, len(allCharacters))
	for _, character := range allCharacters {
		allPositions[models.Position{X: int(character.X), Y: int(character.Y)}] = struct{}{}
	}
	found := false
	if _, exists := allPositions[character.Position]; !exists {
		found = true
	} else {
		for x := range game.Width {
			for y := range game.Height {
				newPosition := models.Position{X: int(x), Y: int(y)}
				if _, exists := allPositions[newPosition]; !exists {
					character.Position = newPosition
					found = true
					break
				}
			}
		}
	}
	if !found {
		return errors.New("no free position left in game")
	}

	// if character id is not passed, then its new
	if character.ID == 0 {
		characterDB, err := tx.CreateCharacter(ctx, gen.CreateCharacterParams{
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
	_, err = tx.CreateGameCharacter(ctx, gen.CreateGameCharacterParams{
		GameID:      gameID,
		CharacterID: int64(character.ID),
		X:           int32(character.Position.X),
		Y:           int32(character.Position.Y),
	})
	return err
}
