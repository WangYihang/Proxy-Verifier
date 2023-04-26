package model

import (
	"bytes"
	"fmt"
)

type HttpProxyRequest struct {
	HttpRequest
	TargetMethod   string
	TargetProtocol string
	TargetHost     string
	TargetPort     uint16
	TargetPath     string
	TargetParams   map[string]string
	TargetFragment string
	TargetBody     []byte
}

func (httpProxyRequest *HttpProxyRequest) Bytes() []byte {
	var buffer bytes.Buffer
	httpProxyRequest.HttpRequest.Construct()
	// Build HTTP request line
	httpProxyRequest.HttpRequestLine.Path.Scheme = httpProxyRequest.TargetProtocol
	httpProxyRequest.HttpRequestLine.Path.Host = httpProxyRequest.TargetHost
	httpProxyRequest.HttpRequestLine.Path.Port = httpProxyRequest.TargetPort
	httpProxyRequest.HttpRequestLine.Path.Path = httpProxyRequest.TargetPath
	httpProxyRequest.HttpRequestLine.Path.Params = httpProxyRequest.TargetParams
	httpProxyRequest.HttpRequestLine.Path.Fragment = httpProxyRequest.TargetFragment
	buffer.WriteString(string(httpProxyRequest.HttpRequestLine.Bytes()))
	// Build HTTP headers
	for key, value := range httpProxyRequest.Headers {
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	buffer.WriteString("\r\n")
	// Build HTTP body
	buffer.Write(httpProxyRequest.TargetBody)
	return buffer.Bytes()
}

func NewHttpProxyRequest(
	proxyHost string, proxyPort uint16,
	targetMethod string, targetProtocol string, targetHost string, targetPort uint16, targetPath string, targetParams map[string]string, targetFragment string,
	targetHeaders map[string]string,
	targetBody []byte,
) (*HttpProxyRequest, error) {
	httpRequest, err := NewHttpRequest(
		targetMethod,
		fmt.Sprintf("http://%s:%d/", proxyHost, proxyPort),
		map[string]string{},
		targetHeaders,
		targetBody,
	)
	if err != nil {
		return nil, err
	}

	HttpProxyRequest := &HttpProxyRequest{
		HttpRequest:    *httpRequest,
		TargetMethod:   targetMethod,
		TargetProtocol: targetProtocol,
		TargetHost:     targetHost,
		TargetPort:     targetPort,
		TargetPath:     targetPath,
		TargetParams:   targetParams,
		TargetFragment: targetFragment,
		TargetBody:     targetBody,
	}
	return HttpProxyRequest, nil
}
