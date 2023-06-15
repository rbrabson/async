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
	"sync"
	"time"
)

// rateLimitedWorkerPool is an all-purpose scheduler to restrict the number and frequency of functions being processed
type rateLimitedWorkerPool struct {
	maxConcurrent   int           // Maximum of concurrent, outstanding requests
	minDuration     time.Duration // Minimum interval between sending requests
	completeChannel chan bool     // Buffered channel to notifiy when a piece of work completes
	workerPool      WorkerPool    // Worker pool for running the tasks
	mutex           sync.Mutex    // Lock for updating rateLimiter and lastRequestTime
	rateLimiter     int           // Count of the number of outstanding requests
	lastRequestTime time.Time     // Last time a request was sent; used in making sure there is a minDuration between requests
}

// NewRateLimitedWorkerPool creates a new bottleneck that limits outstanding requests to a maximum number as well as a minimum amount of time between requests
func NewRateLimitedWorkerPool(maxConcurrent int, minTimeInMs int) WorkerPool {
	if maxConcurrent == 0 {
		maxConcurrent = 1
	}

	b := &rateLimitedWorkerPool{
		maxConcurrent:   maxConcurrent,
		minDuration:     time.Duration(minTimeInMs * 1000000),
		workerPool:      NewWorkerPool(maxConcurrent),
		completeChannel: make(chan bool, maxConcurrent),
	}

	return b
}

// taskElem is a task to be executed by work pool.
type taskElem struct {
	task            func()
	completeChannel chan bool
}

// runTask is the task function that executes the queued task and sends a notification on the completion channel once the task is finished.
func (task taskElem) run() {
	task.task()
	task.completeChannel <- true
}

// AddTask queues a task to be executed by the worker pool.
func (wp *rateLimitedWorkerPool) AddTask(task func()) {
	w := taskElem{
		task:            task,
		completeChannel: wp.completeChannel,
	}

	// Make sure only a single go routine is adding work to the channel at any given time. This is necessary
	// to make sure the rate limit function works properly.
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	if wp.rateLimiter == wp.maxConcurrent {
		<-wp.completeChannel
		wp.rateLimiter--
	}

	waitDuration := time.Since(wp.lastRequestTime)
	if waitDuration < wp.minDuration {
		sleepTime := wp.minDuration - waitDuration
		time.Sleep(sleepTime)
	}

	wp.workerPool.AddTask(w.run)
	wp.lastRequestTime = time.Now()
	wp.rateLimiter++
}

// String prints a string representation of the worker pool
func (b *rateLimitedWorkerPool) String() string {

	var sb strings.Builder
	sb.WriteString("{")
	sb.WriteString("maxConcurrent: " + strconv.Itoa(b.maxConcurrent))
	sb.WriteString(", minTime: " + b.minDuration.String())
	sb.WriteString("}")

	str := sb.String()
	return str
}
