package pox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled(), UseJSONTagNameAsFieldName)
	Mod      = modifiers.New()
)

// Copied from https://github.com/go-playground/validator/issues/258#issuecomment-257281334
func UseJSONTagNameAsFieldName(validate *validator.Validate) {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Decodes the request body as JSON, modifies it according to the `mod` tag, then validates it.
//
// It returns an error in case of failure, or the value
func Process[T any](w http.ResponseWriter, r *http.Request, value *T) (T, bool) {
	if !processJSON(w, r, value) {
		return *new(T), false
	} else if !processMod(w, r, value) {
		return *new(T), false
	} else if !processValidate(w, r, *value) {
		return *new(T), false
	} else {
		return *value, true
	}
}

func processJSON[T any](w http.ResponseWriter, r *http.Request, value *T) bool {
	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		JSON(http.StatusBadRequest, map[string]any{"cause": "invalid_json_body"}).ServeHTTP(w, r)
		ReportRequest(r, fmt.Errorf("json.Decode: %w", err))
		return false
	}
	return true
}

func processMod[T any](w http.ResponseWriter, r *http.Request, value *T) bool {
	if err := Mod.Struct(r.Context(), value); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ReportRequest(r, fmt.Errorf("Mod.Struct: %w", err))
		return false
	}
	return true
}

func processValidate(w http.ResponseWriter, r *http.Request, value any) bool {
	var errs validator.ValidationErrors
	if err := validate.StructCtx(r.Context(), value); err == nil {
		return true
	} else {
		errs = err.(validator.ValidationErrors)
	}

	mmap := map[string]any{}
	for _, err := range errs {
		mmap[err.Field()] = err.Tag()
	}

	JSON(http.StatusUnprocessableEntity, mmap).ServeHTTP(w, r)
	return false
}
