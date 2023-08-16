package model

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type BaseTask struct {
	Index         int    `json:"index"`
	TaskId        string `json:"task_id"`
	MeasurementId string `json:"measurement_id"`
	StartTime     int64  `json:"start_time"`
	Duration      int64  `json:"duration"`
}

type Task struct {
	BaseTask
	ProxyProtocol string                           `json:"proxy_protocol"`
	ProxyHost     string                           `json:"proxy_host"`
	ProxyPort     uint16                           `json:"proxy_port"`
	ProxyUri      string                           `json:"proxy_uri"`
	Timeout       int                              `json:"timeout"`
	TargetUri     string                           `json:"target_uri"`
	ProcessFunc   func(task *Task) (result Result) `json:"-"`
}

func NewTask(index int, measurement_id string, proxyProtocol string, proxyHost string, proxyPort uint16, targetUri string, timeout int, processFunc func(task *Task) (result Result)) (task Task) {
	task = Task{}
	task.TaskId = uuid.New().String()
	task.MeasurementId = measurement_id
	task.Index = index

	task.ProxyProtocol = proxyProtocol
	task.ProxyHost = proxyHost
	task.ProxyPort = uint16(proxyPort)
	task.TargetUri = targetUri
	task.Timeout = timeout
	task.ProcessFunc = processFunc

	task.setProxyUri()
	return task
}

func (task *Task) setProxyUri() {
	task.ProxyUri = fmt.Sprintf("%s://%s:%d", task.ProxyProtocol, task.ProxyHost, task.ProxyPort)
	if task.Timeout > 0 && strings.HasPrefix(task.ProxyProtocol, "socks") {
		task.ProxyUri = fmt.Sprintf("%s?timeout=%ds", task.ProxyUri, task.Timeout)
	}
}
