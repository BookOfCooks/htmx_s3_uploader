package pox

import (
	"encoding/json"
	"main/utils"
	"net/http"
	"os"
	"time"
)

func Wrap(compute func(w http.ResponseWriter, r *http.Request) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := compute(w, r)
		if err != nil {
			utils.NewPrettyJsonEncoder(os.Stdout).Encode(map[string]any{
				"_time": time.Now().Format(time.DateTime),
				"_uri":  r.URL.String(),
				"error": err.Error(),
			})
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if resp != nil {
			if err := EncodeJSON(w, http.StatusOK, resp); err != nil {
				panic(err)
			}
		}
	}
}

func EncodeJSON(w http.ResponseWriter, code int, data interface{}) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}
