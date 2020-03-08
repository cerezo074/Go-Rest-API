package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
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

func postOptionsResponse(writer http.ResponseWriter, allowedMethods []string, body jsonResponse) {
	writer.Header().Set("Allow", strings.Join(allowedMethods, ","))
	postBodyResponse(writer, http.StatusAccepted, body)
}

func postMethodNotAllowed(writer http.ResponseWriter, allowedMethods []string) {
	writer.Header().Set("Allow", strings.Join(allowedMethods, ","))
	code := http.StatusMethodNotAllowed
	http.Error(writer, http.StatusText(code), code)
}
