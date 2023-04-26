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
)

var downloadOptions model.DownloadOptions

// downloadProxies downloads proxies from a list of URLs
func downloadProxies(outputFilepath string) error {
	// subscriptions is a map of protocol to URLs
	subscriptions := map[string][]string{
		"http": {
			"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
			"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-http.txt",
			"https://raw.githubusercontent.com/mertguvencli/http-proxy-list/main/proxy-list/data.txt",
			"https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/main/http.txt",
			"https://raw.githubusercontent.com/prxchk/proxy-list/main/http.txt",
			"https://raw.githubusercontent.com/proxylist-to/proxy-list/main/http.txt",
			"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt",
			"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies_anonymous/http.txt",
		},
		"https": {
			"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-https.txt",
			"https://raw.githubusercontent.com/roosterkid/openproxylist/main/HTTPS_RAW.txt",
		},
		"socks4": {
			"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks4.txt",
			"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-socks4.txt",
			"https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/main/socks4.txt",
			"https://raw.githubusercontent.com/prxchk/proxy-list/main/socks4.txt",
			"https://raw.githubusercontent.com/proxylist-to/proxy-list/main/socks4.txt",
			"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks4.txt",
			"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies_anonymous/socks4.txt",
			"https://raw.githubusercontent.com/roosterkid/openproxylist/main/SOCKS4_RAW.txt",
		},
		"socks5": {
			"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks5.txt",
			"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-socks5.txt",
			"https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/main/socks5.txt",
			"https://raw.githubusercontent.com/prxchk/proxy-list/main/socks5.txt",
			"https://raw.githubusercontent.com/proxylist-to/proxy-list/main/socks5.txt",
			"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks5.txt",
			"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies_anonymous/socks5.txt",
			"https://raw.githubusercontent.com/roosterkid/openproxylist/main/SOCKS5_RAW.txt",
		},
	}

	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// compile regex firstly to get better performance
	regex := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`)

	// iterate over subscriptions
	for protocol, urls := range subscriptions {
		// iterate over URLs for each protocol
		for _, url := range urls {
			resp, err := http.Get(url)
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
					_, err = tempFile.WriteString(fmt.Sprintf("%s,%s,%s\n", protocol, ip, port))
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
			fmt.Println(numProxies, "\t", protocol, "\t", url)
		}
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
	err = downloadProxies(downloadOptions.OutputFile)
	if err != nil {
		os.Exit(1)
	}
}
