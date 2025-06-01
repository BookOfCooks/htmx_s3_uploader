package pox

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

func ReportRequest(r *http.Request, err error) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(map[string]any{
		"_time": time.Now().Format(time.DateTime),
		"_uri":  r.URL.String(),
		"error": err.Error(),
	})
}

// Returns an http.HandlerFunc that wraps the result of the h.
// h returns a handler that handles a request, and an error that is treated as an internal server error.
//
// There a several cases that can be reached:
//   - handler != nil && err == nil: most common case, `handler` is called to handle the request.
//   - handler == nil && err != nil: 2nd most common case, Wrap will log `err` and respond with an appropriate response.
//   - handler != nil && err != nil: `handler` is called and `err` is logged.
//   - handler == nil && err == nil: nothing happens (could be useful in the case of Websockets)
func Wrap(h func(w http.ResponseWriter, r *http.Request) (http.Handler, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler, err := h(w, r)
		if handler != nil {
			handler.ServeHTTP(w, r)
		}

		if err != nil {
			ReportRequest(r, err)
			if handler == nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}
}
