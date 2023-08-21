package model

import "net/url"

type ExportOptions struct {
	InputFile        string `short:"i" long:"input-file" description:"The input file" required:"true"`
	OutputFile       string `short:"o" long:"output-file" description:"The output file" required:"true"`
	RequireIdentical bool   `short:"r" long:"require-identical" description:"If provided, the frontend IP and backend IP are required to be identical"`
}

type DownloadOptions struct {
	InputFile  string `short:"i" long:"input-file" description:"The input file in yaml format" required:"true" default:"-"`
	OutputFile string `short:"o" long:"output-file" description:"The output file" required:"true" default:"-"`
	NumWorkers int    `short:"n" long:"num-workers" description:"Number of workers" default:"4"`
	MaxRetries int    `short:"m" long:"max-retries" description:"Maximum number of retries" default:"3"`
}

type SecretOptions struct {
	Secret string `short:"s" long:"secret" description:"The secret used to verify the integrity of the proxy" default:"2d7c29dd-cecb-4454-a4ec-ae2734771a60"`
}

type MainOptions struct {
	InputFilepath   string `short:"i" long:"input-file" description:"The input file" required:"true"`
	OutputFile      string `short:"o" long:"output-file" description:"The output file" required:"true"`
	TargetUrl       string `short:"u" long:"url" description:"The target URL to connect through the proxy, e.g., http://www.google.com, smtp://mails.tsinghua.edu.cn" required:"true"`
	Timeout         int    `short:"t" long:"timeout" description:"Timeout in seconds" default:"16"`
	NumWorkers      int    `short:"n" long:"num-workers" description:"Number of workers" default:"256"`
	MonitorInterval int    `short:"m" long:"monitor-interval" description:"Interval to output the current running state (in seconds)" default:"1"`
	Verbose         []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	MeasurementId   string `short:"d" long:"measurement-id" description:"The measurement ID used to seperate different measurements in logs" default:""`
	TargetUrlObject *url.URL
	SecretOptions
}

type HTTPEchoServerOptions struct {
	BindHost    string `short:"b" long:"bind-host" description:"The host to bind" required:"true" default:"127.0.0.1"`
	BindPort    int    `short:"p" long:"bind-port" description:"The port to bind" required:"true" default:"80"`
	LogFilename string `short:"l" long:"log-filename" description:"The filename to log to" required:"true" default:"gin.log"`
	SecretOptions
}

type SMTPEchoServerOptions struct {
	BindHost    string `short:"b" long:"bind-host" description:"The host to bind" required:"true" default:"127.0.0.1"`
	BindPort    int    `short:"p" long:"bind-port" description:"The port to bind" required:"true" default:"25"`
	LogFilename string `short:"l" long:"log-filename" description:"The filename to log to" required:"true" default:"smtp.log"`
	SecretOptions
}
