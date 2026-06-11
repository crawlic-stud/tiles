package game

import (
	"tiles/pkg/store"
)

type Game struct {
	Store *store.Store
}

func New(store *store.Store) *Game {
	return &Game{
		Store: store,
	}
}
