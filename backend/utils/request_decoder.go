package utils

import (
	"encoding/json"
	"net/http"
)

// DecodeJSONBody decodes a JSON request body into the provided destination struct
func DecodeJSONBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

// DecodeFormBody decodes a form request body into a specified form field
func DecodeFormBody(r *http.Request, formField string, dst *string) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	*dst = r.FormValue(formField)
	return nil
}
