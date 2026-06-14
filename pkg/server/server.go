package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"tiles/pkg/store"
)

const assetsDir = "assets"

func Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, err := store.New("tiles.db")
	if err != nil {
		panic(err)
	}
	app := New(storage)

	go app.RunSocketHub(ctx)

	// games handles
	app.Handle("POST /games", app.CreateGame)
	app.Handle("GET /games/{id}", app.RenderGame)
	app.Handle("POST /games/tiles", app.PlaceTiles)
	app.Handle("POST /games/characters", app.AddCharacterToGame)

	// settings handles
	app.Handle("POST /settings/hideTiles", app.HideTiles)

	// characters handles
	app.Handle("POST /characters/move", app.MoveCharacter)

	// images handles
	app.Handle("POST /images", app.UploadImage)

	http.HandleFunc("/ws/{id}", app.ConnectToHub)

	// handle for static files
	assetsRoute := "/" + assetsDir + "/"
	http.Handle(assetsRoute, http.StripPrefix(assetsRoute, http.FileServer(http.Dir(assetsDir))))

	// create uploads dir for user uploads
	if err := os.MkdirAll(assetsDir+"/uploads", 0755); err != nil {
		panic(err)
	}

	port := "3000"
	fmt.Println("Listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
