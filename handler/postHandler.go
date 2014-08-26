package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"time"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	"fmt"

	"github.com/fengjian0106/gomgo/appcontext"
	"github.com/fengjian0106/gomgo/database"
)

func CreatePostHandler(appCtx *appcontext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	//log.Printf("postHandler.go, CreateposteatePostSigninHandler")
	//log.Println(r.Header["Content-Type"])

	//<1> get token string
	tokenStr, err := parseJwtTokenStrFromHeaderOrUrlQuery(r.Header, r.URL)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorTokenNotFound, err}
	}

	//<2> jwt token verify, parse to UserEmbed
	var userEmbed *database.UserEmbed
	userEmbed, err = parseJwtTokenStrToUserEmbed(tokenStr)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorTokenParseErr, err}
	}

	//log.Printf("*************** %#v", userEmbed)

	//<3> parse form and get post content created by http client
	type postHttpClientParams struct {
		From    database.UserEmbed
		Message string
	}
	//http://stackoverflow.com/questions/15672556/handling-json-post-request-in-go
	var clientPost postHttpClientParams
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&clientPost)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, err}
	}

	if clientPost.From.Id.Hex() == "" {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("post must have a creater user id")}
	}
	if clientPost.From.Id.Hex() != userEmbed.Id.Hex() {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("no permission for create post with other's id")}
	}

	if clientPost.Message == "" {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("message of post can not be nil")}
	}
	//http://stackoverflow.com/questions/12668681/go-language-string-length
	if utf8.RuneCountInString(clientPost.Message) > 1024 {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("message of post can not longger than 1024 unicode string")}
	}

	clientPost.From.Name = userEmbed.Name

	//<4> insert to db
	post := database.Post{bson.NewObjectId(),
		clientPost.From,
		clientPost.Message,
		time.Now(),
		time.Now(),
		nil}
	err = appCtx.Db.CreatePost(&post)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "application/json")

	jsonStr := fmt.Sprintf("{\"postId\": \"%s\"}", post.Id.Hex())
	w.Write([]byte(jsonStr))
	return http.StatusOK, nil
}

func GetPostByPostIdHandler(appCtx *appcontext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	//<1> get token string
	tokenStr, err := parseJwtTokenStrFromHeaderOrUrlQuery(r.Header, r.URL)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorTokenNotFound, err}
	}

	//<2> jwt token verify, parse to UserEmbed
	_, err = parseJwtTokenStrToUserEmbed(tokenStr)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorTokenParseErr, err}
	}

	//<3>
	vars := mux.Vars(r)
	id := vars["postId"]

	if bson.IsObjectIdHex(id) == false {
		return http.StatusBadRequest, &ApiError{ApiErrorParamIdFmtErr, errors.New("wrong format id")}
	}

	post, err := appCtx.Db.GetPostByPostId(bson.ObjectIdHex(id))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&post)

	return http.StatusOK, nil
}

func GetPostsByUserIdHandler(appCtx *appcontext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	//<1> get token string
	tokenStr, err := parseJwtTokenStrFromHeaderOrUrlQuery(r.Header, r.URL)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorTokenNotFound, err}
	}

	//<2> jwt token verify, parse to UserEmbed
	_, err = parseJwtTokenStrToUserEmbed(tokenStr)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorTokenParseErr, err}
	}

	//<3>
	vars := mux.Vars(r)
	id := vars["userId"]

	if bson.IsObjectIdHex(id) == false {
		return http.StatusBadRequest, &ApiError{ApiErrorParamIdFmtErr, errors.New("wrong format id")}
	}

	posts, err := appCtx.Db.GetPostsByUserId(bson.ObjectIdHex(id))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&posts)

	return http.StatusOK, nil
}

func CreateCommentForPostIdHandler(appCtx *appcontext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	//<1> get token string
	tokenStr, err := parseJwtTokenStrFromHeaderOrUrlQuery(r.Header, r.URL)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorTokenNotFound, err}
	}

	//<2> jwt token verify, parse to UserEmbed
	var userEmbed *database.UserEmbed
	userEmbed, err = parseJwtTokenStrToUserEmbed(tokenStr)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorTokenParseErr, err}
	}

	//<3> parse form and get post content created by http client
	type commentHttpClientParams struct {
		From    database.UserEmbed
		Message string
	}
	//http://stackoverflow.com/questions/15672556/handling-json-post-request-in-go
	var clientComment commentHttpClientParams
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&clientComment)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, err}
	}

	if clientComment.From.Id.Hex() == "" {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("comment must have a creater user id")}
	}
	if clientComment.From.Id.Hex() != userEmbed.Id.Hex() {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("comment's creater user id must equal to token's user id")}
	}

	if clientComment.Message == "" {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("message of comment can not be nil")}
	}
	//http://stackoverflow.com/questions/12668681/go-language-string-length
	if utf8.RuneCountInString(clientComment.Message) > 1024 {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("message of comment can not longger than 1024 unicode string")}
	}

	clientComment.From.Name = userEmbed.Name

	//<5> insert to db
	vars := mux.Vars(r)
	postId := vars["postId"]

	comment := database.Comment{bson.NewObjectId(),
		clientComment.From,
		clientComment.Message,
		time.Now()}
	err = appCtx.Db.CreateCommentForPostId(&comment, bson.ObjectIdHex(postId))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "application/json")

	jsonStr := fmt.Sprintf("{\"postId\": \"%s\", \"commentId\": \"%s\"}", postId, comment.Id.Hex())
	w.Write([]byte(jsonStr))
	return http.StatusOK, nil
}
