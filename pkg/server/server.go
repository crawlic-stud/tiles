package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"tiles/pkg/store"
)

func Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, err := store.New("tiles.db")
	if err != nil {
		panic(err)
	}
	app := New(storage)

	go app.RunSocketHub(ctx)

	app.Handle("POST /games", app.CreateGame)
	app.Handle("GET /games/{id}", app.RenderGame)
	app.Handle("POST /games/tiles", app.PlaceTiles)

	app.Handle("POST /settings/hideTiles", app.HideTiles)

	app.Handle("POST /characters/move", app.MoveCharacter)

	http.HandleFunc("/ws/{id}", app.ConnectToHub)

	// handle for static files
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	port := "3000"
	fmt.Println("Listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
