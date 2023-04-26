package model

import (
	"bytes"
	"fmt"
	neturl "net/url"
	"strings"
)

type HttpRequestPath struct {
	Scheme   string
	Host     string
	Port     uint16
	Path     string
	Params   map[string]string
	Fragment string
}

func (path *HttpRequestPath) Bytes() []byte {
	var buffer bytes.Buffer
	// Write scheme
	if path.Scheme != "" {
		buffer.WriteString(path.Scheme)
		buffer.WriteString("://")
	}
	// Write host
	if path.Host != "" {
		buffer.WriteString(path.Host)
	}
	// Write port
	if path.Port != 0 {
		buffer.WriteString(fmt.Sprintf(":%d", path.Port))
	}
	// Write path
	buffer.WriteString(path.Path)
	// Write params
	var paramStrings []string
	if len(path.Params) != 0 {
		buffer.WriteString("?")
	}
	for key, value := range path.Params {
		paramStrings = append(paramStrings, fmt.Sprintf("%s=%s", key, value))
	}
	buffer.WriteString(strings.Join(paramStrings, "&"))
	// Write fragment
	if path.Fragment != "" {
		buffer.WriteString("#")
		buffer.WriteString(path.Fragment)
	}
	return buffer.Bytes()
}

type HttpRequestLine struct {
	Method  string
	Path    *HttpRequestPath
	Version string
}

func (requestLine *HttpRequestLine) Bytes() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s %s %s\r\n", requestLine.Method, requestLine.Path.Bytes(), requestLine.Version))
	return buffer.Bytes()
}

type HttpRequest struct {
	HttpRequestLine *HttpRequestLine
	Headers         map[string]string
	RawBody         []byte
}

func (httpRequest *HttpRequest) Construct() {
	if len(httpRequest.RawBody) != 0 {
		httpRequest.Headers["Content-Length"] = fmt.Sprintf("%d", len(httpRequest.RawBody))
	}
}

func (httpRequest *HttpRequest) Bytes() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(string(httpRequest.HttpRequestLine.Bytes()))
	buffer.WriteString(fmt.Sprintf("Host: %s\r\n", httpRequest.Headers["Host"]))
	delete(httpRequest.Headers, "Host")
	httpRequest.Construct()
	for key, value := range httpRequest.Headers {
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	buffer.WriteString("\r\n")
	buffer.Write(httpRequest.RawBody)
	return buffer.Bytes()
}

func NewHttpRequest(method string, url string, params map[string]string, headers map[string]string, body []byte) (*HttpRequest, error) {
	urlObj, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	headers["Host"] = urlObj.Host
	return &HttpRequest{
		HttpRequestLine: &HttpRequestLine{
			Method: method,
			Path: &HttpRequestPath{
				Path:     urlObj.Path,
				Params:   params,
				Fragment: urlObj.Fragment,
			},
			Version: "HTTP/1.1",
		},
		Headers: headers,
		RawBody: body,
	}, nil
}
