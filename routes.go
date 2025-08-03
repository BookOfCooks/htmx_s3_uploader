package main

import (
	"fmt"
	"main/pox"
	"main/templates"
	"net/http"
	"os"
	"strings"
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

	os.Mkdir("public/audios", 0o700)

	r.Mount("/public/",
		disableCacheInDevMode(
			http.StripPrefix("/public",
				http.FileServer(http.Dir("public")))))

	r.Get("/", pox.Wrap(home))
	r.Post("/form/step1", pox.Wrap(formStep1))
	r.Post("/form/step2", pox.Wrap(formStep2))
	r.Post("/form/step3", pox.Wrap(formStep3))

	return r
}

func home(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	return pox.Templ(http.StatusOK, templates.Home()), nil
}

func formStep1(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	name := strings.TrimSpace(r.FormValue("name"))
	if len(name) == 0 {
		return pox.Templ(http.StatusOK, templates.AlertError("Name cannot be empty")), nil
	}
	return pox.Templ(http.StatusOK, templates.Stepper(2, templates.Step2())), nil
}

func formStep2(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	_, _, err := r.FormFile("file")
	if err != nil {
		return pox.Templ(http.StatusOK, templates.AlertError("Cannot open file, please contact tech support!")), fmt.Errorf("r.FormFile: %w", err)
	}
	return pox.Templ(http.StatusOK, templates.Stepper(3, templates.Step3())), nil
}

func formStep3(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	return pox.Templ(http.StatusOK, templates.Stepper(4, templates.Step4())), nil
}
