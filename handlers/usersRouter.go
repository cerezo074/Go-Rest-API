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
		usersPostOne(writer, request)
		return
	default:
		postError(writer, http.StatusMethodNotAllowed)
	}
}

func routeUser(idPath string, writer http.ResponseWriter, request *http.Request) {
	userID := bson.ObjectIdHex(idPath)
	switch request.Method {
	case http.MethodGet:
		usersGetOne(writer, request, userID)
		return
	case http.MethodPut:
		usersPutOne(writer, request, userID)
		return
	case http.MethodPatch:
		usersPatchOne(writer, request, userID)
		return
	case http.MethodDelete:
		usersDeleteOne(writer, request, userID)
		return
	default:
		postError(writer, http.StatusMethodNotAllowed)
	}
}
