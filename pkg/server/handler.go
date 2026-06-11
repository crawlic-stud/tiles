package server

import (
	"context"
	"encoding/json"
	"net/http"
	"tiles/pkg/store"
)

type Handler struct {
	Store *store.Store
	hub   *GameHub
}

func New(store *store.Store) *Handler {
	return &Handler{
		Store: store,
		hub:   NewHub(),
	}
}

type Response interface {
	Write(w http.ResponseWriter, r *http.Request)
}

func ReadBody[T any](r *http.Request) (body T, err error) {
	err = json.NewDecoder(r.Body).Decode(&body)
	return
}

func (h *Handler) Handle(
	pattern string,
	fn func(r *http.Request) Response,
) {
	f := func(w http.ResponseWriter, r *http.Request) {
		resp := fn(r)
		resp.Write(w, r)
	}
	http.HandleFunc(pattern, f)
}

func (h *Handler) RunSocketHub(ctx context.Context) {
	h.hub.Run(ctx)
}
