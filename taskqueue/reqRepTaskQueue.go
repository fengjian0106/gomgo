package taskqueue

import (
	"fmt"
	"math/rand"
	"time"
)

type ReqRepTask struct {
	Msg        string
	ResultChan chan string
}

type ReqRepTaskQueue chan *ReqRepTask

func NewReqRepTaskQueue(workerNum int) ReqRepTaskQueue {
	queue := make(ReqRepTaskQueue, 500)

	for i := 0; i < workerNum; i++ {
		go reqRepTaskQueueWorker(queue, i)
	}

	return queue
}

func reqRepTaskQueueWorker(queue ReqRepTaskQueue, index int) {
	for task := range queue {
		//worker can also connect to a remote ZeroMQ server
		//Abstractively, we distribute MESSAGE between goroutine by channel
		//Or, between servers by zeromq(RPC, tcp, udp, http...)
		random := rand.Intn(3)
		time.Sleep(time.Duration(random) * time.Second)

		echo := fmt.Sprintf("ReqRepTaskQueueWorker [%d] echo: %s", index, task.Msg)
		task.ResultChan <- echo

		close(task.ResultChan)
	}
}
