package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"

	"github.com/WangYihang/Proxy-Verifier/internal/model"
	"github.com/WangYihang/Proxy-Verifier/internal/util"
	"github.com/fatih/color"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
)

var downloadOptions model.DownloadOptions
var downloadTaskQueue chan *downloadTask
var writeTaskQueue chan *writeTask
var regex *regexp.Regexp

type proxySource struct {
	Protocol string `yaml:"type"`
	Url      string `yaml:"url"`
}

type downloadTask struct {
	source     *proxySource
	maxRetries int
}

type writeTask struct {
	content string
}

func createDownloadTask(source *proxySource) *downloadTask {
	return &downloadTask{
		source:     source,
		maxRetries: downloadOptions.MaxRetries,
	}
}

func errorLog(err error) string {
	return color.New(color.FgRed).Sprint(err.Error())
}

func httpStatusCodeLog(statusCode int) string {
	if statusCode == 200 {
		return color.New(color.FgGreen).Sprintf("[%d]", statusCode)
	} else {
		return color.New(color.FgRed).Sprintf("[%d]", statusCode)
	}
}

func loader() error {
	// Load the YAML file
	yamlFile, err := os.Open(downloadOptions.InputFile)
	if err != nil {
		fmt.Println(errorLog(err))
		return err
	}
	defer yamlFile.Close()

	// Parse the YAML file into a slice of Source structs
	sources := []proxySource{}
	err = yaml.NewDecoder(yamlFile).Decode(&sources)
	if err != nil {
		fmt.Println(errorLog(err))
		return err
	}

	// Dispatch download tasks
	for _, source := range sources {
		task := createDownloadTask(&source)
		downloadTaskQueue <- task
	}

	// Add stop signals for all workers
	for i := 0; i < downloadOptions.NumWorkers; i++ {
		downloadTaskQueue <- nil
	}
	return nil
}

func worker() {
	for {
		task := <-downloadTaskQueue
		if task == nil {
			break
		}
		for ; task.maxRetries > 0; task.maxRetries-- {
			// Download the proxy list
			resp, err := http.Get(task.source.Url)
			if err != nil {
				fmt.Println(errorLog(err))
				continue
			}
			defer resp.Body.Close()
			// Parse the proxy list
			scanner := bufio.NewScanner(resp.Body)
			numProxies := 0
			for scanner.Scan() {
				line := scanner.Text()
				// Match each proxy
				matched := regex.MatchString(line)
				if matched {
					ip, port, err := net.SplitHostPort(line)
					if err != nil {
						fmt.Println(errorLog(err))
						continue
					}
					writeTaskQueue <- &writeTask{
						content: fmt.Sprintf("%s,%s,%s\n", task.source.Protocol, ip, port),
					}
					numProxies += 1
				}
			}
			// color.New(color.FgGreen).Printf("[%d] %s (%d proxies)\n", , ,
			fmt.Printf(
				"%s %s %s\n",
				httpStatusCodeLog(resp.StatusCode),
				color.New(color.FgYellow).Sprintf("%s", task.source.Url),
				color.New(color.FgBlue).Sprintf("(%d proxy services)", numProxies),
			)
			// Exit the loop if the proxy list is downloaded successfully
			break
		}
	}
	writeTaskQueue <- nil
}

func writer() error {
	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		fmt.Println(errorLog(err))
		return err
	}
	defer tempFile.Close()

	// Write tasks to the temp file
	numFinishedWorkers := 0
	for numFinishedWorkers < downloadOptions.NumWorkers {
		task := <-writeTaskQueue
		if task == nil {
			numFinishedWorkers += 1
			continue
		}
		tempFile.WriteString(task.content)
	}

	// Deduplicate the temp file
	numLines, err := util.DeduplicateLinesRandomly(tempFile.Name(), downloadOptions.OutputFile)
	if err != nil {
		fmt.Println(errorLog(err))
		return err
	}
	color.New(color.FgBlue).Printf("%d unique proxy services found totally\n", numLines)
	return nil
}

func init() {
	// initialize the task queues
	downloadTaskQueue = make(chan *downloadTask)
	writeTaskQueue = make(chan *writeTask)
	// compile regex firstly to get better performance
	regex = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`)
}

func main() {
	_, err := flags.ParseArgs(&downloadOptions, os.Args)
	if err != nil {
		os.Exit(1)
	}

	// Start loader
	go loader()

	// Start workers
	for i := 0; i < downloadOptions.NumWorkers; i++ {
		go worker()
	}

	// Start writer
	writer()
}
