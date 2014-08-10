package handler

import "fmt"

type ApiError struct {
	Code int   `json: "code"`
	Err  error `json: "err"`
}

func (err *ApiError) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("ApiErrorCode %d -- %s", err.Code, err.Err.Error())
	}

	return fmt.Sprintf("ApiErrorCode %d -- Unknow error ", err.Code)
}

//error code define  _ __ __
const (
	ApiErrorNotAuth               = 10101
	ApiErrorAuthPwdError          = 10102
	ApiErrorAuthIdentifieNotFound = 10103 //name or email

	ApiErrorParamNeedId     = 10201
	ApiErrorParamIdNotFound = 10202
	ApiErrorParamIdFmtErr   = 10203
	ApiErrorParamErr        = 10203

	ApiErrorTokenSignErr  = 10301
	ApiErrorTokenParseErr = 10302
	ApiErrorTokenNotFound = 10303

	ApiErrorUnknow = 19901
)
