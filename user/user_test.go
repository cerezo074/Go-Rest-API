package user

import (
	"os"
	"reflect"
	"testing"

	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
)

func TestMain(testSuite *testing.M) {
	testSuite.Run()
	os.Remove(dbPath)
}

func BenchmarkCRUD(testCase *testing.B) {
	os.Remove((dbPath))
	testCase.ResetTimer()

	for i := 0; i < testCase.N; i++ {

	}
}

func TestCRUD(testCase *testing.T) {
	///////////////////
	testCase.Log("Create")

	user := &User{
		ID:   bson.NewObjectId(),
		Name: "John",
		Role: "Tester",
	}

	error := user.Save()

	if error != nil {
		testCase.Errorf("Error saving a record: %s", error)
	}

	///////////////////
	testCase.Log("Read")

	user2, error := One(user.ID)

	if error != nil {
		testCase.Errorf("Error retrieving a record %s", error)
	}

	if !reflect.DeepEqual(user2, user) {
		testCase.Error("Records do not match")
	}

	///////////////////
	testCase.Log("Update")

	user.Role = "Developer"
	error = user.Save()

	if error != nil {
		testCase.Errorf("Error saving a record: %s", error)
	}

	user3, error := One(user.ID)

	if error != nil {
		testCase.Errorf("Error retrieving a record %s", error)
	}

	if !reflect.DeepEqual(user3, user) {
		testCase.Error("Record does not match")
	}

	///////////////////
	testCase.Log("Delete")

	error = Delete(user.ID)

	if error != nil {
		testCase.Errorf("Error deleting a record: %s", error)
	}

	_, error = One(user.ID)

	if error == nil {
		testCase.Error("Record should not exist anymore")
	}

	if error != storm.ErrNotFound {
		testCase.Errorf("Error retrieving non-existing record: %s", error)
	}

	///////////////////
	testCase.Log("Read All")

	user.ID = bson.NewObjectId()
	user2.ID = bson.NewObjectId()
	user3.ID = bson.NewObjectId()

	error = user.Save()

	if error != nil {
		testCase.Errorf("Error saving a record: %s", error)
	}

	error = user2.Save()

	if error != nil {
		testCase.Errorf("Error saving a record: %s", error)
	}

	error = user3.Save()

	if error != nil {
		testCase.Errorf("Error saving a record: %s", error)
	}

	users, error := All()

	if error != nil {
		testCase.Errorf("Error reading all records: %s", error)
	}

	if len(users) != 3 {
		testCase.Error("All records in database should be 3")
	}
}
