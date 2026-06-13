package store

import (
	"context"
	"fmt"
	"tiles/pkg/models"
	"tiles/pkg/server/mapper"
)

func (s *Store) GetGameWithGrid(ctx context.Context, gameID int64) (game models.Game, err error) {
	gameDB, err := s.GetGameByID(ctx, gameID)
	if err != nil {
		return game, fmt.Errorf("cant find game with ID=%d: %v", gameID, err)
	}

	game, err = mapper.GameFromDB(gameDB)
	if err != nil {
		return game, fmt.Errorf("failed to get game from db: %v", err)
	}

	return game, nil
}
