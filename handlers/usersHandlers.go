package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"projects/GoAPI/cache"
	"projects/GoAPI/user"

	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
)

const (
	usersKey string = "users"
	userKey  string = "user"
)

func bodyToUser(request *http.Request, user *user.User) error {
	if request == nil {
		return errors.New("Request is required")
	}

	if request.Body == nil {
		return errors.New("Request body is empty")
	}

	if user == nil {
		return errors.New("A user is required")
	}

	rawBody, err := ioutil.ReadAll(request.Body)

	if err != nil {
		return err
	}

	return json.Unmarshal(rawBody, user)
}

func usersGetAll(writer http.ResponseWriter, request *http.Request) {
	if cache.Serve(writer, request) {
		return
	}

	users, error := user.All()

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	if request.Method == http.MethodHead {
		postBodyResponse(writer, http.StatusAccepted, jsonResponse{})
		return
	}

	customWriter := cache.NewWritter(writer, request)
	postBodyResponse(customWriter, http.StatusAccepted, jsonResponse{usersKey: users})
}

func usersPostOne(writer http.ResponseWriter, request *http.Request) {
	newUser := new(user.User)
	error := bodyToUser(request, newUser)

	if error != nil {
		postError(writer, http.StatusBadRequest)
		return
	}

	newUser.ID = bson.NewObjectId()
	error = newUser.Save()

	if error == user.ErrRecordInvalid {
		postError(writer, http.StatusBadRequest)
		return
	}

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	cache.Drop(UsersPathSlashed)
	writer.Header().Set("Location", UsersPathSlashed+newUser.ID.Hex())
	writer.WriteHeader(http.StatusCreated)
}

func usersGetOne(writer http.ResponseWriter, request *http.Request, userID bson.ObjectId) {
	if cache.Serve(writer, request) {
		return
	}

	user, error := user.One(userID)

	if error == storm.ErrNotFound {
		postError(writer, http.StatusNotFound)
		return
	}

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	if request.Method == http.MethodHead {
		postBodyResponse(writer, http.StatusAccepted, jsonResponse{})
		return
	}

	customWriter := cache.NewWritter(writer, request)
	postBodyResponse(customWriter, http.StatusAccepted, jsonResponse{userKey: user})
}

func usersPutOne(writer http.ResponseWriter, request *http.Request, userID bson.ObjectId) {
	updatedUser := new(user.User)
	error := bodyToUser(request, updatedUser)

	if error != nil {
		postError(writer, http.StatusBadRequest)
		return
	}

	updatedUser.ID = userID
	error = updatedUser.Save()

	if error == user.ErrRecordInvalid {
		postError(writer, http.StatusBadRequest)
		return
	}

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	cache.Drop(UsersPathSlashed)
	customWriter := cache.NewWritter(writer, request)
	postBodyResponse(customWriter, http.StatusAccepted, jsonResponse{userKey: updatedUser})
}

func usersPatchOne(writer http.ResponseWriter, request *http.Request, userID bson.ObjectId) {
	updatedUser, error := user.One(userID)

	if error == storm.ErrNotFound {
		postError(writer, http.StatusNotFound)
		return
	}

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	error = bodyToUser(request, updatedUser)

	if error != nil {
		postError(writer, http.StatusBadRequest)
		return
	}

	error = updatedUser.Save()

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	cache.Drop(UsersPathSlashed)
	customWriter := cache.NewWritter(writer, request)
	postBodyResponse(customWriter, http.StatusAccepted, jsonResponse{userKey: updatedUser})
}

func usersDeleteOne(writer http.ResponseWriter, request *http.Request, userID bson.ObjectId) {
	error := user.Delete(userID)

	if error == storm.ErrNotFound {
		postError(writer, http.StatusNotFound)
		return
	}

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	cache.Drop(UsersPathSlashed)
	cache.Drop(cache.MakeResource(request))
	writer.WriteHeader(http.StatusAccepted)
}
