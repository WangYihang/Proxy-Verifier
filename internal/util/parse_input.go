package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"
)

func ParseRecordCsv(record []string) (string, string, uint16) {
	var err error
	var proxyProtocol, proxyHost string
	var proxyPort int

	proxyProtocol = record[0]
	proxyHost = record[1]
	proxyPort, err = strconv.Atoi(record[2])
	if err != nil {
		logger.Error(err)
		protocolName2DefaultPort := map[string]int{
			"ftp":   21,
			"ssh":   22,
			"smtp":  25,
			"http":  80,
			"https": 443,
		}
		if defaultPort, ok := protocolName2DefaultPort[proxyProtocol]; ok {
			proxyPort = defaultPort
		} else {
			proxyPort = 8080
		}
	}
	return proxyProtocol, proxyHost, uint16(proxyPort)
}

func ParseRecordXmap(record []string) (string, string, uint16) {
	var proxyProtocol, proxyHost string
	var proxyPort int

	proxyProtocol = "http"
	proxyHost = record[0]
	proxyPort = 8080
	return proxyProtocol, proxyHost, uint16(proxyPort)
}

func DetectFileFormat(inputFilepath string) (string, error) {
	file, err := os.OpenFile(inputFilepath, os.O_RDONLY, 0)
	if err != nil {
		logger.Error(err)
	}
	defer file.Close()

	// read first line
	scanner := bufio.NewScanner(file)
	ok := scanner.Scan()
	if ok {
		line := scanner.Text()
		items := strings.Split(line, ",")
		switch len(items) {
		case 1:
			return "xmap", nil
		case 3:
			return "csv", nil
		default:
			return "", fmt.Errorf("unknown file format: %s", line)
		}
	} else {
		return "", fmt.Errorf("failed to read first line")
	}
}
