package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/WangYihang/Proxy-Verifier/internal/model"
	"github.com/WangYihang/Proxy-Verifier/internal/util"
	flags "github.com/jessevdk/go-flags"
	logger "github.com/sirupsen/logrus"
)

var exportOptions model.ExportOptions

func transformLine(oldline string) (newline string, err error) {
	var result model.Result
	// Unmarshal line in JSON format
	err = json.Unmarshal([]byte(oldline), &result)
	if err != nil {
		return "", err
	}

	// Extract information from the server header
	protocol := result.Task.ProxyProtocol
	proxyHost := result.Task.ProxyHost
	proxyPort := result.Task.ProxyPort

	if result.Error != "" {
		// Check if the proxy and backend are the same
		return "", fmt.Errorf(result.Error)
	}

	return fmt.Sprintf("%s://%s:%d/\n", protocol, proxyHost, proxyPort), nil
}

func export(inputFilepath string, outputFilepath string) error {
	// Open input file
	inputFile, err := os.Open(inputFilepath)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer inputFile.Close()

	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tempFile.Close()

	// Read input file line by line
	scanner := bufio.NewScanner(inputFile)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024*1024)

	for scanner.Scan() {
		oldline := scanner.Text()
		// Transform json to proxy format
		newline, err := transformLine(oldline)
		if err != nil {
			// fmt.Printf("%s\r", err)
			continue
		}
		fmt.Println(strings.TrimSpace(newline))
		// Write to output file
		_, err = tempFile.WriteString(newline)
		if err != nil {
			// fmt.Printf("%s\r", err)
			continue
		}
	}

	err = util.DeduplicateLinesRandomly(tempFile.Name(), outputFilepath)
	if err != nil {
		return err
	}
	return nil
}

// Export available proxies from the log file to the output file
func main() {
	_, err := flags.ParseArgs(&exportOptions, os.Args)
	if err != nil {
		os.Exit(1)
	}
	err = export(exportOptions.InputFile, exportOptions.OutputFile)
	if err != nil {
		os.Exit(1)
	}
}
