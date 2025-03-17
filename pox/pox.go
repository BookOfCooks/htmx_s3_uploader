package pox

import (
	"encoding/json"
	"main/utils"
	"net/http"
	"os"
	"time"
)

// if err != nil: client receives 500 response, and the error is logged
// else if err == nil && resp != nli: resp is encoded and sent to client
// else if err == nill && resp == nil: does nothing
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
			if v, ok := resp.(statusOk); ok {
				if err := EncodeJSON(w, http.StatusOK, v.data); err != nil {
					panic(err)
				}
			} else if _, ok := resp.(statusNotFound); ok {
				EncodeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			} else {
				if err := EncodeJSON(w, http.StatusOK, resp); err != nil {
					panic(err)
				}
			}

		}
	}
}

func EncodeJSON(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}
