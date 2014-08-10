package database

import "fmt"

type DbError struct {
	Code int
	Err  error
}

func (err *DbError) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("DbErrorCode %d -- %s", err.Code, err.Err.Error())
	}

	return fmt.Sprintf("DbErrorCode %d -- Unknow error ", err.Code)
}

const (
	DbErrorModalLackRequiredPrpty = 20101
	DbErrorModalPrptyErr          = 20102

	DbErrorInsertErr = 20201
	DbErrorDeleteErr = 20202
	DbErrorFindErr   = 20203
	DbErrorUpdateErr = 20204

	DbErrorUnknow = 29901
)
