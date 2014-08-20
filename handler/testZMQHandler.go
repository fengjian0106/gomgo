package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	zmq "github.com/pebbe/zmq4"

	"github.com/fengjian0106/gomgo/context"
)

//https://github.com/pebbe/zmq4/blob/master/examples/hwclient.go
func TestZMQHandler(context *context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	msg := vars["message"]
	//log.Println("TestZMQHandler  ", msg)

	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	requester.Connect("tcp://localhost:5678")

	//You do your own serialization using protocol buffers, msgpack, JSON, or whatever else your applications need to speak
	requester.Send(msg, 0)

	// Wait for reply:
	reply, err := requester.Recv(0)
	if err != nil {
		log.Println("TestZMQHandler, zmq err:  ", err)
		return http.StatusInternalServerError, err
	}
	//fmt.Println("Received ", reply)

	w.Header().Set("Content-Type", "application/json")

	jsonStr := fmt.Sprintf("{\"retMsg\": \"%s\"}", reply)
	w.Write([]byte(jsonStr))
	return http.StatusOK, nil

}
