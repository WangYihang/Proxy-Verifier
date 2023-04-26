package protocol

import "github.com/WangYihang/Proxy-Verifier/internal/model"

var TargetProtocl2ProxyProtocol2SubProxyProtocol2Handler = map[string]map[string]map[string]func(task *model.Task) model.Result{
	"http": {
		"http": {
			"proxy":             HttpViaHttpProxy,
			"transparent_proxy": HttpViaHttpTransparentProxy,
			"tunnel":            HttpViaHttpTunnel,
		},
		// "https": {
		// 	"proxy":             HttpViaHttpsProxy,
		// 	"transparent_proxy": HttpViaHttpsTransparentProxy,
		// 	"tunnel":            HttpViaHttpsTunnel,
		// },
		"socks4": {
			"proxy": HttpViaSocksProxy,
		},
		"socks4a": {
			"proxy": HttpViaSocksProxy,
		},
		"socks5": {
			"proxy": HttpViaSocksProxy,
		},
	},
}
