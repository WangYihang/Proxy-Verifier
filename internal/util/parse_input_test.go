package util_test

import (
	"fmt"
	"testing"

	"github.com/WangYihang/Proxy-Verifier/internal/util"
)

func TestParseRecordCsv(t *testing.T) {
	testcases := []struct {
		record        []string
		expectedProto string
		expectedHost  string
		expectedPort  uint16
	}{
		{
			record:        []string{"http", "127.0.0.1", "8080"},
			expectedProto: "http",
			expectedHost:  "127.0.0.1",
			expectedPort:  8080,
		},
		{
			record:        []string{"http", "192.168.1.1", ""},
			expectedProto: "http",
			expectedHost:  "192.168.1.1",
			expectedPort:  80,
		},
		{
			record:        []string{"https", "www.baidu.com", "8080"},
			expectedProto: "https",
			expectedHost:  "www.baidu.com",
			expectedPort:  8080,
		},
		{
			record:        []string{"https", "127.0.0.1", ""},
			expectedProto: "https",
			expectedHost:  "127.0.0.1",
			expectedPort:  443,
		},
		{
			record:        []string{"ssh", "www.baidu.com", "8080"},
			expectedProto: "ssh",
			expectedHost:  "www.baidu.com",
			expectedPort:  8080,
		},
		{
			record:        []string{"ssh", "127.0.0.1", ""},
			expectedProto: "ssh",
			expectedHost:  "127.0.0.1",
			expectedPort:  22,
		},
		{
			record:        []string{"smtp", "www.baidu.com", "8080"},
			expectedProto: "smtp",
			expectedHost:  "www.baidu.com",
			expectedPort:  8080,
		},
		{
			record:        []string{"smtp", "127.0.0.1", ""},
			expectedProto: "smtp",
			expectedHost:  "127.0.0.1",
			expectedPort:  25,
		},
	}
	for _, testcase := range testcases {
		proxyProtocol, proxyHost, proxyPort := util.ParseRecordCsv(testcase.record)
		fmt.Println(testcase.record, proxyProtocol, proxyHost, proxyPort)
		if proxyProtocol != testcase.expectedProto {
			t.Errorf("proxyProtocol should be %s, got %s", testcase.expectedProto, proxyProtocol)
		}
		if proxyHost != testcase.expectedHost {
			t.Errorf("proxyHost should be %s, got %s", testcase.expectedHost, proxyHost)
		}
		if proxyPort != testcase.expectedPort {
			t.Errorf("proxyPort should be %d, got %d", testcase.expectedPort, proxyPort)
		}
	}
}
