package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/WangYihang/Proxy-Verifier/internal"
	"github.com/WangYihang/Proxy-Verifier/internal/config"
	"github.com/WangYihang/Proxy-Verifier/internal/model"
	"github.com/google/uuid"
)

func BuildUrl(urlString string, challenge, proxyProtocol, proxyHost string, proxyPort uint16) string {
	urlObject, err := url.Parse(urlString)
	if err != nil {
		fmt.Println(err)
	}
	urlObject.Path = fmt.Sprintf("/%s.php", uuid.New().String())
	q := urlObject.Query()
	q.Set(config.ChallengeParamName, challenge)
	q.Set("proxy_protocol", proxyProtocol)
	q.Set("proxy_host", proxyHost)
	q.Set("proxy_port", fmt.Sprintf("%d", proxyPort))
	urlObject.RawQuery = q.Encode()
	return urlObject.String()
}

func ExtractChallengeFromUrl(urlString string) string {
	urlObject, err := url.Parse(urlString)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return urlObject.Query().Get(config.ChallengeParamName)
}

func HandleChallenge(challenge string) string {
	var secret = internal.Options.Secret
	hash := md5.Sum([]byte(challenge + secret))
	answer := hex.EncodeToString(hash[:])
	return answer
}

func VerifyChallenge(challenge string, answer string) bool {
	expectedAnswer := HandleChallenge(challenge)
	return expectedAnswer == answer
}

func VerifyIntegrity(responseBody []byte, challenge string) bool {
	type EchoServerResponse struct {
		BackendAddress string              `json:"backend_address"`
		RequestUri     string              `json:"request_uri"`
		RequestHeaders map[string][]string `json:"request_headers"`
		RequestBody    string              `json:"request_body"`
		Challenge      string              `json:"challenge"`
		Response       string              `json:"response"`
	}
	echoServerResponse := EchoServerResponse{}
	err := json.Unmarshal(responseBody, &echoServerResponse)
	if err != nil {
		return false
	}
	return VerifyChallenge(challenge, echoServerResponse.Response)
}

func SendVerificationHttpRequest(client *http.Client, url string, result *model.Result) error {
	// Create request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Record the raw request in result
	rawRequest, err := httputil.DumpRequest(request, true)
	if err != nil {
		return err
	}
	result.RawProxyOriginRequest = string(rawRequest)

	// Send the request
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	// Record the raw response in result
	rawResponse, err := httputil.DumpResponse(response, true)
	if err != nil {
		return err
	}

	// Record the raw response in result
	result.RawProxyOriginResponse = string(rawResponse)
	rawBody := rawResponse[bytes.Index(rawResponse, []byte("\r\n\r\n"))+4:]

	// Verify integrity
	challenge := ExtractChallengeFromUrl(url)
	if !VerifyIntegrity(rawBody, challenge) {
		result.Error = "integrity verification failed"
		return err
	}

	// No error
	return nil
}
