package handler

import (
	"errors"
	"fmt"
	"time"

	"net/http"

	"github.com/fengjian0106/gomgo/appcontext"
	"github.com/fengjian0106/gomgo/taskqueue"
)

func ReqRepTaskHandler(appCtx *appcontext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	msg := r.FormValue("msg")
	if msg == "" {
		return http.StatusBadRequest, &ApiError{ApiErrorZmqMsgNotFound, errors.New("no msg")}
	}

	task := &taskqueue.ReqRepTask{Msg: msg, ResultChan: make(chan string)}
	appCtx.ReqRepTaskQueue <- task

	var reply string
	select {
	case reply = <-task.ResultChan:
		//fmt.Println("Received ", reply)
	case <-time.After(time.Second * 2):
		return http.StatusInternalServerError, &ApiError{ApiErrorTaskQueueRequestTimeout, errors.New("task queue request timeout")}
	}

	w.Header().Set("Content-Type", "application/json")

	jsonStr := fmt.Sprintf("{\"retMsg\": \"%s\"}", reply)
	w.Write([]byte(jsonStr))
	return http.StatusOK, nil
}

func PubTaskHandler(appCtx *appcontext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	msg := r.FormValue("msg")
	if msg == "" {
		return http.StatusBadRequest, &ApiError{ApiErrorZmqMsgNotFound, errors.New("no msg")}
	}

	task := &taskqueue.PubTask{Msg: msg}
	appCtx.PubTaskQueue <- task

	w.Header().Set("Content-Type", "application/json")

	jsonStr := fmt.Sprintf("{\"retMsg\": \"%s\"}", "publish success")
	w.Write([]byte(jsonStr))
	return http.StatusOK, nil

}
