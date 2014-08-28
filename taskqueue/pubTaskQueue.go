package taskqueue

import (
	"log"
	"math/rand"
	"time"
)

type PubTask struct {
	Msg string
}

type PubTaskQueue chan *PubTask

func NewPubTaskQueue(workerNum int) PubTaskQueue {
	queue := make(PubTaskQueue, 500)

	for i := 0; i < workerNum; i++ {
		go pubTaskQueueWorker(queue, i)
	}

	return queue
}

func pubTaskQueueWorker(queue PubTaskQueue, index int) {
	for task := range queue {
		random := rand.Intn(3)
		time.Sleep(time.Duration(random) * time.Second)

		log.Printf("PubTaskQueueWorker [%d] do task: %s", index, task.Msg)
	}
}
