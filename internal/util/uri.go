package util

import (
	"net"
	"net/url"
	"strconv"
)

func ExtractHostPortFromUri(uri string) (host string, port uint16, err error) {
	var portString string
	var portInt int
	var splitHostPortError, lookupErr, parseIntErr error

	// Parse Uri
	url, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", 0, err
	}

	// Extract host and port from TargetUri
	host, portString, splitHostPortError = net.SplitHostPort(url.Host)

	// Copied from go1.20.2/src/net/ipsock.go:SplitHostPort
	const (
		missingPort = "missing port in address"
	)
	if splitHostPortError != nil {
		if splitHostPortError.(*net.AddrError).Err != missingPort {
			return "", 0, splitHostPortError
		} else {
			// Find default port of the current protocol
			portInt, lookupErr = net.LookupPort("tcp", url.Scheme)
			if lookupErr != nil {
				return "", 0, lookupErr
			}
			portString = strconv.Itoa(portInt)
			host = url.Host
		}
	}

	portInt, parseIntErr = strconv.Atoi(portString)
	if parseIntErr != nil {
		return "", 0, parseIntErr
	}
	return host, uint16(portInt), nil
}
