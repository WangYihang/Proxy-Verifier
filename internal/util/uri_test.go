package util_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/WangYihang/Proxy-Verifier/internal/util"
)

func TestExtractHostPort(t *testing.T) {
	testCases := []struct {
		uri          string
		expectedHost string
		expectedPort uint16
		expectedErr  error
	}{
		{
			uri:          "http://localhost/",
			expectedHost: "localhost",
			expectedPort: 80,
			expectedErr:  nil,
		},
		{
			uri:          "http://localhost:8080",
			expectedHost: "localhost",
			expectedPort: 8080,
			expectedErr:  nil,
		},
		{
			uri:          "https://example.com",
			expectedHost: "example.com",
			expectedPort: 443,
			expectedErr:  nil,
		},
		{
			uri:          "ftp://ftp.example.com:21",
			expectedHost: "ftp.example.com",
			expectedPort: 21,
			expectedErr:  nil,
		},
		{
			uri:          "invalid uri",
			expectedHost: "",
			expectedPort: 0,
			expectedErr:  errors.New(`parse "invalid uri": invalid URI for request`),
		},
	}

	for _, tc := range testCases {
		host, port, err := util.ExtractHostPortFromUri(tc.uri)
		fmt.Println(tc.uri, host, port, err)
		if host != tc.expectedHost {
			t.Errorf("Expected host %s but got %s", tc.expectedHost, host)
		}
		if port != tc.expectedPort {
			t.Errorf("Expected port %d but got %d", tc.expectedPort, port)
		}
		if err != nil && tc.expectedErr == nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if err == nil && tc.expectedErr != nil {
			t.Errorf("Expected error %v but got no error", tc.expectedErr)
		}
		if err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error() {
			t.Errorf("Expected error %v but got %v", tc.expectedErr, err)
		}
	}
}
