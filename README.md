# Tiles

Tiles is a D&D inspired game, in which you can create maps, put characters on them and move them around. It's good for representing visuals for D&D sessions.

<img width="2550" height="1892" alt="image" src="https://github.com/user-attachments/assets/171510fe-e424-4812-90d4-0ef759bb5487" />

## Run

Simplest and fastest way to run is:

```sh
sqlc generate && templ generate && go run main.go
```

## Create game

To create a simple map you can play around with use this CURL:

```sh
curl --request POST \
  --url http://localhost:3000/games \
  --header 'content-type: application/json' \
  --data '{
  "tileSize": 50,
  "width": 34,
  "height": 88,
  "customTiles": [],
  "characters": [
    {
      "id": 0,
      "name": "Hero",
      "type": "Hero",
      "scale": 1.5,
      "position": {
        "x": 5,
        "y": 5
      },
      "image": "/assets/characters/hero.png"
    },
    {
      "id": 0,
      "name": "Enemy",
      "type": "Enemy",
      "scale": 2,
      "position": {
        "x": 10,
        "y": 10
      },
      "image": "/assets/characters/skeleton.png"
    }
  ],
  "backgroundImage": "/assets/maps/bigmap.jpg"
}'
```
