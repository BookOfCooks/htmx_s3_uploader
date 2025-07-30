package main

import (
	"main/pox"
	"main/templates"
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

	r.Mount("/public/",
		disableCacheInDevMode(
			http.StripPrefix("/public",
				http.FileServer(http.Dir("public")))))

	r.Get("/", pox.Wrap(home))

	return r
}

func home(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	return pox.Templ(http.StatusOK, templates.Home()), nil
}

func disableCacheInDevMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
