package protocol

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/WangYihang/Proxy-Verifier/internal/model"
	"github.com/WangYihang/Proxy-Verifier/internal/util"
)

func EstablishHttpsTransparentProxy(proxyHost string, proxyPort uint16, timeoutDuration time.Duration) (net.Conn, error) {
	// Establish tunnel
	conn, err := tls.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", proxyHost, proxyPort),
		&tls.Config{
			InsecureSkipVerify: true,
		},
	)

	if err != nil {
		return nil, err
	}
	return conn, nil
}

func HttpViaHttpsTransparentProxy(task *model.Task) (result model.Result) {
	result = model.Result{Task: task}
	timeoutDuration := time.Second * time.Duration(task.Timeout)

	// Establish tunnel
	conn, err := EstablishHttpsTransparentProxy(task.ProxyHost, task.ProxyPort, timeoutDuration)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// Create an HTTP client that uses the tunnel connection
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return conn, nil
			},
		},
		Timeout: timeoutDuration,
	}

	// Send request
	newUrl := util.BuildUrl(task.TargetUri, task.TaskId, task.ProxyProtocol, task.ProxyHost, task.ProxyPort)
	err = util.SendVerificationHttpRequest(client, newUrl, &result)
	if err != nil {
		result.Error = err.Error()
		return
	}
	return result
}
