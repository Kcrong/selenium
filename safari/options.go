package safari

import (
	"github.com/Kcrong/selenigo"
)

const (
	// OptionsKey is the capability key for safari options.
	OptionsKey = "safari:options"
	// BrowserName is the name of the Safari browser.
	BrowserName = selenigo.BrowserType("safari")
	// TechnologyPreviewBrowserName is the name of the Safari Technology Preview browser.
	TechnologyPreviewBrowserName = selenigo.BrowserType("Safari Technology Preview")

	// AutomaticInspection is the capability key for automatic inspection.
	AutomaticInspection = "safari:automaticInspection"
	// AutomaticProfiling is the capability key for automatic profiling.
	AutomaticProfiling = "safari:automaticProfiling"
)

// Options contains the options for Safari browser.
type Options struct {
	caps map[string]interface{}
}

// NewOptions creates a new Safari options instance.
func NewOptions() *Options {
	return &Options{
		caps: make(map[string]interface{}),
	}
}

// SetAutomaticInspection sets whether to enable automatic inspection.
func (o *Options) SetAutomaticInspection(enable bool) {
	o.caps[AutomaticInspection] = enable
}

// GetAutomaticInspection returns whether automatic inspection is enabled.
func (o *Options) GetAutomaticInspection() bool {
	if enable, ok := o.caps[AutomaticInspection].(bool); ok {
		return enable
	}

	return false
}

// SetAutomaticProfiling sets whether to enable automatic profiling.
func (o *Options) SetAutomaticProfiling(enable bool) {
	o.caps[AutomaticProfiling] = enable
}

// GetAutomaticProfiling returns whether automatic profiling is enabled.
func (o *Options) GetAutomaticProfiling() bool {
	if enable, ok := o.caps[AutomaticProfiling].(bool); ok {
		return enable
	}

	return false
}

// SetUseTechnologyPreview sets whether to use Safari Technology Preview.
func (o *Options) SetUseTechnologyPreview(use bool) {
	if use {
		o.caps["browserName"] = TechnologyPreviewBrowserName
	} else {
		o.caps["browserName"] = BrowserName
	}
}

// GetUseTechnologyPreview returns whether Safari Technology Preview is being used.
func (o *Options) GetUseTechnologyPreview() bool {
	if browserName, ok := o.caps["browserName"].(selenigo.BrowserType); ok {
		return browserName == TechnologyPreviewBrowserName
	}

	return false
}

// ToCapabilities converts the options to a capabilities map.
func (o *Options) ToCapabilities() map[string]interface{} {
	caps := selenigo.NewCapabilities()

	// Set browser name if not already set
	if browserName, ok := o.caps["browserName"].(selenigo.BrowserType); ok {
		caps.Capabilities.BrowserName = browserName
	} else {
		caps.Capabilities.BrowserName = BrowserName
	}

	// Add all options to capabilities
	for k, v := range o.caps {
		if k != "browserName" {
			caps.SetBrowserOptions(OptionsKey, map[string]interface{}{k: v})
		}
	}

	return caps.ToCapabilities()
}

// DefaultCapabilities returns the default capabilities for Safari.
func (o *Options) DefaultCapabilities() map[string]interface{} {
	caps := selenigo.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	return caps.ToCapabilities()
}
