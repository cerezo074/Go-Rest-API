package handlers

import (
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

//Place all paths related with user
const (
	UsersPath        string = "/users"
	UsersPathSlashed string = "/users/"
)

//UsersRotuer handles all routes for user path
func UsersRotuer(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")

	if path == UsersPath {
		routeUsers(w, r)
		return
	}

	userIDPath := strings.TrimPrefix(path, UsersPathSlashed)

	if !bson.IsObjectIdHex(userIDPath) {
		postError(w, http.StatusNotFound)
		return
	}

	routeUser(userIDPath, w, r)
}

func routeUsers(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		usersGetAll(writer, request)
		return
	case http.MethodPost:
		return
	default:
		postError(writer, http.StatusMethodNotAllowed)
	}
}

func routeUser(idPath string, writer http.ResponseWriter, request *http.Request) {
	// id := bson.ObjectIdHex(idPath)
	switch request.Method {
	case http.MethodGet:
		return
	case http.MethodPut:
		return
	case http.MethodPatch:
		return
	case http.MethodDelete:
		return
	default:
		postError(writer, http.StatusMethodNotAllowed)
	}
}
