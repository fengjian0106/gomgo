package handler

import (
	"encoding/json"
	"errors"

	"code.google.com/p/go.crypto/bcrypt"

	"log"
	"net/http"

	"github.com/fengjian0106/gomgo/appcontext"
	"github.com/fengjian0106/gomgo/database"
)

type ApiToken struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

func PostSigninHandler(appCtx *appcontext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	type clientUserSigninInfo struct {
		Email    string
		Password string
	}

	//<1> parse json
	decoder := json.NewDecoder(r.Body)
	var userSigninInfo clientUserSigninInfo
	err := decoder.Decode(&userSigninInfo)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("user signin info json parse error")}
	}

	//<2>get user
	user, err := appCtx.Db.GetUserByEmail(userSigninInfo.Email)
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("user not found")}
	}

	//<3>compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userSigninInfo.Password))
	if err != nil {
		return http.StatusBadRequest, &ApiError{ApiErrorParamErr, errors.New("wrong password")}
	}

	//<4> create jwt token
	userEmbed := database.UserEmbed{user.Id, user.Name}
	tokenString, err := createJwtTokenStrWithUserEmbed(&userEmbed)
	//log.Printf("===>> %s", tokenString)

	if err != nil {
		log.Printf("Token Signing error: %v\n", err)
		return http.StatusInternalServerError, &ApiError{ApiErrorTokenParseErr, err}
	}

	w.Header().Set("Content-Type", "application/json")

	tokenJson := ApiToken{user.Id.Hex(), tokenString}
	json.NewEncoder(w).Encode(&tokenJson)

	return http.StatusOK, nil
}
