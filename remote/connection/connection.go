package connection

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Kcrong/selenium-go"
	"github.com/Kcrong/selenium-go/remote/command"
)

const (
	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 60 * time.Second
)

// ClientConfig represents the configuration for the remote connection
type ClientConfig struct {
	// RemoteServerAddr is the address of the remote server
	RemoteServerAddr string
	// KeepAlive indicates whether to keep the connection alive
	KeepAlive bool
	// IgnoreCertificates indicates whether to ignore SSL certificates
	IgnoreCertificates bool
	// Timeout is the timeout for HTTP requests
	Timeout time.Duration
	// CACerts is the path to the CA certificates file
	CACerts string
	// ExtraHeaders are additional headers to include in requests
	ExtraHeaders map[string]string
	// UserAgent is the user agent string to use
	UserAgent string
	// ProxyURL is the URL of the proxy server to use
	ProxyURL string
}

// NewClientConfig creates a new ClientConfig with default values
func NewClientConfig(remoteServerAddr string) *ClientConfig {
	system := runtime.GOOS
	if system == "darwin" {
		system = "mac"
	}

	return &ClientConfig{
		RemoteServerAddr:   remoteServerAddr,
		KeepAlive:          true,
		IgnoreCertificates: false,
		Timeout:            DefaultTimeout,
		CACerts:            os.Getenv("REQUESTS_CA_BUNDLE"),
		UserAgent:          fmt.Sprintf("selenium/%s (golang %s)", selenium.Version, system),
	}
}

// RemoteConnection represents a connection to a remote WebDriver server
type RemoteConnection struct {
	client       *http.Client
	config       *ClientConfig
	commandMap   map[command.Command]command.CommandInfo
	proxyURL     *url.URL
	proxyAuth    string
	extraHeaders map[string]string
}

// New creates a new RemoteConnection
func New(config *ClientConfig) (*RemoteConnection, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	rc := &RemoteConnection{
		config:       config,
		commandMap:   command.CommandInfoMap,
		extraHeaders: make(map[string]string),
	}

	if config.ExtraHeaders != nil {
		for k, v := range config.ExtraHeaders {
			rc.extraHeaders[k] = v
		}
	}

	if config.ProxyURL != "" {
		proxyURL, err := url.Parse(config.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %v", err)
		}
		rc.proxyURL = proxyURL

		if proxyURL.User != nil {
			password, _ := proxyURL.User.Password()
			rc.proxyAuth = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", proxyURL.User.Username(), password)))
		}
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.IgnoreCertificates,
		},
	}

	if rc.proxyURL != nil {
		transport.Proxy = http.ProxyURL(rc.proxyURL)
	}

	rc.client = &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	return rc, nil
}

// AddCommand adds a new command to the command map
func (rc *RemoteConnection) AddCommand(cmd command.Command, method command.HTTPMethod, path string) {
	rc.commandMap[cmd] = command.CommandInfo{
		Method: method,
		Path:   path,
	}
}

// GetCommand returns the command info for the given command
func (rc *RemoteConnection) GetCommand(cmd command.Command) (command.CommandInfo, bool) {
	info, ok := rc.commandMap[cmd]
	return info, ok
}

// Execute executes a command with the given parameters
func (rc *RemoteConnection) Execute(cmd command.Command, params map[string]interface{}) (map[string]interface{}, error) {
	cmdInfo, ok := rc.GetCommand(cmd)
	if !ok {
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}

	path := cmdInfo.Path
	for k, v := range params {
		path = strings.ReplaceAll(path, fmt.Sprintf("$%s", k), fmt.Sprint(v))
	}

	url := fmt.Sprintf("%s%s", rc.config.RemoteServerAddr, path)

	var body io.Reader
	if cmdInfo.Method == command.HTTPPost {
		jsonData, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %v", err)
		}
		body = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(string(cmdInfo.Method), url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	rc.addHeaders(req)

	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result, nil
}

// Close closes the connection
func (rc *RemoteConnection) Close() error {
	rc.client.CloseIdleConnections()
	return nil
}

// addHeaders adds the required headers to the request
func (rc *RemoteConnection) addHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("User-Agent", rc.config.UserAgent)

	if rc.config.KeepAlive {
		req.Header.Set("Connection", "keep-alive")
	}

	if rc.proxyAuth != "" {
		req.Header.Set("Proxy-Authorization", fmt.Sprintf("Basic %s", rc.proxyAuth))
	}

	for k, v := range rc.extraHeaders {
		req.Header.Set(k, v)
	}
}
