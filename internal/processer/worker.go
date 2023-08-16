package processer

import (
	"time"

	"github.com/WangYihang/Proxy-Verifier/internal"
	"github.com/WangYihang/Proxy-Verifier/internal/model"
)

func Worker(taskQueue chan *model.Task, resultQueue chan *model.Result) {
	for {
		// Fetch a task
		task := <-taskQueue
		if task == nil {
			break
		}
		// Update start time
		startTime := time.Now()
		task.StartTime = startTime.UnixMilli()
		// Update max index
		internal.State.CurrentIndex = task.Index
		// Process Task
		result := task.ProcessFunc(task)
		// Update end time
		task.Duration = time.Since(startTime).Milliseconds()
		// Save result
		resultQueue <- &result
	}
	// Signal worker is finished
	resultQueue <- nil
}
