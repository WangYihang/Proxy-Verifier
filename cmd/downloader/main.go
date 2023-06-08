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
	flags "github.com/jessevdk/go-flags"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var downloadOptions model.DownloadOptions

type Source struct {
	Type string `yaml:"type"`
	Url  string `yaml:"url"`
}

// downloadProxies downloads proxies from a list of URLs
func downloadProxies(inputFilepath, outputFilepath string) error {
	yamlFile, err := os.Open(inputFilepath)
	if err != nil {
		return err
	}
	defer yamlFile.Close()

	// Parse the YAML file into a slice of Source structs
	sources := []Source{}
	err = yaml.NewDecoder(yamlFile).Decode(&sources)
	if err != nil {
		return err
	}

	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// compile regex firstly to get better performance
	regex := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`)

	for _, source := range sources {
		fmt.Printf("[%s] %s\n", source.Type, source.Url)
		resp, err := http.Get(source.Url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		numProxies := 0
		for scanner.Scan() {
			line := scanner.Text()
			matched := regex.MatchString(line)
			if matched {
				ip, port, err := net.SplitHostPort(line)
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, err = tempFile.WriteString(fmt.Sprintf("%s,%s,%s\n", source.Type, ip, port))
				if err != nil {
					logger.Error(err)
					continue
				}
				numProxies++
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(numProxies, "\t", source.Type, "\t", source.Url)
	}

	err = util.DeduplicateLinesRandomly(tempFile.Name(), outputFilepath)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	_, err := flags.ParseArgs(&downloadOptions, os.Args)
	if err != nil {
		os.Exit(1)
	}
	err = downloadProxies(downloadOptions.InputFile, downloadOptions.OutputFile)
	if err != nil {
		os.Exit(1)
	}
}
