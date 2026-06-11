package mapper

import (
	"tiles/pkg/db"
	"tiles/pkg/models"
)

func CharacterFromDb(character db.Character) *models.Character {
	return &models.Character{
		ID:    int(character.ID),
		Name:  character.Name,
		Type:  models.CharacterType(character.Type),
		Scale: character.Scale,
		Image: character.Image,
	}
}

func GameCharacterFromDb[T db.GetGameCharactersRow | db.GetGameCharacterByIDRow](value T) *models.Character {
	character := db.GetGameCharactersRow(value)
	return &models.Character{
		ID:    int(character.ID),
		Name:  character.Name,
		Type:  models.CharacterType(character.Type),
		Scale: character.Scale,
		Position: models.Position{
			X: int(character.X),
			Y: int(character.Y),
		},
		Image: character.Image,
	}
}
