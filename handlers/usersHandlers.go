package handlers

import (
	"encoding/json"
	"errors"
	"github/cerezo074/GoAPI/user"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

const (
	usersKey string = "users"
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

	if error != nil {
		if error == user.ErrorRecordInvalid {
			postError(writer, http.StatusBadRequest)
			return
		}

		postError(writer, http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Location", UsersPathSlashed+newUser.ID.Hex())
	writer.WriteHeader(http.StatusCreated)
}
