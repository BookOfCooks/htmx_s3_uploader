package main

import (
	"main/pox"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", pox.Wrap(home))

	return r
}

func home(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	return pox.JSON(http.StatusOK, map[string]any{
		"hello": "world",
	}), nil
}
