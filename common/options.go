package common

import (
	"fmt"
)

// PageLoadStrategy represents different strategies for page loading
type PageLoadStrategy string

const (
	// Normal waits for all resources to download
	Normal PageLoadStrategy = "normal"
	// Eager DOM access is ready, but other resources like images may still be loading
	Eager PageLoadStrategy = "eager"
	// None does not block WebDriver at all
	None PageLoadStrategy = "none"
)

// Capabilities represents the capabilities of the browser session
type Capabilities map[string]interface{}

// BaseOptions represents the base options for all browser drivers
type BaseOptions struct {
	capabilities Capabilities
	proxy        *Proxy
}

// NewBaseOptions creates a new BaseOptions instance
func NewBaseOptions() *BaseOptions {
	return &BaseOptions{
		capabilities: make(Capabilities),
	}
}

// GetCapability returns the value of a capability
func (o *BaseOptions) GetCapability(name string) interface{} {
	return o.capabilities[name]
}

// SetCapability sets a capability
func (o *BaseOptions) SetCapability(name string, value interface{}) {
	o.capabilities[name] = value
}

// GetBrowserVersion returns the browser version
func (o *BaseOptions) GetBrowserVersion() string {
	if v, ok := o.capabilities["browserVersion"].(string); ok {
		return v
	}
	return ""
}

// SetBrowserVersion sets the browser version
func (o *BaseOptions) SetBrowserVersion(version string) {
	o.SetCapability("browserVersion", version)
}

// GetPlatformName returns the platform name
func (o *BaseOptions) GetPlatformName() string {
	if v, ok := o.capabilities["platformName"].(string); ok {
		return v
	}
	return ""
}

// SetPlatformName sets the platform name
func (o *BaseOptions) SetPlatformName(platform string) {
	o.SetCapability("platformName", platform)
}

// GetAcceptInsecureCerts returns whether insecure certificates are accepted
func (o *BaseOptions) GetAcceptInsecureCerts() bool {
	if v, ok := o.capabilities["acceptInsecureCerts"].(bool); ok {
		return v
	}
	return false
}

// SetAcceptInsecureCerts sets whether insecure certificates are accepted
func (o *BaseOptions) SetAcceptInsecureCerts(accept bool) {
	o.SetCapability("acceptInsecureCerts", accept)
}

// GetStrictFileInteractability returns whether strict file interactability is enabled
func (o *BaseOptions) GetStrictFileInteractability() bool {
	if v, ok := o.capabilities["strictFileInteractability"].(bool); ok {
		return v
	}
	return false
}

// SetStrictFileInteractability sets whether strict file interactability is enabled
func (o *BaseOptions) SetStrictFileInteractability(strict bool) {
	o.SetCapability("strictFileInteractability", strict)
}

// GetSetWindowRect returns whether window rect settings are enabled
func (o *BaseOptions) GetSetWindowRect() bool {
	if v, ok := o.capabilities["setWindowRect"].(bool); ok {
		return v
	}
	return false
}

// SetSetWindowRect sets whether window rect settings are enabled
func (o *BaseOptions) SetSetWindowRect(enabled bool) {
	o.SetCapability("setWindowRect", enabled)
}

// GetEnableBiDi returns whether BiDi support is enabled
func (o *BaseOptions) GetEnableBiDi() bool {
	if v, ok := o.capabilities["webSocketUrl"].(bool); ok {
		return v
	}
	if v, ok := o.capabilities["webSocketUrl"].(string); ok {
		return v != ""
	}
	return false
}

// SetEnableBiDi sets whether BiDi support is enabled
func (o *BaseOptions) SetEnableBiDi(enabled bool) {
	o.SetCapability("webSocketUrl", enabled)
}

// GetWebSocketUrl returns the WebSocket URL for BiDi support
func (o *BaseOptions) GetWebSocketUrl() string {
	if v, ok := o.capabilities["webSocketUrl"].(string); ok {
		return v
	}
	return ""
}

// GetPageLoadStrategy returns the page load strategy
func (o *BaseOptions) GetPageLoadStrategy() PageLoadStrategy {
	if v, ok := o.capabilities["pageLoadStrategy"].(string); ok {
		return PageLoadStrategy(v)
	}
	return Normal
}

// SetPageLoadStrategy sets the page load strategy
func (o *BaseOptions) SetPageLoadStrategy(strategy PageLoadStrategy) error {
	switch strategy {
	case Normal, Eager, None:
		o.SetCapability("pageLoadStrategy", string(strategy))
		return nil
	default:
		return fmt.Errorf("invalid page load strategy: %s", strategy)
	}
}

// GetUnhandledPromptBehavior returns the unhandled prompt behavior
func (o *BaseOptions) GetUnhandledPromptBehavior() string {
	if v, ok := o.capabilities["unhandledPromptBehavior"].(string); ok {
		return v
	}
	return "dismiss and notify"
}

// SetUnhandledPromptBehavior sets the unhandled prompt behavior
func (o *BaseOptions) SetUnhandledPromptBehavior(behavior string) error {
	validBehaviors := map[string]bool{
		"dismiss":            true,
		"accept":             true,
		"dismiss and notify": true,
		"accept and notify":  true,
		"ignore":             true,
	}

	if !validBehaviors[behavior] {
		return fmt.Errorf("invalid unhandled prompt behavior: %s", behavior)
	}

	o.SetCapability("unhandledPromptBehavior", behavior)
	return nil
}

// GetTimeouts returns the timeouts configuration
func (o *BaseOptions) GetTimeouts() map[string]interface{} {
	if v, ok := o.capabilities["timeouts"].(map[string]interface{}); ok {
		return v
	}
	return nil
}

// SetTimeouts sets the timeouts configuration
func (o *BaseOptions) SetTimeouts(timeouts map[string]interface{}) error {
	validKeys := map[string]bool{
		"implicit": true,
		"pageLoad": true,
		"script":   true,
	}

	for key := range timeouts {
		if !validKeys[key] {
			return fmt.Errorf("invalid timeout key: %s", key)
		}
	}

	o.SetCapability("timeouts", timeouts)
	return nil
}

// GetProxy returns the proxy configuration
func (o *BaseOptions) GetProxy() *Proxy {
	return o.proxy
}

// SetProxy sets the proxy configuration
func (o *BaseOptions) SetProxy(proxy *Proxy) error {
	if proxy == nil {
		return fmt.Errorf("proxy cannot be nil")
	}
	o.proxy = proxy
	o.SetCapability("proxy", proxy.ToCapabilities())
	return nil
}

// GetEnableDownloads returns whether downloads are enabled
func (o *BaseOptions) GetEnableDownloads() bool {
	if v, ok := o.capabilities["se:downloadsEnabled"].(bool); ok {
		return v
	}
	return false
}

// SetEnableDownloads sets whether downloads are enabled
func (o *BaseOptions) SetEnableDownloads(enabled bool) {
	o.SetCapability("se:downloadsEnabled", enabled)
}

// ToCapabilities returns the capabilities as a map
func (o *BaseOptions) ToCapabilities() Capabilities {
	return o.capabilities
}
