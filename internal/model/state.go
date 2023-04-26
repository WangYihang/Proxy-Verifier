package model

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

type State struct {
	MonitorInterval int `json:"monitor_interval"`

	CurrentIndex  int `json:"current_index"`
	PreviousIndex int `json:"previous_index"`

	NumScheduledTasks int `json:"num_scheduled_tasks"`
	NumFinishedTasks  int `json:"num_finished_tasks"`
	NumFailedTasks    int `json:"num_failed_tasks"`
	NumSucceedTasks   int `json:"num_succeed_tasks"`

	lock *sync.Mutex `json:"-"`
}

func NewState(monitorInterval int) *State {
	return &State{
		MonitorInterval: monitorInterval,

		CurrentIndex:  0,
		PreviousIndex: 0,

		NumScheduledTasks: 0,
		NumFinishedTasks:  0,
		NumFailedTasks:    0,
		NumSucceedTasks:   0,

		lock: &sync.Mutex{},
	}
}

func (s *State) String() string {
	var m runtime.MemStats
	var scheduleSpeed, hitRate float64
	scheduleSpeed = float64(s.CurrentIndex-s.PreviousIndex) / float64(s.MonitorInterval)
	hitRate = float64(s.NumSucceedTasks) / float64(s.NumSucceedTasks+s.NumFailedTasks)
	runtime.ReadMemStats(&m)
	return fmt.Sprintf(
		"[%s] current index: %s, schedule speed: %s/s, %s failed, %s succeed, %s / %s overall, hitrate: %s, memory usage: %s, go routine number: %s",
		color.New(color.FgYellow).Sprint(time.Now().Format("2006-01-02 15:04:05")),
		color.New(color.FgGreen).Sprint(s.CurrentIndex),
		color.New(color.FgBlue).Sprint(scheduleSpeed),
		color.New(color.FgRed).Sprint(s.NumFailedTasks),
		color.New(color.FgGreen).Sprint(s.NumSucceedTasks),
		color.New(color.FgYellow).Sprint(s.NumFinishedTasks),
		color.New(color.FgYellow).Sprint(s.NumScheduledTasks),
		color.New(color.FgCyan).Sprintf("%.2f%%", hitRate*100),
		color.New(color.FgYellow).Sprint(humanize.Bytes(m.Alloc)),
		color.New(color.FgYellow).Sprint(runtime.NumGoroutine()),
	)
}

func (s *State) TaskScheduled() {
	s.lock.Lock()
	s.NumScheduledTasks += 1
	s.lock.Unlock()
}

func (s *State) TaskDone(isSucceed bool) {
	s.lock.Lock()
	s.NumFinishedTasks += 1
	if isSucceed {
		s.NumSucceedTasks += 1
	} else {
		s.NumFailedTasks += 1
	}
	s.lock.Unlock()
}

func (s *State) Snapshot() {
	fmt.Println(s.String())
	s.PreviousIndex = s.CurrentIndex
}
