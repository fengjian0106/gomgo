package appcontext

import (
	"log"

	zmq "github.com/pebbe/zmq4"
)

func init() {
	log.Println("start zmqMsgQueue")
	go msgQueue()
}

func msgQueue() {
	var err error

	//  Socket facing clients
	frontend, _ := zmq.NewSocket(zmq.ROUTER)
	defer frontend.Close()
	err = frontend.Bind("inproc://zmq/msg/queue")
	if err != nil {
		log.Fatalln("Binding frontend:", err)
	}

	//  Socket facing services
	backend, _ := zmq.NewSocket(zmq.DEALER)
	defer backend.Close()
	backend.Connect("tcp://localhost:5559")

	//  Start the proxy
	err = zmq.Proxy(frontend, backend, nil)
	log.Fatalln("Proxy interrupted:", err)
}
