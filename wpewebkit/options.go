package wpewebkit

import (
	"errors"

	"github.com/Kcrong/selenium"
)

const (
	// OptionsKey is the capability key for WPEWebKit options.
	OptionsKey = "wpe:browserOptions"
	// BrowserName is the name of the WPEWebKit browser.
	BrowserName = selenium.BrowserType("wpewebkit")
)

// ErrInvalidBinaryLocation is returned when an invalid binary location is provided.
var ErrInvalidBinaryLocation = errors.New("binary location must be a string")

// Options contains the options for WPEWebKit browser.
type Options struct {
	binaryLocation string
	arguments      []string
}

// NewOptions creates a new WPEWebKit options instance.
func NewOptions() *Options {
	return &Options{
		binaryLocation: "",
		arguments:      make([]string, 0),
	}
}

// SetBinaryLocation sets the path to WPEWebKit binary.
func (o *Options) SetBinaryLocation(path string) error {
	if path == "" {
		return ErrInvalidBinaryLocation
	}

	o.binaryLocation = path

	return nil
}

// GetBinaryLocation returns the path to WPEWebKit binary.
func (o *Options) GetBinaryLocation() string {
	return o.binaryLocation
}

// AddArgument adds a command-line argument.
func (o *Options) AddArgument(arg string) {
	o.arguments = append(o.arguments, arg)
}

// ToCapabilities converts the options to a capabilities map.
func (o *Options) ToCapabilities() map[string]interface{} {
	caps := selenium.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	browserOptions := make(map[string]interface{})

	if o.binaryLocation != "" {
		browserOptions["binary"] = o.binaryLocation
	}

	if len(o.arguments) > 0 {
		browserOptions["args"] = o.arguments
	}

	caps.SetBrowserOptions(OptionsKey, browserOptions)

	return caps.ToCapabilities()
}

// DefaultCapabilities returns the default capabilities for WPEWebKit.
func (o *Options) DefaultCapabilities() map[string]interface{} {
	caps := selenium.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	return caps.ToCapabilities()
}
