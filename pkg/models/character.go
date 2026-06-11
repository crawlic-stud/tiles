package models

type CharacterType string

const (
	CharacterHero  CharacterType = "Hero"
	CharacterEnemy CharacterType = "Enemy"
)

type Character struct {
	ID       int           `json:"id"`
	Name     string        `json:"name"`
	Type     CharacterType `json:"type"`
	Scale    float64       `json:"scale"`
	Position Position      `json:"position"`
	Image    string        `json:"image"`
}
