package handlers

import "net/http"

//HomePath is the root route
const (
	HomePath string = "/"
)

//RootHandler Handles the root / route
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != HomePath {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Asset not found \n"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Running API V1\n"))
}
