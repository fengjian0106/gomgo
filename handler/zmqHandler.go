package handler

import (
	"errors"
	"fmt"
	"log"
	"time"

	"net/http"

	zmq "github.com/pebbe/zmq4"

	"github.com/fengjian0106/gomgo/appcontext"
)

const (
	REQUEST_TIMEOUT = 2000 * time.Millisecond //  msecs, (> 1000!)
	REQUEST_RETRIES = 3                       //  Before we abandon
	SERVER_ENDPOINT = "inproc://zmq/msg/queue"
)

//https://github.com/pebbe/zmq4/blob/master/examples/hwclient.go
func ZMQHandler(appCtx *appcontext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	msg := r.FormValue("msg")
	if msg == "" {
		return http.StatusBadRequest, &ApiError{ApiErrorZmqMsgNotFound, errors.New("no msg")}
	}
	//log.Println("TestZMQHandler  ", msg)

	requester, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		return http.StatusInternalServerError, &ApiError{ApiErrorZmqCanNotCreateSocket, err}
	}
	defer requester.Close()
	requester.Connect(SERVER_ENDPOINT)

	//You do your own serialization using protocol buffers, msgpack, JSON, or whatever else your applications need to speak
	requester.Send(msg, 0)

	// Wait for reply:
	reply, err := requester.Recv(0)
	if err != nil {
		log.Println("TestZMQHandler, zmq err:  ", err)
		return http.StatusInternalServerError, &ApiError{ApiErrorZmqErr, err}
	}
	//fmt.Println("Received ", reply)

	w.Header().Set("Content-Type", "application/json")

	jsonStr := fmt.Sprintf("{\"retMsg\": \"%s\"}", reply)
	w.Write([]byte(jsonStr))
	return http.StatusOK, nil
}

func ZMQWithTimeoutHandler(appCtx *appcontext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	msg := r.FormValue("msg")
	if msg == "" {
		return http.StatusBadRequest, &ApiError{ApiErrorZmqMsgNotFound, errors.New("no msg")}
	}

	//////////////////////////////
	client, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		return http.StatusInternalServerError, &ApiError{ApiErrorZmqCanNotCreateSocket, err}
	}
	client.Connect(SERVER_ENDPOINT)

	poller := zmq.NewPoller()
	poller.Add(client, zmq.POLLIN)

	var reply string
	err = nil
	sequence := 0
	retries_left := REQUEST_RETRIES
	for retries_left > 0 {
		//log.Println("@@ ", sequence)

		//You do your own serialization using protocol buffers, msgpack, JSON, or whatever else your applications need to speak
		//client.Send(msg, 0)
		client.Send(fmt.Sprintf("%s [%d]", msg, sequence), 0)

		for expect_reply := true; expect_reply; {
			//  Poll socket for a reply, with timeout
			sockets, err := poller.Poll(REQUEST_TIMEOUT)
			if err != nil {
				break //  Interrupted
			}

			//  Here we process a server reply and exit our loop if the
			//  reply is valid. If we didn't a reply we close the client
			//  socket and resend the request. We try a number of times
			//  before finally abandoning:
			if len(sockets) > 0 {
				reply, err = client.Recv(0)
				if err != nil {
					break //  Interrupted
				}

				//log.Printf("E: malformed reply from server: %s\n", reply)
				retries_left = 0
				expect_reply = false

			} else {
				retries_left--
				if retries_left == 0 {
					//log.Println("E: server seems to be offline, abandoning")
					break
				} else {
					//log.Println("W: no response from server, retrying...")
					//  Old socket is confused; close it and open a new one
					client.Close()
					client, _ = zmq.NewSocket(zmq.REQ)
					client.Connect(SERVER_ENDPOINT)
					// Recreate poller for new client
					poller = zmq.NewPoller()
					poller.Add(client, zmq.POLLIN)

					//  Send request again, on new socket
					sequence++
					//log.Println("@@ ", sequence)
					//client.Send(msg, 0)
					client.Send(fmt.Sprintf("%s [%d]", msg, sequence), 0)
				}
			}
		}
	}
	client.Close()

	if err != nil {
		return http.StatusInternalServerError, &ApiError{ApiErrorZmqErr, err}
	}
	if err == nil && reply == "" {
		return http.StatusInternalServerError, &ApiError{ApiErrorZmqRequestTimeout, errors.New("zmq request timeout")}
	}

	//
	w.Header().Set("Content-Type", "application/json")

	jsonStr := fmt.Sprintf("{\"retMsg\": \"%s\"}", reply)
	w.Write([]byte(jsonStr))
	return http.StatusOK, nil
}
