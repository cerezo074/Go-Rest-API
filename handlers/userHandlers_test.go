package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github/cerezo074/GoAPI/user"

	"gopkg.in/mgo.v2/bson"
)

type UserTestCase struct {
	title       string
	request     *http.Request
	userInput   *user.User
	expectError bool
	userOutput  *user.User
}

func TestBodyToUser(test *testing.T) {
	validUser := &user.User{
		ID:   bson.NewObjectId(),
		Name: "John",
		Role: "Tester",
	}

	validUser2 := &user.User{
		ID:   validUser.ID,
		Name: "John",
		Role: "Developer",
	}

	validJSON, error := json.Marshal(validUser)

	if error != nil {
		test.Errorf("Error marshalling a valid user %s", error)
		test.FailNow()
	}

	ts := []UserTestCase{
		{
			title:       "nil request",
			expectError: true,
		},
		{
			title:       "empty request body",
			request:     &http.Request{},
			expectError: true,
		},
		{
			title: "empty user",
			request: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString("{}")),
			},
			expectError: true,
		},
		{
			title: "malformated data",
			request: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"id": "12"}`)),
			},
			userInput:   &user.User{},
			expectError: true,
		},
		{
			title: "valid request",
			request: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBuffer(validJSON)),
			},
			userInput:  &user.User{},
			userOutput: validUser,
		}, {
			title: "valid partial request",
			request: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"role": "Developer", "age": 22}`)),
			},
			userInput:  validUser,
			userOutput: validUser2,
		},
	}

	for _, testCase := range ts {
		test.Log(testCase.title)
		error := bodyToUser(testCase.request, testCase.userInput)

		if testCase.expectError && error == nil {
			test.Error("Expect to find an error")
			continue
		}

		if !testCase.expectError && error != nil {
			test.Errorf("Error found! %s", error)
			continue
		}

		if error == nil && !reflect.DeepEqual(testCase.userInput, testCase.userOutput) {
			test.Error("Unmarshalled data is different: ")
			test.Error("Input: ", testCase.userInput)
			test.Error("Output: ", testCase.userOutput)
		}
	}
}
