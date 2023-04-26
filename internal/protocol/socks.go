package protocol

import (
	"net/http"
	"time"

	"github.com/WangYihang/Proxy-Verifier/internal/model"
	"github.com/WangYihang/Proxy-Verifier/internal/util"
	"h12.io/socks"
)

func HttpViaSocksProxy(task *model.Task) (result model.Result) {
	result = model.Result{Task: task}
	timeoutDuration := time.Second * time.Duration(task.Timeout)

	// Create transport
	transport := &http.Transport{Dial: socks.Dial(task.ProxyUri)}

	// Create client
	client := &http.Client{
		Transport: transport,
		Timeout:   timeoutDuration,
	}

	// Send request
	newUrl := util.BuildUrl(task.TargetUri, task.TaskId, task.ProxyProtocol, task.ProxyHost, task.ProxyPort)
	err := util.SendVerificationHttpRequest(client, newUrl, &result)
	if err != nil {
		result.Error = err.Error()
		return
	}
	return result
}
