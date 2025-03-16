package chrome

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Kcrong/selenium"
)

const (
	// OptionsKey is the capability key for chrome options.
	OptionsKey = "goog:chromeOptions"
)

// Options contains the options for Chrome browser.
type Options struct {
	experimentalOpts map[string]interface{}
	mobileOptions    map[string]interface{}
	binaryLocation   string
	debuggerAddress  string
	extensionFiles   []string
	extensions       []string
	arguments        []string
}

// NewOptions creates a new Chrome options instance.
func NewOptions() *Options {
	//nolint:exhaustruct // Initialize required fields only for better readability.
	return &Options{
		experimentalOpts: make(map[string]interface{}),
		mobileOptions:    make(map[string]interface{}),
		extensionFiles:   make([]string, 0),
		extensions:       make([]string, 0),
		arguments:        make([]string, 0),
	}
}

// SetBinaryLocation sets the path to Chrome binary.
func (o *Options) SetBinaryLocation(path string) {
	o.binaryLocation = path
}

// GetBinaryLocation returns the path to Chrome binary.
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
		return fmt.Errorf("%w: %s", os.ErrNotExist, absPath)
	}

	o.extensionFiles = append(o.extensionFiles, absPath)

	return nil
}

var ErrEmptyExtension = errors.New("encoded extension cannot be empty")

// AddEncodedExtension adds Base64 encoded string with extension data.
func (o *Options) AddEncodedExtension(encodedExtension string) error {
	if encodedExtension == "" {
		return ErrEmptyExtension
	}

	o.extensions = append(o.extensions, encodedExtension)

	return nil
}

// AddArgument adds a command-line argument to pass to Chrome.
func (o *Options) AddArgument(arg string) {
	o.arguments = append(o.arguments, arg)
}

// AddExperimentalOption adds an experimental option to Chrome.
func (o *Options) AddExperimentalOption(name string, value interface{}) {
	o.experimentalOpts[name] = value
}

// EnableMobile enables mobile emulation.
func (o *Options) EnableMobile(androidPackage, androidActivity, deviceSerial string) {
	if androidPackage != "" {
		o.mobileOptions["androidPackage"] = androidPackage
	}

	if androidActivity != "" {
		o.mobileOptions["androidActivity"] = androidActivity
	}

	if deviceSerial != "" {
		o.mobileOptions["deviceSerial"] = deviceSerial
	}
}

// getEncodedExtensions returns a list of encoded extensions.
func (o *Options) getEncodedExtensions() ([]string, error) {
	encodedExts := make([]string, 0, len(o.extensions)+len(o.extensionFiles))

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
func (o *Options) ToCapabilities() selenium.Capabilities {
	chromeOptions := make(map[string]interface{})

	// Add experimental options
	for k, v := range o.experimentalOpts {
		chromeOptions[k] = v
	}

	// Add mobile options if any
	if len(o.mobileOptions) > 0 {
		for k, v := range o.mobileOptions {
			chromeOptions[k] = v
		}
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

	//nolint:exhaustruct // Use explicit fields for better readability.
	return selenium.Capabilities{
		BrowserName:  selenium.Chrome,
		PlatformName: "ANY",
		BrowserOptions: map[string]interface{}{
			OptionsKey: chromeOptions,
		},
	}
}
