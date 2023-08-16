package processer

import (
	"github.com/WangYihang/Proxy-Verifier/internal"
	"github.com/WangYihang/Proxy-Verifier/internal/model"
)

func ScheduleTask(task *model.Task, taskQueue chan *model.Task) {
	taskQueue <- task
	internal.State.TaskScheduled()
}
