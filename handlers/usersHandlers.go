package handlers

import (
	"encoding/json"
	"errors"
	"github/cerezo074/GoAPI/user"
	"io/ioutil"
	"net/http"

	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
)

const (
	usersKey string = "users"
	userKey  string = "user"
)

func bodyToUser(request *http.Request, user *user.User) error {
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
	users, error := user.All()

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	postBodyResponse(writer, http.StatusAccepted, jsonResponse{usersKey: users})
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

	if error == user.ErrorRecordInvalid {
		postError(writer, http.StatusBadRequest)
		return
	}

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Location", UsersPathSlashed+newUser.ID.Hex())
	writer.WriteHeader(http.StatusCreated)
}

func usersGetOne(writer http.ResponseWriter, _ *http.Request, userID bson.ObjectId) {
	user, error := user.One(userID)

	if error == storm.ErrNotFound {
		postError(writer, http.StatusNotFound)
		return
	}

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	postBodyResponse(writer, http.StatusAccepted, jsonResponse{userKey: user})
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

	if error == user.ErrorRecordInvalid {
		postError(writer, http.StatusBadRequest)
		return
	}

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	postBodyResponse(writer, http.StatusAccepted, jsonResponse{userKey: updatedUser})
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

	postBodyResponse(writer, http.StatusAccepted, jsonResponse{userKey: updatedUser})
}

func usersDeleteOne(writer http.ResponseWriter, _ *http.Request, userID bson.ObjectId) {
	error := user.Delete(userID)

	if error == storm.ErrNotFound {
		postError(writer, http.StatusNotFound)
		return
	}

	if error != nil {
		postError(writer, http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
}
