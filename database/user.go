package database

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

type UserEmbed struct {
	Id   bson.ObjectId `json:"id" bson:"_id"`
	Name string        `json:"name" bson:"name"`
}

type User struct {
	Id       bson.ObjectId `json:"id"                  bson:"_id"`
	Email    string        `json:"email"               bson:"email"`
	Password string        `json:"password,omitempty"  bson:"password"` //when User convert to json string, password should be skipped
	Name     string        `json:"name"                bson:"name"`
	//http://stackoverflow.com/questions/17306358/golang-removing-fields-from-struct-or-hiding-them-in-json-response
}

func (user User) verify() bool {
	//TODO, strict Verify
	if user.Email != "" && user.Password != "" && user.Name != "" {
		return true
	} else {
		return false
	}
}

func (db *Database) CheckEmailUnique(email string) bool {
	uc := db.collections[cUser]

	user := User{}
	err := uc.Find(bson.M{"email": email}).One(&user)
	if err != nil {
		//if email not found, also will enter this if branch
		return true
	}

	//FIXME user.Email == "" when should happen
	if user.Email == "" {
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////
func (db *Database) CreateUser(user *User) error {
	if user.verify() == false {
		return &DbError{DbErrorModalPrptyErr, errors.New("user property error")}
		//TODO, return more meaningful error info
	}

	uc := db.collections[cUser]
	err := uc.Insert(&user)
	if err != nil {
		return &DbError{DbErrorInsertErr, err}
	}

	return nil
}

func (db *Database) GetUserById(id bson.ObjectId) (*User, error) {
	uc := db.collections[cUser]

	user := User{}
	err := uc.Find(bson.M{"_id": id}).One(&user) // err := uc.FindId(id).One(&user)

	if err != nil {
		return nil, &DbError{DbErrorFindErr, err}
	}

	return &user, nil
}

func (db *Database) GetUserByEmail(email string) (*User, error) {
	uc := db.collections[cUser]

	user := User{}
	err := uc.Find(bson.M{"email": email}).One(&user)

	if err != nil {
		return nil, &DbError{DbErrorFindErr, err}
	}

	return &user, nil
}

func (db *Database) GetUsers() ([]*User, error) {
	uc := db.collections[cUser]

	var users []*User
	err := uc.Find(bson.M{}).All(&users)

	if err != nil {
		return nil, &DbError{DbErrorFindErr, err}
	}

	return users, nil
}

/** TODO
func (db *Database) UpdateUser(user *User) error {
	//Can not update password
	return nil
}
*/
