package user

import (
	"errors"

	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
)

//User holds data for a single user
type User struct {
	ID   bson.ObjectId `json:"id" storm:"id"`
	Name string        `json:"name"`
	Role string        `json:"role"`
}

const (
	dbPath string = "users.db"
)

//ErrRecordInvalid used when the integrity of the holded data is not correct(e.g. id not preset)
var (
	ErrRecordInvalid = errors.New("record is invalid")
)

//All fetches all users from the database
func All() ([]User, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	users := []User{}
	err = db.All(&users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

//One gets a user by id
func One(id bson.ObjectId) (*User, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	user := new(User)
	err = db.One("ID", id, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

//Delete removes a given record from the database
func Delete(id bson.ObjectId) error {
	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}

	defer db.Close()
	user := new(User)
	err = db.One("ID", id, user)

	if err != nil {
		return err
	}

	return db.DeleteStruct(user)
}

//Save updates or creates a given record in the database
func (u *User) Save() error {
	if err := u.validate(); err != nil {
		return err
	}

	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}

	defer db.Close()
	return db.Save(u)
}

//validate makes sure a user contains valid data
func (u *User) validate() error {
	if u.Name == "" {
		return ErrRecordInvalid
	}

	return nil
}
