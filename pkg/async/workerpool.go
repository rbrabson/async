// Package async implements utility routines for running tasks asynchronously.
//
// The async package provides different WorkerPool implementations. A WorkerPool
// allows for a defined number of tasks to be run asynchronously. A rate-limited
// WorkerPool also allows for a defined number of taks to be run asynchronously,
// but also ensures there is a minimum amount of time between scheduling tasks.
package async

import (
	"strconv"
	"strings"
)

// WorkerPool is a contract for Worker Pool implementation
type WorkerPool interface {
	// AddTask adds a task to the worker pool to run asynchronously
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
	sb.WriteString("{")
	sb.WriteString("maxWorker: " + strconv.Itoa(wp.maxWorker))
	sb.WriteString("}")

	str := sb.String()
	return str
}
