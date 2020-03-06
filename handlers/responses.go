package handlers

import (
	"encoding/json"
	"net/http"
)

type jsonResponse map[string]interface{}

func postError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func postBodyResponse(writer http.ResponseWriter, code int, content jsonResponse) {
	if content == nil {
		writer.WriteHeader(code)
		writer.Write([]byte(http.StatusText(code)))
		return
	}

	JSON, error := json.Marshal(content)

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(JSON)
}
