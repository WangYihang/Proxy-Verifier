# Proxy Verifier

> Read the project `README.md` in other languages: [English](README.en.md), [简体中文](README.zh.md)

Many organizations on the Internet regularly publish proxy server lists, such as [TheSpeedX/PROXY-List](https://github.com/TheSpeedX/PROXY-List) repository, etc. The authors of these repositories regularly update the latest proxy server lists, and we thank them for their outstanding work. However, **not all of these proxies are available**. For example, some require a username and password, and some proxies will expire in a short time.

This project can help you **quickly obtain public and genuinely available proxy server resources** on the Internet, including HTTP proxies, HTTP tunnels, transparent HTTP proxies, Socks4 proxies, Socks4a proxies, and Socks5 proxies.

This project is divided into 4 parts, or four command-line tools, which are `downloader`, `verifier`, `server`, and `exporter`.

1. `downloader`: Download real-time updated raw proxy files from the repository and merge them;
2. `server`: An HTTP server that helps the `verifier` determine if the proxy server is available. It will accept a random string parameter passed in by the verifier, perform a hash operation on the parameter, and return the result as an HTTP response to ensure that the proxy server is indeed accessed by the `verifier` and that there is no cache server on the link;
3. `verifier`: Verify whether the proxy servers downloaded by the `downloader` are available. Specifically, the verifier will connect to a controlled HTTP server through the proxy, and judge whether the proxy is available by judging whether the server's response is valid;
4. `exporter`: Read the log files of the `verifier` and export the list of available proxy servers from them.

It is easy to see that the `downloader` has no input (the download address list is built into the tool), the output of the `downloader` is the input of the `verifier`, the output of the `verifier` is the input of the `exporter`, and the output of the `exporter` is **the list of available proxy servers** filtered by this project for you.

## Installation

```bash
go install github.com/WangYihang/Proxy-Verifier/cmd/downloader@latest
go install github.com/WangYihang/Proxy-Verifier/cmd/exporter@latest
go install github.com/WangYihang/Proxy-Verifier/cmd/server@latest
go install github.com/WangYihang/Proxy-Verifier/cmd/verifier@latest
```

## Usage

### Downloader

```bash
Usage:
  main [OPTIONS]

Application Options:
  -i, --input-file=  The input file in yaml format (default: -)
  -o, --output-file= The output file (default: -)
  -n, --num-workers= Number of workers (default: 4)
  -m, --max-retries= Maximum number of retries (default: 3)

Help Options:
  -h, --help         Show this help message
```

### Server

```bash
Usage:
  main [OPTIONS]

Application Options:
  -b, --bind-host=    The host to bind (default: 127.0.0.1)
  -p, --bind-port=    The port to bind (default: 80)
  -l, --log-filename= The filename to log to (default: gin.log)
  -s, --secret=       The secret used to verify the integrity of the proxy (default:
                      2d7c29dd-cecb-4454-a4ec-ae2734771a60)

Help Options:
  -h, --help          Show this help message
```

### Verifier

```bash
Usage:
  main [OPTIONS]

Application Options:
  -i, --input-file=       The input file
  -o, --output-file=      The output file
  -u, --url=              The target URL to connect through the proxy, e.g., http://www.google.com,
                          smtp://mails.tsinghua.edu.cn
  -t, --timeout=          Timeout in seconds (default: 16)
  -n, --num-workers=      Number of workers (default: 256)
  -m, --monitor-interval= Interval to output the current running state (in seconds) (default: 1)
  -v, --verbose           Show verbose debug information
  -d, --measurement-id=   The measurement ID used to seperate different measurements in logs
  -s, --secret=           The secret used to verify the integrity of the proxy (default:
                          2d7c29dd-cecb-4454-a4ec-ae2734771a60)

Help Options:
  -h, --help              Show this help message
```

### Exporter

```bash
Usage:
  main [OPTIONS]

Application Options:
  -i, --input-file=        The input file
  -o, --output-file=       The output file
  -r, --require-identical  If provided, the frontend IP and backend IP are required to be identical

Help Options:
  -h, --help               Show this help message
```

## Example

0. Define output and input filenames as environment variables

```bash
TODAY=$(date -u '+%Y-%m-%d')
FREE_OPEN_PROXIES_FILEPATH="free-open-proxies-v${TODAY}.txt"
FREE_OPEN_PROXIES_LOG_FILEPATH="free-open-proxies-v${TODAY}.log"
AVAILABLE_FREE_OPEN_PROXIES_FILEPATH="available-free-open-proxies-v${TODAY}.txt"
```

1. First, download the proxy list using the `downloader`.

```bash
$ # Download the default proxies source lists yaml file
$ wget https://raw.githubusercontent.com/WangYihang/Proxy-Verifier/main/sources.yaml
$ ./downloader --input-file source.yaml --output-file ${FREE_OPEN_PROXIES_FILEPATH}
```

2. Next, start the `server` on a machine with a public IP (e.g., 1.2.3.4).

```bash
$ ./server -b 0.0.0.0 -p 80
```

3. Once again, start the `verifier` to verify the availability of the proxies.

```bash
$ ./verifier \
    --input-file ${FREE_OPEN_PROXIES_FILEPATH} \
    --output-file ${FREE_OPEN_PROXIES_LOG_FILEPATH} \
    --url http://1.2.3.4:80/ \
    --num-workers 1024 \
    --timeout 8
```

4. Finally, export the available proxy servers using the `exporter`.

```bash
$ ./exporter \
    --input-file ${FREE_OPEN_PROXIES_LOG_FILEPATH} \
    --output-file ${AVAILABLE_FREE_OPEN_PROXIES_FILEPATH}
http://20.210.113.32:8123/
http://113.254.50.31:80/
...
http://113.252.10.120:8118/
http://121.132.95.7:80/
```

## Limitations

Due to the proxy data sources of this project coming from other projects, it is not possible to completely obtain all the proxy server lists on the entire internet.

## To-Do

- [ ] Support HTTPS proxies
- [ ] Support transparent HTTPS proxies
- [ ] Support HTTPS tunnels
- [x] Allow the downloader to read proxy list download links through a configuration file

## Acknowledgement

* [clarketm/proxy-list](https://github.com/clarketm/proxy-list)
* [jetkai/proxy-list](https://github.com/jetkai/proxy-list)
* [mertguvencli/http-proxy-list](https://github.com/mertguvencli/http-proxy-list)
* [monosans/proxy-list](https://github.com/monosans/proxy-list)
* [MuRongPIG/Proxy-Master](https://github.com/MuRongPIG/Proxy-Master)
* [proxylist-to/proxy-list](https://github.com/proxylist-to/proxy-list)
* [prxchk/proxy-list](https://github.com/prxchk/proxy-list)
* [roosterkid/openproxylist](https://github.com/roosterkid/openproxylist)
* [TheSpeedX/PROXY-List](https://github.com/TheSpeedX/PROXY-List)
