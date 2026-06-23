package mapper

import (
	"tiles/pkg/db/gen"
	"tiles/pkg/models"
)

func CharacterFromDB(character gen.Character) *models.Character {
	return &models.Character{
		ID:    int(character.ID),
		Name:  character.Name,
		Type:  models.CharacterType(character.Type),
		Scale: character.Scale,
		Image: character.Image,
	}
}

func GameCharacterFromDB[T gen.GetGameCharactersRow | gen.GetGameCharacterByIDRow](value T) *models.Character {
	character := gen.GetGameCharactersRow(value)
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
