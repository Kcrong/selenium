package selenium

import (
	"reflect"

	"github.com/Kcrong/caseconv"
)

// PageLoadStrategy represents different strategies for page loading.
type PageLoadStrategy string

const (
	// Normal waits for all resources to download.
	Normal PageLoadStrategy = "normal"
	// Eager DOM access is ready, but other resources like images may still be loading.
	Eager PageLoadStrategy = "eager"
	// None does not block WebDriver at all.
	None PageLoadStrategy = "none"
)

// BrowserType represents the type of browser.
type BrowserType string

const (
	// Chrome browser.
	Chrome BrowserType = "chrome"
	// Firefox browser.
	Firefox BrowserType = "firefox"
	// Edge browser.
	Edge BrowserType = "edge"
	// Safari browser.
	Safari BrowserType = "safari"
	// IE browser.
	IE BrowserType = "ie"
	// WebKitGTK browser.
	WebKitGTK BrowserType = "webkitgtk"
	// WPEWebKit browser.
	WPEWebKit BrowserType = "wpewebkit"
)

type HandlePromptBehaviorType string

const (
	// HandlePromptBehaviorTypeDismiss dismisses the alert.
	HandlePromptBehaviorTypeDismiss HandlePromptBehaviorType = "dismiss"
	// HandlePromptBehaviorTypeAccept accepts the alert.
	HandlePromptBehaviorTypeAccept HandlePromptBehaviorType = "accept"
	// HandlePromptBehaviorTypeDismissAndNotify dismisses the alert and notifies the user.
	HandlePromptBehaviorTypeDismissAndNotify HandlePromptBehaviorType = "dismiss and notify"
	// HandlePromptBehaviorTypeAcceptAndNotify accepts the alert and notifies the user.
	HandlePromptBehaviorTypeAcceptAndNotify HandlePromptBehaviorType = "accept and notify"
	// HandlePromptBehaviorTypeIgnore ignores the alert.
	HandlePromptBehaviorTypeIgnore HandlePromptBehaviorType = "ignore"
)

const (
	PlatformANY = "ANY"
)

// Capabilities represents the basic desired capabilities for WebDriver implementations.
type Capabilities struct {
	Proxy                     Proxy `json:"proxy"`
	BrowserOptions            map[string]interface{}
	UnhandledPromptBehavior   HandlePromptBehaviorType `json:"unhandledPromptBehavior"`
	PageLoadStrategy          PageLoadStrategy         `json:"pageLoadStrategy"`
	PlatformName              string
	BrowserName               BrowserType `json:"browserName"`
	BidiWebSocketURL          string      `json:"webSocketUrl"`
	BrowserVersion            string      `json:"browserVersion"`
	Timeouts                  Timeouts    `json:"timeouts"`
	AcceptInsecureCerts       bool        `json:"acceptInsecureCerts"`
	StrictFileInteractAbility bool        `json:"strictFileInteractbility"`
	SetWindowRect             bool        `json:"setWindowRect"`
	IsDownloadsEnabled        bool        `json:"se:downloadsEnabled"`
	IsJavaScriptEnabled       bool        `json:"javascriptEnabled"`
}

// BaseOptions represents the base options for all browser drivers.
type BaseOptions struct {
	Proxy        Proxy
	Capabilities Capabilities
}

// NewCapabilities creates a new BaseOptions instance.
func NewCapabilities() *BaseOptions {
	return &BaseOptions{
		Proxy: nil,
		//nolint:exhaustruct // This will be set by the user.
		Capabilities: Capabilities{},
	}
}

func (o *BaseOptions) SetBrowserOptions(browserOptionsKey string, value map[string]interface{}) {
	o.Capabilities.BrowserOptions[browserOptionsKey] = value
}

// ToCapabilities returns the capabilities as a map.
func (o *BaseOptions) ToCapabilities() map[string]interface{} {
	caps := make(map[string]interface{})

	caps["browserName"] = o.Capabilities.BrowserName
	caps["browserVersion"] = o.Capabilities.BrowserVersion
	caps["platformName"] = o.Capabilities.PlatformName
	caps["platform"] = o.Capabilities.PlatformName
	caps["acceptInsecureCerts"] = o.Capabilities.AcceptInsecureCerts
	caps["pageLoadStrategy"] = o.Capabilities.PageLoadStrategy
	caps["strictFileInteractability"] = o.Capabilities.StrictFileInteractAbility
	caps["proxy"] = o.Capabilities.Proxy.ToCapabilities()
	caps["setWindowRect"] = o.Capabilities.SetWindowRect
	caps["timeouts"] = o.Capabilities.Timeouts.ToCapabilities()
	caps["unhandledPromptBehavior"] = o.Capabilities.UnhandledPromptBehavior
	caps["webSocketUrl"] = o.Capabilities.BidiWebSocketURL
	caps["se:downloadsEnabled"] = o.Capabilities.IsDownloadsEnabled

	for k, v := range o.Capabilities.BrowserOptions {
		caps[k] = v
	}

	return caps
}

type Convertible interface {
	ToCapabilities() map[string]interface{}
}

func Marshal(c any) map[string]interface{} {
	if c == nil {
		return nil
	}

	if c, ok := c.(Convertible); ok {
		return c.ToCapabilities()
	}

	return StructToMap(c)
}

// StructToMap converts any struct to map with camelCase keys.
func StructToMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()

	for i := range val.NumField() {
		field := typ.Field(i)

		// Ignore unexported fields
		if !val.Field(i).CanInterface() {
			continue
		}

		key := caseconv.ToCamelCase(field.Name)
		result[key] = val.Field(i).Interface()
	}

	return result
}
