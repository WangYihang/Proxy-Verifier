package model_test

import (
	"bytes"
	"testing"

	"github.com/WangYihang/Proxy-Verifier/internal/model"
)

func TestHttpProxyRequest(t *testing.T) {
	type testcase struct {
		HttpProxyRequest model.HttpProxyRequest
		Expected         []byte
	}
	var testcases = []testcase{
		{
			HttpProxyRequest: model.HttpProxyRequest{
				HttpRequest:    model.HttpRequest{HttpRequestLine: &model.HttpRequestLine{Method: "GET", Path: &model.HttpRequestPath{Path: "/", Params: map[string]string{}, Fragment: ""}, Version: "HTTP/1.1"}, Headers: map[string]string{"Host": "example.com"}, RawBody: []byte("")},
				TargetMethod:   "",
				TargetProtocol: "http",
				TargetHost:     "1.1.1.1",
				TargetPort:     80,
				TargetPath:     "/",
				TargetParams:   map[string]string{},
				TargetFragment: "",
				TargetBody:     []byte{},
			},
			Expected: []byte("GET http://1.1.1.1:80/ HTTP/1.1\r\nHost: example.com\r\n\r\n"),
		},
	}
	for _, testcase := range testcases {
		actual := testcase.HttpProxyRequest.Bytes()
		if !bytes.Equal(actual, testcase.Expected) {
			t.Errorf("Expected %s, got %s", testcase.Expected, actual)
		}
	}
}

func TestNewHttpProxyRequest(t *testing.T) {
	type testcase struct {
		ProxyHost      string
		ProxyPort      uint16
		TargetMethod   string
		Targetmodel    string
		TargetHost     string
		TargetPort     uint16
		TargetPath     string
		TargetParams   map[string]string
		targetFragment string
		TargetHeaders  map[string]string
		TargetBody     []byte
		Expected       []byte
	}
	var testcases = []testcase{
		{
			ProxyHost:      "1.1.1.1",
			ProxyPort:      80,
			TargetMethod:   "POST",
			Targetmodel:    "http",
			TargetHost:     "2.2.2.2",
			TargetPort:     80,
			TargetPath:     "/",
			TargetParams:   map[string]string{},
			targetFragment: "",
			TargetHeaders:  map[string]string{},
			TargetBody:     []byte("Hello World"),
			Expected: []byte("POST http://2.2.2.2:80/ HTTP/1.1\r\n" +
				"Host: 1.1.1.1:80\r\n" +
				"Content-Length: 11\r\n" +
				"\r\n" +
				"Hello World"),
		},
	}
	for _, testcase := range testcases {
		actual, err := model.NewHttpProxyRequest(
			testcase.ProxyHost, testcase.ProxyPort,
			testcase.TargetMethod, testcase.Targetmodel, testcase.TargetHost, testcase.TargetPort, testcase.TargetPath, testcase.TargetParams, testcase.targetFragment,
			testcase.TargetHeaders,
			testcase.TargetBody,
		)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if !bytes.Equal(actual.Bytes(), testcase.Expected) {
			t.Errorf("Expected %s, got %s", testcase.Expected, actual.Bytes())
		}
	}

	type NewHttpRequestTest struct {
		Method   string
		Path     string
		Params   map[string]string
		Headers  map[string]string
		Body     []byte
		Expected []byte
	}
	var newHttpRequestTests = []NewHttpRequestTest{
		{
			Method: "GET",
			Path:   "http://www.example.com:25/index.php",
			Params: map[string]string{
				"id": "1",
			},
			Headers: map[string]string{
				"Custom-Header": "Custom-Value",
			},
			Body:     []byte(""),
			Expected: []byte("GET /index.php?id=1 HTTP/1.1\r\nHost: www.example.com:25\r\nCustom-Header: Custom-Value\r\n\r\n"),
		},
	}
	for _, testCase := range newHttpRequestTests {
		actual, err := model.NewHttpRequest(testCase.Method, testCase.Path, testCase.Params, testCase.Headers, testCase.Body)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if !bytes.Equal(actual.Bytes(), testCase.Expected) {
			t.Errorf("Expected %s, got %s", testCase.Expected, actual.Bytes())
		}
	}
}
