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
	"github.com/google/uuid"
	"github.com/sanity-io/litter"
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

type FormSession struct {
	Id         string
	Name       string
	Soundtrack string
	Markers    []Marker
}

type Marker struct {
	Name  string
	Start time.Duration
}

var sessions = map[string]*FormSession{}

func home(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	// In production-grade systems, you'll might save this into a Database
	sessionId := uuid.NewString()
	// In real code, don't write to `sessions` without mutex
	sessions[sessionId] = &FormSession{Id: sessionId}
	return pox.Templ(http.StatusOK, templates.Home(sessionId)), nil
}

func formStep1(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	name := strings.TrimSpace(r.FormValue("name"))
	if len(name) == 0 {
		return pox.Templ(http.StatusOK, templates.AlertError("Name cannot be empty")), nil
	}

	sessionId := r.FormValue("sessionId")
	// In real code, protect with a mutex!
	// WARNING! ENSURE `sessionId` EXISTS IN `sessions`, OTHERWISE YOUR SERVER WILL PANIC AND CRASH!!!
	sessions[sessionId].Name = name

	// Passing `name` over here
	return pox.Templ(http.StatusOK, templates.Stepper(2, templates.Step2(sessionId, name))), nil
}

func formStep2(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	file, headers, err := r.FormFile("file")
	if err != nil {
		return pox.Templ(http.StatusOK, templates.AlertError("Cannot open file, please contact tech support!")), fmt.Errorf("r.FormFile: %w", err)
	}

	// In real code, avoid using `headers.Filename` (or at least with safety precautions)
	uploadUrl, err := FakeS3Upload(headers.Filename, file)
	if err != nil {
		return pox.Templ(http.StatusOK, templates.AlertError("Cannot upload file, try again")), fmt.Errorf("FakeS3Upload: %w", err)
	}

	sessionId := r.FormValue("sessionId")
	sessions[sessionId].Soundtrack = uploadUrl

	// Let's log the data we have so far
	litter.Dump(sessions[sessionId])

	// Don't forget to pass `sessionId` to next step
	return pox.Templ(http.StatusOK, templates.Stepper(3, templates.Step3(sessionId))), nil
}

func formStep3(w http.ResponseWriter, r *http.Request) (http.Handler, error) {
	sessionId := r.FormValue("sessionId")
	session := sessions[sessionId]
	session.Markers = []Marker{}

	return pox.Templ(http.StatusOK, templates.Stepper(4,
		// This time we pass extra information, as this is the review step
		// Markers are TODO
		templates.Step4(sessionId, session.Name, session.Soundtrack),
	)), nil
}
