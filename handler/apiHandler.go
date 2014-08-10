package handler

import (
	"fmt"
	"net/http"

	"github.com/fengjian0106/gomgo/context"
	"github.com/fengjian0106/gomgo/database"
)

type ApiHandler struct {
	Context *context.Context
	Handler func(*context.Context, http.ResponseWriter, *http.Request) (int, error)
}

//if you want server-end rendering view, maybe you can use ViewHandler

func (ah ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.Handler(ah.Context, w, r)
	if err != nil {
		//log.Printf("HTTP %d: %v", status, err)
		var errJson string
		switch err.(type) {
		case *ApiError:
			apiErr, _ := err.(*ApiError)
			errJson = fmt.Sprintf("{\"code\": %d, \"err\": \"%s\"}", apiErr.Code, apiErr.Err.Error())
		case *database.DbError:
			dbErr, _ := err.(*database.DbError)
			errJson = fmt.Sprintf("{\"code\": %d, \"err\": \"%s\"}", dbErr.Code, dbErr.Err.Error())
		default:
			apiErr := ApiError{ApiErrorUnknow, err}
			errJson = fmt.Sprintf("{\"code\": %d, \"err\": \"%s\"}", apiErr.Code, apiErr.Err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(errJson))
		//http.Error(w, errJson, status)  //this func just return plane textmimetype
	}
}
