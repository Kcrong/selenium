package chromium

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Kcrong/selenigo"
)

const (
	// OptionsKey is the capability key for chrome options.
	OptionsKey = "goog:chromeOptions"
	// BrowserName is the name of the Chrome browser.
	BrowserName = "chrome"
)

// ErrEmptyExtension is returned when an empty extension is provided.
var ErrEmptyExtension = errors.New("encoded extension cannot be empty")

// Options contains the base options for Chromium-based browsers.
type Options struct {
	binaryLocation   string
	extensionFiles   []string
	extensions       []string
	experimentalOpts map[string]interface{}
	debuggerAddress  string
	arguments        []string
}

// NewOptions creates a new Chromium options instance.
func NewOptions() *Options {
	return &Options{
		binaryLocation:   "",
		extensionFiles:   make([]string, 0),
		extensions:       make([]string, 0),
		experimentalOpts: make(map[string]interface{}),
		debuggerAddress:  "",
		arguments:        make([]string, 0),
	}
}

// SetBinaryLocation sets the path to Chromium binary.
func (o *Options) SetBinaryLocation(path string) {
	o.binaryLocation = path
}

// GetBinaryLocation returns the path to Chromium binary.
func (o *Options) GetBinaryLocation() string {
	return o.binaryLocation
}

// SetDebuggerAddress sets the address of the remote devtools instance.
func (o *Options) SetDebuggerAddress(address string) {
	o.debuggerAddress = address
}

// GetDebuggerAddress returns the address of the remote devtools instance.
func (o *Options) GetDebuggerAddress() string {
	return o.debuggerAddress
}

// AddExtension adds the path to the extension to load.
func (o *Options) AddExtension(path string) error {
	if path == "" {
		return ErrEmptyExtension
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", err, absPath)
	}

	o.extensionFiles = append(o.extensionFiles, absPath)

	return nil
}

// AddEncodedExtension adds Base64 encoded string with extension data.
func (o *Options) AddEncodedExtension(encodedExtension string) error {
	if encodedExtension == "" {
		return ErrEmptyExtension
	}

	o.extensions = append(o.extensions, encodedExtension)

	return nil
}

// AddArgument adds a command-line argument.
func (o *Options) AddArgument(arg string) {
	o.arguments = append(o.arguments, arg)
}

// AddExperimentalOption adds an experimental option.
func (o *Options) AddExperimentalOption(name string, value interface{}) {
	o.experimentalOpts[name] = value
}

// getEncodedExtensions returns a list of encoded extensions.
func (o *Options) getEncodedExtensions() ([]string, error) {
	encodedExts := make([]string, 0)

	// Add pre-encoded extensions
	encodedExts = append(encodedExts, o.extensions...)

	// Encode extension files
	for _, path := range o.extensionFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read extension file %s: %w", path, err)
		}

		encoded := base64.StdEncoding.EncodeToString(data)
		encodedExts = append(encodedExts, encoded)
	}

	return encodedExts, nil
}

// ToCapabilities converts the options to a capabilities map.
func (o *Options) ToCapabilities() map[string]interface{} {
	caps := selenigo.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	chromeOptions := make(map[string]interface{})

	// Add experimental options
	for k, v := range o.experimentalOpts {
		chromeOptions[k] = v
	}

	// Add extensions
	if encodedExts, err := o.getEncodedExtensions(); err == nil && len(encodedExts) > 0 {
		chromeOptions["extensions"] = encodedExts
	}

	// Add binary location if set
	if o.binaryLocation != "" {
		chromeOptions["binary"] = o.binaryLocation
	}

	// Add arguments if any
	if len(o.arguments) > 0 {
		chromeOptions["args"] = o.arguments
	}

	// Add debugger address if set
	if o.debuggerAddress != "" {
		chromeOptions["debuggerAddress"] = o.debuggerAddress
	}

	caps.SetBrowserOptions(OptionsKey, chromeOptions)

	return caps.ToCapabilities()
}

// GetExtensions returns a list of encoded extensions.
func (o *Options) GetExtensions() ([]string, error) {
	return o.getEncodedExtensions()
}

// GetExperimentalOptions returns a map of experimental options.
func (o *Options) GetExperimentalOptions() map[string]interface{} {
	return o.experimentalOpts
}

// DefaultCapabilities returns the default capabilities for Chromium.
func (o *Options) DefaultCapabilities() map[string]interface{} {
	caps := selenigo.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	return caps.ToCapabilities()
}
