package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	"github.com/fengjian0106/gomgo/context"
	"github.com/fengjian0106/gomgo/database"
)

func CreateUserHandler(context *context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	type userHttpClientParams struct {
		Email    string
		Password string
		Name     string
	}

	//<1> parse json
	decoder := json.NewDecoder(r.Body)
	var clientUser userHttpClientParams
	err := decoder.Decode(&clientUser)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("user info json parse error")}
	}

	//<2>check
	if clientUser.Email == "" || clientUser.Password == "" || clientUser.Name == "" {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("email paassword and name can not be nil")}
	}

	if context.Db.CheckEmailUnique(clientUser.Email) == false {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("email has been registered")}
	}

	//<3>hash password
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(clientUser.Password), bcrypt.DefaultCost)
	newUser := database.User{bson.NewObjectId(), clientUser.Email, string(hashPwd), clientUser.Name}

	err = context.Db.CreateUser(&newUser)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	//password should not return to http client
	newUser.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&newUser)
	return http.StatusOK, nil
}

func GetUserByUserIdHandler(context *context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
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
	/**no need to do like this
	vars := mux.Vars(r)
	id, ok := vars["userId"]
	if ok == false {
		return http.StatusBadRequest, &ApiError{ApiErrorParamNeedId, errors.New("you must pass an id")}
	}
	*/
	vars := mux.Vars(r)
	id := vars["userId"]

	if bson.IsObjectIdHex(id) == false {
		return http.StatusBadRequest, &ApiError{ApiErrorParamIdFmtErr, errors.New("wrong format id")}
	}

	user, err := context.Db.GetUserById(bson.ObjectIdHex(id))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	/** structure.Map() can convert struct to map, but not use the json tag
	// if want to convert and use json tag as map key,
	// please look at https://gist.github.com/tonyhb/5819315  and  http://www.golangtc.com/t/53256166320b523f0a000093
	userMap := structure.Map(user)
	log.Println(userMap)
	log.Println(userMap["Password"])
	delete(userMap, "Password")
	*/
	//password should not return to http client
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&user)
	return http.StatusOK, nil
}

func GetUsersHandler(context *context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
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
	users, err := context.Db.GetUsers()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	/**
	userMaps := make([]map[string]interface{}, len(users))
	for index, user := range users {
		userMap := structure.Map(user)
		delete(userMap, "Password")
		userMaps[index] = userMap
	}
	*/
	//delete password
	for _, user := range users {
		user.Password = ""
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&users)
	return http.StatusOK, nil
}

/** TODO
func ChangeUserPwdHandler(context *context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
}

func RestUserPwdHandler(context *context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
}
*/
