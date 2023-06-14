package async

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type testTask struct {
	number   int
	waitTime time.Duration
}

func (task testTask) Run() {
	fmt.Println("-> Work item: " + task.String())
	time.Sleep(task.waitTime)
	fmt.Println("<- Work item " + task.String())
}

func (task *testTask) String() string {
	return strconv.Itoa(task.number)
}

func TestRatedLimitedWorkerPool(t *testing.T) {
	wp := NewRateLimitedWorkerPool(0, 25) // For CoC, the rate is 40, the minRate is 25ms
	fmt.Println(wp)

	waitTime, _ := time.ParseDuration("1s")
	for i := 0; i < 20; i++ {
		task := testTask{
			number:   i + 1,
			waitTime: waitTime,
		}
		wp.AddTask(task.Run)
	}
	wait, _ := time.ParseDuration("5s")
	time.Sleep(wait)
}
