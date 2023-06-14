package async

import (
	"strconv"
	"strings"
)

// WorkerPool is a contract for Worker Pool implementation
type WorkerPool interface {
	AddTask(task func())
}

type workerPool struct {
	maxWorker         int
	queuedTaskChannel chan func()
}

func NewWorkerPool(maxRequests int) WorkerPool {
	wp := &workerPool{
		maxWorker:         maxRequests,
		queuedTaskChannel: make(chan func(), maxRequests),
	}

	for i := 0; i < wp.maxWorker; i++ {
		go func(workerID int) {
			for task := range wp.queuedTaskChannel {
				task()
			}
		}(i + 1)
	}

	return wp
}

func (wp *workerPool) AddTask(task func()) {
	wp.queuedTaskChannel <- task
}

func (wp workerPool) String() string {
	sb := strings.Builder{}
	sb.WriteString("{}")
	sb.WriteString("maxWorker: " + strconv.Itoa(wp.maxWorker))
	sb.WriteString("}")

	str := sb.String()
	return str
}
