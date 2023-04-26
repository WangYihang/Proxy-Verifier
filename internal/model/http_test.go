package model_test

import (
	"bytes"
	"testing"

	"github.com/WangYihang/Proxy-Verifier/internal/model"
)

func TestHttpRequestPath(t *testing.T) {
	type HttpRequestPathToBytesTest struct {
		HttpRequestPath model.HttpRequestPath
		Expected        []byte
	}
	var httpRequestPathToBytesTests = []HttpRequestPathToBytesTest{
		{
			HttpRequestPath: model.HttpRequestPath{
				Path:     "/",
				Params:   map[string]string{},
				Fragment: "",
			},
			Expected: []byte("/"),
		},
		{
			HttpRequestPath: model.HttpRequestPath{
				Path: "/",
				Params: map[string]string{
					"key": "value",
				},
				Fragment: "fragment",
			},
			Expected: []byte("/?key=value#fragment"),
		},
		{
			HttpRequestPath: model.HttpRequestPath{
				Path: "/",
				Params: map[string]string{
					"a":   "b",
					"key": "value",
				},
				Fragment: "fragment",
			},
			Expected: []byte("/?a=b&key=value#fragment"),
		},
		{
			HttpRequestPath: model.HttpRequestPath{
				Path: "/index.php",
				Params: map[string]string{
					"a": "b",
				},
				Fragment: "top",
			},
			Expected: []byte("/index.php?a=b#top"),
		},
	}
	for _, testCase := range httpRequestPathToBytesTests {
		actual := testCase.HttpRequestPath.Bytes()
		if !bytes.Equal(actual, testCase.Expected) {
			t.Errorf("Expected %s, got %s", testCase.Expected, actual)
		}
	}
}

func TestHttpRequestLine(t *testing.T) {
	type HttpRequestLineToBytesTest struct {
		HttpRequestLine model.HttpRequestLine
		Expected        []byte
	}
	var httpRequestLineToBytesTests = []HttpRequestLineToBytesTest{
		{
			HttpRequestLine: model.HttpRequestLine{
				Method: "GET",
				Path: &model.HttpRequestPath{
					Path:     "/",
					Params:   map[string]string{},
					Fragment: "",
				},
				Version: "HTTP/1.1",
			},
			Expected: []byte("GET / HTTP/1.1\r\n"),
		},
	}
	for _, testCase := range httpRequestLineToBytesTests {
		actual := testCase.HttpRequestLine.Bytes()
		if !bytes.Equal(actual, testCase.Expected) {
			t.Errorf("Expected %s, got %s", testCase.Expected, actual)
		}
	}
}

func TestHttp(t *testing.T) {
	type HttpToBytesTest struct {
		Http     model.HttpRequest
		Expected []byte
	}
	var httpToBytesTests = []HttpToBytesTest{
		{
			Http: model.HttpRequest{
				HttpRequestLine: &model.HttpRequestLine{
					Method: "GET",
					Path: &model.HttpRequestPath{
						Path:     "/",
						Params:   map[string]string{},
						Fragment: "",
					},
					Version: "HTTP/1.1",
				},
				Headers: map[string]string{
					"Host": "example.com",
				},
				RawBody: []byte(""),
			},
			Expected: []byte("GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"),
		},
		{
			Http: model.HttpRequest{
				HttpRequestLine: &model.HttpRequestLine{
					Method: "POST",
					Path: &model.HttpRequestPath{
						Path:     "/",
						Params:   map[string]string{},
						Fragment: "",
					},
					Version: "HTTP/1.1",
				},
				Headers: map[string]string{
					"Host": "example.com",
				},
				RawBody: []byte("Hello World!"),
			},
			Expected: []byte("POST / HTTP/1.1\r\nHost: example.com\r\nContent-Length: 12\r\n\r\nHello World!"),
		},
		{
			Http: model.HttpRequest{
				HttpRequestLine: &model.HttpRequestLine{
					Method: "POST",
					Path: &model.HttpRequestPath{
						Path:     "/",
						Params:   map[string]string{},
						Fragment: "",
					},
					Version: "HTTP/1.1",
				},
				Headers: map[string]string{
					"Host": "example.com",
				},
				RawBody: []byte(""),
			},
			Expected: []byte("POST / HTTP/1.1\r\nHost: example.com\r\n\r\n"),
		},
	}
	for _, testCase := range httpToBytesTests {
		actual := testCase.Http.Bytes()
		if !bytes.Equal(actual, testCase.Expected) {
			t.Errorf("Expected %s, got %s", testCase.Expected, actual)
		}
	}
}

func TestNewHttpRequest(t *testing.T) {
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
