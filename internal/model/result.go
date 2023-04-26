package model

import (
	"fmt"
)

type Result struct {
	Task  *Task  `json:"task"`
	Error string `json:"error"`
	// Client -> Proxy
	RawClientProxyRequest  string `json:"raw_client_proxy_request"`
	RawClientProxyResponse string `json:"raw_client_proxy_response"`
	// Proxy -> Server
	RawProxyOriginRequest  string `json:"raw_proxy_origin_request"`
	RawProxyOriginResponse string `json:"raw_proxy_origin_response"`
}

func (r *Result) String() string {
	if r.Error == "" {
		return fmt.Sprintf("[%s -> %d] %s://%s:%d (succeed)", r.Task.TaskId, r.Task.Index, r.Task.ProxyProtocol, r.Task.ProxyHost, r.Task.ProxyPort)
	} else {
		return fmt.Sprintf("[%s -> %d] %s://%s:%d (failed)", r.Task.TaskId, r.Task.Index, r.Task.ProxyProtocol, r.Task.ProxyHost, r.Task.ProxyPort)
	}
}
