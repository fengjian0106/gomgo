package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fengjian0106/gomgo/database"

	"github.com/dgrijalva/jwt-go"
)

// location of the files used for signing and verification
const (
	privKeyPath = "keys/app.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

// keys are held in global variables
// i havn't seen a memory corruption/info leakage in go yet
// but maybe it's a better idea, just to store the public key in ram?
// and load the signKey on every signing request? depends on  your usage i guess
var (
	verifyKey, signKey []byte
)

// read the key files before starting http handlers
func init() {
	//log.Println("package handler - token.go - init()")
	var err error

	signKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	verifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}
}

/////////////////////////////////////////////
func createJwtTokenStrWithUserEmbed(userEmbed *database.UserEmbed) (string, error) {
	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	// set our claims
	t.Claims["TokenType"] = "AccessToken"
	jsonStr, _ := json.Marshal(userEmbed)
	t.Claims["Padding"] = string(jsonStr)
	//log.Printf("-----------------%s", string(jsonStr))

	// set the expire time
	// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
	t.Claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	tokenString, err := t.SignedString(signKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func parseJwtTokenStrToUserEmbed(tokenStr string) (*database.UserEmbed, error) {
	// validate the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) ([]byte, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	})

	// branch out into the possible error from signing
	switch err.(type) {

	case nil: // no error

		if !token.Valid { // but may still be invalid
			testErr := errors.New("token invalid")
			return nil, testErr
		}

		// see stdout and watch for the CustomUserEmbed, nicely unmarshalled
	//	log.Printf("Someone accessed resricted area! Token:%+v\n", token)

	case *jwt.ValidationError: // something was wrong during the validation
		vErr := err.(*jwt.ValidationError)

		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			testErr := errors.New("Token Expired, get a new one")
			return nil, testErr
		default:
			log.Printf("ValidationError error: %+v\n", vErr.Errors)
			testErr := errors.New("Error while Parsing Token")
			return nil, testErr
		}

	default: // something else went wrong
		testErr := errors.New("Error while Parsing Token")
		return nil, testErr
	}

	paddingStr, ok := token.Claims["Padding"].(string)
	//	log.Printf(">>>>>>-----------%s", paddingStr)
	if ok == false {
		return nil, errors.New("bad token")
	}
	var userEmbed database.UserEmbed
	if err = json.Unmarshal([]byte(paddingStr), &userEmbed); err != nil {
		return nil, err
	}

	return &userEmbed, nil
}

func parseJwtTokenStrFromHeaderOrUrlQuery(header http.Header, url *url.URL) (string, error) {
	tokenStr := ""

	//<1> try to get jwt token string from http header
	authStr := header.Get("Authorization")
	if authStr != "" {
		tokenStr = strings.TrimPrefix(authStr, "Bearer ")
	}

	//<2> if not found in header, try to find it from url query
	if tokenStr == "" {
		tokens, ok := url.Query()["token"]
		if ok == true && len(tokens) != 0 {
			tokenStr = tokens[0]
		}
	}

	//
	if tokenStr == "" {
		return "", errors.New("token not found")
	}

	return tokenStr, nil
}
