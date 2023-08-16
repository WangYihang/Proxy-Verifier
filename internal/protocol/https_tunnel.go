package protocol

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/WangYihang/Proxy-Verifier/internal/model"
	"github.com/WangYihang/Proxy-Verifier/internal/util"
)

func EstablishHttpsTunnel(proxyHost string, proxyPort uint16, targetUri string, timeoutDuration time.Duration, result *model.Result) (net.Conn, error) {
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

	host, port, err := util.ExtractHostPortFromUri(targetUri)
	if err != nil {
		result.Error = err.Error()
		return nil, err
	}

	rawClientProxyRequest := fmt.Sprintf("CONNECT %s:%d HTTP/1.1\r\nHost: %s:%d\r\n\r\n", host, port, proxyHost, proxyPort)
	remainingBytes := []byte(rawClientProxyRequest)
	for len(remainingBytes) > 0 {
		_ = conn.SetDeadline(time.Now().Add(timeoutDuration))
		n, err := conn.Write(remainingBytes)
		_ = conn.SetDeadline(time.Time{})
		if err != nil {
			return nil, err
		}
		remainingBytes = remainingBytes[n:]
	}
	result.RawClientProxyRequest = rawClientProxyRequest

	// Read the first line of response
	buffer := make([]byte, 1024)
	_ = conn.SetDeadline(time.Now().Add(timeoutDuration))
	n, err := conn.Read(buffer)
	_ = conn.SetDeadline(time.Time{})
	if err != nil {
		return nil, err
	}

	// Check the first line of response
	firstLine := string(buffer[:n])
	result.RawClientProxyResponse = firstLine

	// HTTP/1.0 200 Connection established
	// HTTP/1.1 200 Connection established
	// HTTP/1.0 200 Connected
	// HTTP/1.1 200 Connected
	// HTTP/1.1 200 OK
	if strings.Contains(result.RawClientProxyResponse, " 200 ") {
		return conn, nil
	} else {
		return nil, fmt.Errorf("invalid Tunnel Server Response: %s", firstLine)
	}
}

func HttpViaHttpsTunnel(task *model.Task) (result model.Result) {
	result = model.Result{Task: task}
	timeoutDuration := time.Second * time.Duration(task.Timeout)

	// Connect to the proxy
	conn, err := EstablishHttpsTunnel(task.ProxyHost, task.ProxyPort, task.TargetUri, timeoutDuration, &result)
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
