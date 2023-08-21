package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/WangYihang/Proxy-Verifier/internal/config"
	"github.com/WangYihang/Proxy-Verifier/internal/model"
	"github.com/gin-gonic/gin"
	flags "github.com/jessevdk/go-flags"
)

type Response struct {
	BackendAddress string              `json:"backend_address"`
	RequestUri     string              `json:"request_uri"`
	RequestHeaders map[string][]string `json:"request_headers"`
	RequestBody    string              `json:"request_body"`
	Challenge      string              `json:"challenge"`
	Response       string              `json:"response"`
}

var logFd *os.File
var options *model.HTTPEchoServerOptions

func HandleChallenge(challenge string) string {
	var secret = options.Secret
	hash := md5.Sum([]byte(challenge + secret))
	answer := hex.EncodeToString(hash[:])
	return answer
}

func Handler(c *gin.Context) {
	challenge := c.Query(config.ChallengeParamName)
	response := HandleChallenge(challenge)
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Set custom log format
	logLine, err := json.Marshal(map[string]interface{}{
		"remote_addr":      c.Request.RemoteAddr,
		"client_ip":        c.ClientIP(),
		"timestamp":        time.Now().UnixMilli(),
		"request_method":   c.Request.Method,
		"request_path":     c.Request.RequestURI,
		"request_protocol": c.Request.Proto,
		"request_headers":  c.Request.Header,
		"request_body":     string(requestBody),
		"challenge":        challenge,
		"response":         response,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		_, err := logFd.Write([]byte(fmt.Sprintf("%s\n", logLine)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, Response{
		BackendAddress: c.Request.RemoteAddr,
		RequestUri:     c.Request.RequestURI,
		RequestHeaders: c.Request.Header,
		RequestBody:    string(requestBody),
		Challenge:      challenge,
		Response:       response,
	})
}

func init() {
	var err error
	options = &model.HTTPEchoServerOptions{}
	_, err = flags.Parse(options)
	if err != nil {
		os.Exit(1)
	}

	logFd, err = os.OpenFile(options.LogFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	r := gin.Default()
	r.NoRoute(Handler)
	err := r.Run(fmt.Sprintf("%s:%d", options.BindHost, options.BindPort))
	if err != nil {
		log.Fatal(err.Error())
	}
}
