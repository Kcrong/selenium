package edge

import (
	"github.com/Kcrong/selenigo"
	"github.com/Kcrong/selenigo/chromium"
)

const (
	// OptionsKey is the capability key for edge options.
	OptionsKey = "ms:edgeOptions"
	// BrowserName is the name of the Microsoft Edge browser.
	BrowserName = "MicrosoftEdge"
)

// Options contains the options for Microsoft Edge browser.
type Options struct {
	*chromium.Options
	useWebView bool
}

// NewOptions creates a new Edge options instance.
func NewOptions() *Options {
	return &Options{
		Options:    chromium.NewOptions(),
		useWebView: false,
	}
}

// SetUseWebView sets whether to use WebView2.
func (o *Options) SetUseWebView(useWebView bool) {
	o.useWebView = useWebView
}

// GetUseWebView returns whether WebView2 is being used.
func (o *Options) GetUseWebView() bool {
	return o.useWebView
}

// ToCapabilities converts the options to a capabilities map.
func (o *Options) ToCapabilities() map[string]interface{} {
	caps := o.Options.ToCapabilities()

	if o.useWebView {
		caps["browserName"] = "webview2"
	}

	return caps
}

// DefaultCapabilities returns the default capabilities for Edge.
func (o *Options) DefaultCapabilities() map[string]interface{} {
	caps := selenigo.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	return caps.ToCapabilities()
}
