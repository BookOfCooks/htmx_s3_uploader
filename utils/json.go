package utils

import (
	"encoding/json"
	"io"
)

func NewPrettyJsonEncoder(w io.Writer) *json.Encoder {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc
}
