# 代理验证器

互联网上有许多组织会定期发布代理服务器列表，例如 [TheSpeedX/PROXY-List](https://github.com/TheSpeedX/PROXY-List) 仓库等，这些仓库的作者会定期更新最新的代理服务器列表，感谢他们杰出的工作。然而，这些代理**并不是每一个都可用**。例如：有一些需要提供用户名密码，有一些代理会在短时间内失效。

本项目可以帮助您**快速获得互联网上公开且真实可用的代理服务器资源**，包括 HTTP 代理、HTTP 隧道、透明 HTTP 代理、Socks4 代理、Socks4a 代理以及 Socks5 代理等。

本项目分为 4 个部分，或者说，四个命令行工具，分别是下载器（`downloader`）、验证器（`verifier`）、服务器（`server`）与导出器（`exporter`）。

1. 下载器（`downloader`）：从仓库中下载实时更新的原始代理文件，并将他们进行合并；
2. 服务器（`server`）：一个 HTTP 服务器，用于协助验证器判断代理服务器是否可用。其会接受由验证器传入的随机字符串参数，并对该参数进行哈希运算，将结果作为 HTTP 响应返回，以确保代理服务器确实被验证器访问到了，并且链路上没有经过缓存服务器；
3. 验证器（`verifier`）：验证下载器下载的代理服务器是否可用。具体来说，验证器将会通过代理连接受控的 HTTP 服务器，通过判断服务器的响应是否合法来判断代理是否可用；
4. 导出器（`exporter`）：读取验证器的日志文件，从中导出可用的代理服务器列表。

很容易发现，下载器没有输入（下载地址列表内置在该工具中），下载器的输出是验证器的输入，验证器的输出是导出器的输入，导出器的输出即该项目为您过滤出的**可用代理服务器列表**。

## 安装

```bash
go install github.com/WangYihang/Proxy-Verifier/cmd/downloader@latest
go install github.com/WangYihang/Proxy-Verifier/cmd/exporter@latest
go install github.com/WangYihang/Proxy-Verifier/cmd/server@latest
go install github.com/WangYihang/Proxy-Verifier/cmd/verifier@latest
```

## 使用方式

### 下载器

```bash
Usage:
  main [OPTIONS]

Application Options:
  -o, --output-file= The output file

Help Options:
  -h, --help         Show this help message
```

### 服务器

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

### 验证器

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

### 导出器

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

## 示例

0. 定义输出输入文件名作为环境变量

```bash
TODAY=$(date -u '+%Y-%m-%d')
FREE_OPEN_PROXIES_FILEPATH="free-open-proxies-v${TODAY}.txt"
FREE_OPEN_PROXIES_LOG_FILEPATH="free-open-proxies-v${TODAY}.log"
AVAILABLE_FREE_OPEN_PROXIES_FILEPATH="available-free-open-proxies-v${TODAY}.txt"
```

1. 首先通过下载器（`downloader`）下载代理列表

```bash
$ ./downloader --output-file ${FREE_OPEN_PROXIES_FILEPATH}
```

2. 其次，在具有公网 IP 的服务器（如：1.2.3.4）上启动服务器（`server`）

```bash
$ ./server -b 0.0.0.0 -p 80
```

3. 再次，启动验证器（`verifier`）进行代理的可用性验证

```bash
$ ./verifier \
    --input-file ${FREE_OPEN_PROXIES_FILEPATH} \
    --output-file ${FREE_OPEN_PROXIES_LOG_FILEPATH} \
    --url http://1.2.3.4:80/ \
    --num-workers 1024 \
    --timeout 8
```

4. 最后通过导出器（`exporter`）导出可用的代理服务器。

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

## 局限性

由于本项目的代理数据源来自其他项目，因此并不能完整地获得整个互联网上的所有代理服务器列表。

## 待办

- [x] 支持 HTTPS 代理
- [x] 支持透明 HTTPS 代理
- [x] 支持 HTTPS 隧道
- [x] 让下载器通过配置文件读取代理列表下载链接

## 致谢

* [clarketm/proxy-list](https://github.com/clarketm/proxy-list)
* [jetkai/proxy-list](https://github.com/jetkai/proxy-list)
* [mertguvencli/http-proxy-list](https://github.com/mertguvencli/http-proxy-list)
* [monosans/proxy-list](https://github.com/monosans/proxy-list)
* [MuRongPIG/Proxy-Master](https://github.com/MuRongPIG/Proxy-Master)
* [proxylist-to/proxy-list](https://github.com/proxylist-to/proxy-list)
* [prxchk/proxy-list](https://github.com/prxchk/proxy-list)
* [roosterkid/openproxylist](https://github.com/roosterkid/openproxylist)
* [TheSpeedX/PROXY-List](https://github.com/TheSpeedX/PROXY-List)
