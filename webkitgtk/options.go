package webkitgtk

import (
	"github.com/Kcrong/selenigo"
)

const (
	// OptionsKey is the capability key for WebKitGTK options.
	OptionsKey = "webkitgtk:browserOptions"
	// BrowserName is the name of the WebKitGTK browser.
	BrowserName = selenigo.BrowserType("webkitgtk")
)

// Options contains the options for WebKitGTK browser.
type Options struct {
	binaryLocation           string
	arguments                []string
	overlayScrollbarsEnabled bool
}

// NewOptions creates a new WebKitGTK options instance.
func NewOptions() *Options {
	return &Options{
		binaryLocation:           "",
		overlayScrollbarsEnabled: true,
		arguments:                make([]string, 0),
	}
}

// SetBinaryLocation sets the path to WebKitGTK binary.
func (o *Options) SetBinaryLocation(path string) {
	o.binaryLocation = path
}

// GetBinaryLocation returns the path to WebKitGTK binary.
func (o *Options) GetBinaryLocation() string {
	return o.binaryLocation
}

// SetOverlayScrollbarsEnabled sets whether to enable overlay scrollbars.
func (o *Options) SetOverlayScrollbarsEnabled(enable bool) {
	o.overlayScrollbarsEnabled = enable
}

// GetOverlayScrollbarsEnabled returns whether overlay scrollbars are enabled.
func (o *Options) GetOverlayScrollbarsEnabled() bool {
	return o.overlayScrollbarsEnabled
}

// AddArgument adds a command-line argument.
func (o *Options) AddArgument(arg string) {
	o.arguments = append(o.arguments, arg)
}

// ToCapabilities converts the options to a capabilities map.
func (o *Options) ToCapabilities() map[string]interface{} {
	caps := selenigo.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	browserOptions := make(map[string]interface{})

	if o.binaryLocation != "" {
		browserOptions["binary"] = o.binaryLocation
	}

	if len(o.arguments) > 0 {
		browserOptions["args"] = o.arguments
	}

	browserOptions["useOverlayScrollbars"] = o.overlayScrollbarsEnabled

	caps.SetBrowserOptions(OptionsKey, browserOptions)

	return caps.ToCapabilities()
}

// DefaultCapabilities returns the default capabilities for WebKitGTK.
func (o *Options) DefaultCapabilities() map[string]interface{} {
	caps := selenigo.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	return caps.ToCapabilities()
}
