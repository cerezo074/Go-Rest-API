package handlers

import (
	"github/cerezo074/GoAPI/user"
	"net/http"
)

const (
	usersKey string = "users"
)

func usersGetAll(writer http.ResponseWriter, request *http.Request) {
	users, error := user.All()

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	postBodyResponse(writer, http.StatusAccepted, jsonResponse{usersKey: users})
}
