package database

import (
	"flag"

	"gopkg.in/mgo.v2"
)

//http://blog.denevell.org/golang-constants-enums.html
const (
	cUser = iota
	cPost

	cLast //this is the array length
)

type Database struct {
	session *mgo.Session

	collections [cLast]*mgo.Collection //collections[cUser] will return collection for user
}

var (
	mongoServer = flag.String("mongo-db-server", "127.0.0.1", "IP of mongod.")
)

func New() (*Database, error) {
	session, err := mgo.Dial(*mongoServer)
	if err != nil {
		return nil, err
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	var collections [cLast]*mgo.Collection
	db := session.DB("simple-albums")

	collections[cUser] = db.C("user")
	collections[cPost] = db.C("post")
	//more about MongoDB, see http://edgystuff.tumblr.com/post/93523827905/how-to-implement-robust-and-scalable-transactions

	return &Database{session: session, collections: collections}, nil
}

func (db *Database) CloseSession() {
	db.session.Close()
}
