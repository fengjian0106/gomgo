package database

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Comment struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	From        UserEmbed
	Message     string
	CreatedTime time.Time
}

type Post struct {
	Id          bson.ObjectId `json:"id"    bson:"_id"`
	From        UserEmbed     `json:"from"  bson:"from"`
	Message     string
	CreatedTime time.Time
	UpdatedTime time.Time
	Comments    []Comment
}

func (db *Database) CreatePost(post *Post) error {
	pc := db.collections[cPost]
	err := pc.Insert(post)

	if err != nil {
		return &DbError{DbErrorInsertErr, err}
	}

	return nil
}

func (db *Database) GetPostByPostId(id bson.ObjectId) (*Post, error) {
	pc := db.collections[cPost]

	post := Post{}
	err := pc.Find(bson.M{"_id": id}).One(&post) //err := pc.FindId(id).One(&post)

	if err != nil {
		return nil, &DbError{DbErrorFindErr, err}
	}

	return &post, nil
}

func (db *Database) GetPostsByUserId(id bson.ObjectId) ([]*Post, error) {
	pc := db.collections[cPost]

	var posts []*Post

	//http://stackoverflow.com/questions/10043965/how-to-get-a-specific-embedded-document-inside-a-mongodb-collection
	err := pc.Find(bson.M{"from._id": id}).All(&posts)

	if err != nil {
		return nil, &DbError{DbErrorFindErr, err}
	}
	//log.Println(posts)

	return posts, nil
}

func (db *Database) CreateCommentForPostId(comment *Comment, postId bson.ObjectId) error {
	//https://gist.github.com/border/3489566
	pc := db.collections[cPost]
	colQuerier := bson.M{"_id": postId}
	change := bson.M{"$push": bson.M{"comments": comment}}
	err := pc.Update(colQuerier, change)
	if err != nil {
		return &DbError{DbErrorUpdateErr, err}
	}

	return nil
}
