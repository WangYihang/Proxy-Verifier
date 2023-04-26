package protocol

import (
	"net/http"
	"net/url"
	"time"

	"github.com/WangYihang/Proxy-Verifier/internal/model"
	"github.com/WangYihang/Proxy-Verifier/internal/util"
)

func HttpViaHttpProxy(task *model.Task) (result model.Result) {
	result = model.Result{Task: task}
	timeoutDuration := time.Second * time.Duration(task.Timeout)

	// Create transport
	proxyURIObject, _ := url.Parse(task.ProxyUri)
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURIObject)}

	// Create client
	client := &http.Client{
		Transport: transport,
		Timeout:   timeoutDuration,
	}

	// Send request
	cacheBypassingUrl := util.BuildUrl(task.TargetUri, task.TaskId, task.ProxyProtocol, task.ProxyHost, task.ProxyPort)
	err := util.SendVerificationHttpRequest(client, cacheBypassingUrl, &result)
	if err != nil {
		result.Error = err.Error()
		return
	}

	return result
}
