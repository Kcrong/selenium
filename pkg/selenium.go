package pkg

import (
	"time"

	"github.com/Kcrong/selenium/pkg/chrome"
	"github.com/Kcrong/selenium/pkg/firefox"
	"github.com/Kcrong/selenium/pkg/log"
)

// Methods by which to find elements (W3C, though some locators are effectively synonyms now).
const (
	ByID              = "id"
	ByXPATH           = "xpath"
	ByLinkText        = "link text"
	ByPartialLinkText = "partial link text"
	ByName            = "name"
	ByTagName         = "tag name"
	ByClassName       = "class name"
	ByCSSSelector     = "css selector"
)

// MouseButton 은 마우스 버튼(왼쪽/중간/오른쪽)을 나타냅니다.
type MouseButton int

const (
	LeftButton MouseButton = iota
	MiddleButton
	RightButton
)

// Special keyboard keys, for SendKeys or KeyDown/KeyUp.
const (
	NullKey       = string('\ue000')
	CancelKey     = string('\ue001')
	HelpKey       = string('\ue002')
	BackspaceKey  = string('\ue003')
	TabKey        = string('\ue004')
	ClearKey      = string('\ue005')
	ReturnKey     = string('\ue006')
	EnterKey      = string('\ue007')
	ShiftKey      = string('\ue008')
	ControlKey    = string('\ue009')
	AltKey        = string('\ue00a')
	PauseKey      = string('\ue00b')
	EscapeKey     = string('\ue00c')
	SpaceKey      = string('\ue00d')
	PageUpKey     = string('\ue00e')
	PageDownKey   = string('\ue00f')
	EndKey        = string('\ue010')
	HomeKey       = string('\ue011')
	LeftArrowKey  = string('\ue012')
	UpArrowKey    = string('\ue013')
	RightArrowKey = string('\ue014')
	DownArrowKey  = string('\ue015')
	InsertKey     = string('\ue016')
	DeleteKey     = string('\ue017')
	SemicolonKey  = string('\ue018')
	EqualsKey     = string('\ue019')
	Numpad0Key    = string('\ue01a')
	Numpad1Key    = string('\ue01b')
	Numpad2Key    = string('\ue01c')
	Numpad3Key    = string('\ue01d')
	Numpad4Key    = string('\ue01e')
	Numpad5Key    = string('\ue01f')
	Numpad6Key    = string('\ue020')
	Numpad7Key    = string('\ue021')
	Numpad8Key    = string('\ue022')
	Numpad9Key    = string('\ue023')
	MultiplyKey   = string('\ue024')
	AddKey        = string('\ue025')
	SeparatorKey  = string('\ue026')
	SubstractKey  = string('\ue027')
	DecimalKey    = string('\ue028')
	DivideKey     = string('\ue029')
	F1Key         = string('\ue031')
	F2Key         = string('\ue032')
	F3Key         = string('\ue033')
	F4Key         = string('\ue034')
	F5Key         = string('\ue035')
	F6Key         = string('\ue036')
	F7Key         = string('\ue037')
	F8Key         = string('\ue038')
	F9Key         = string('\ue039')
	F10Key        = string('\ue03a')
	F11Key        = string('\ue03b')
	F12Key        = string('\ue03c')
	MetaKey       = string('\ue03d')
)

// Capabilities configures both the WebDriver process and the target browsers,
// with standard and browser-specific options.
type Capabilities map[string]interface{}

// AddChrome adds Chrome-specific capabilities to a W3C-capable environment.
func (c Capabilities) AddChrome(f chrome.Capabilities) {
	// 최신 W3C Selenium 스펙: "goog:chromeOptions" 키만.
	c[chrome.CapabilitiesKey] = f
}

// AddFirefox adds Firefox-specific capabilities to a W3C-capable environment.
func (c Capabilities) AddFirefox(f firefox.Capabilities) {
	c[firefox.CapabilitiesKey] = f
}

// AddProxy adds proxy configuration to the capabilities.
func (c Capabilities) AddProxy(p Proxy) {
	c["proxy"] = p
}

// AddLogging adds logging configuration to the capabilities.
func (c Capabilities) AddLogging(l log.Capabilities) {
	c[log.CapabilitiesKey] = l
}

// SetLogLevel sets the logging level of a component. It's a shortcut for
// passing a log.Capabilities instance to AddLogging.
func (c Capabilities) SetLogLevel(typ log.Type, level log.Level) {
	if _, ok := c[log.CapabilitiesKey]; !ok {
		c[log.CapabilitiesKey] = make(log.Capabilities)
	}
	m := c[log.CapabilitiesKey].(log.Capabilities)
	m[typ] = level
}

// Proxy specifies configuration for proxies in the browser. Set the key
// "proxy" in Capabilities to an instance of this type.
type Proxy struct {
	Type          ProxyType `json:"proxyType"`
	AutoconfigURL string    `json:"proxyAutoconfigUrl,omitempty"`
	FTP           string    `json:"ftpProxy,omitempty"`
	HTTP          string    `json:"httpProxy,omitempty"`
	SSL           string    `json:"sslProxy,omitempty"`
	SOCKS         string    `json:"socksProxy,omitempty"`
	SOCKSVersion  int       `json:"socksVersion,omitempty"`
	SOCKSUsername string    `json:"socksUsername,omitempty"`
	SOCKSPassword string    `json:"socksPassword,omitempty"`
	NoProxy       []string  `json:"noProxy,omitempty"`
	HTTPPort      int       `json:"httpProxyPort,omitempty"`
	SSLPort       int       `json:"sslProxyPort,omitempty"`
	SocksPort     int       `json:"socksProxyPort,omitempty"`
}

// ProxyType is an enumeration of the types of proxies available.
type ProxyType string

const (
	Direct     ProxyType = "direct"     // Direct connection - no proxy
	Manual               = "manual"     // Manually configured proxy
	Autodetect           = "autodetect" // Autodetect (WPAD)
	System               = "system"     // System settings
	PAC                  = "pac"        // Proxy autoconfig from a URL
)

// Status contains information returned by the Status method (if implemented).
type Status struct {
	Java struct {
		Version string
	}
	Build struct {
		Version, Revision, Time string
	}
	OS struct {
		Arch, Name, Version string
	}

	// W3C (GeckoDriver) fields:
	Ready   bool
	Message string
}

// Point is a 2D point.
type Point struct {
	X, Y int
}

// Size is a size of HTML element.
type Size struct {
	Width, Height int
}

// Cookie represents an HTTP cookie.
type Cookie struct {
	Name     string   `json:"name"`
	Value    string   `json:"value"`
	Path     string   `json:"path"`
	Domain   string   `json:"domain"`
	Secure   bool     `json:"secure"`
	Expiry   uint     `json:"expiry"`
	HTTPOnly bool     `json:"httpOnly"`
	SameSite SameSite `json:"sameSite,omitempty"`
}

// SameSite is the type for the SameSite field in Cookie.
type SameSite string

const (
	SameSiteNone   SameSite = "None"
	SameSiteLax    SameSite = "Lax"
	SameSiteStrict SameSite = "Strict"
	SameSiteEmpty  SameSite = ""
)

// PointerType is the type of pointer used by StorePointerActions.
// According to the W3C spec, there are 3 pointer types.
type PointerType string

const (
	MousePointer PointerType = "mouse"
	PenPointer               = "pen"
	TouchPointer             = "touch"
)

// PointerMoveOrigin controls how the offset for
// the pointer move action is calculated.
type PointerMoveOrigin string

const (
	FromViewport PointerMoveOrigin = "viewport"
	FromPointer                    = "pointer"
)

// KeyAction and PointerAction are for W3C Actions API.
type KeyAction map[string]interface{}

type PointerAction map[string]interface{}

// Actions stores KeyActions and PointerActions for later execution.
type Actions []map[string]interface{}

// WebDriver defines methods supported by WebDriver drivers (Selenium 4, W3C).
type WebDriver interface {
	NewSession() (string, error)
	SessionID() string
	SwitchSession(sessionID string) error
	Capabilities() (Capabilities, error)

	SetAsyncScriptTimeout(timeout time.Duration) error
	SetImplicitWaitTimeout(timeout time.Duration) error
	SetPageLoadTimeout(timeout time.Duration) error
	Quit() error

	CurrentWindowHandle() (string, error)
	WindowHandles() ([]string, error)
	CurrentURL() (string, error)
	Title() (string, error)
	PageSource() (string, error)
	Close() error
	SwitchFrame(frame interface{}) error
	SwitchWindow(name string) error
	CloseCurrentWindow() error
	MaximizeWindow(name string) error
	ResizeWindow(name string, width, height int) error

	Get(url string) error
	Forward() error
	Back() error
	Refresh() error

	FindElement(by, value string) (WebElement, error)
	FindElements(by, value string) ([]WebElement, error)
	ActiveElement() (WebElement, error)

	DecodeElement([]byte) (WebElement, error)
	DecodeElements([]byte) ([]WebElement, error)

	GetCookies() ([]Cookie, error)
	GetCookie(name string) (Cookie, error)
	AddCookie(cookie *Cookie) error
	DeleteAllCookies() error
	DeleteCookie(name string) error

	Click(button int) error
	DoubleClick() error
	ButtonDown() error
	ButtonUp() error

	StoreKeyActions(inputID string, actions ...KeyAction)
	StorePointerActions(inputID string, pointer PointerType, actions ...PointerAction)
	PerformActions() error
	ReleaseActions() error

	KeyDown(keys string) error
	KeyUp(keys string) error

	Screenshot() ([]byte, error)
	Log(typ log.Type) ([]log.Message, error)

	DismissAlert() error
	AcceptAlert() error
	AlertText() (string, error)
	SetAlertText(text string) error

	ExecuteScript(script string, args []interface{}) (interface{}, error)
	ExecuteScriptAsync(script string, args []interface{}) (interface{}, error)
	ExecuteScriptRaw(script string, args []interface{}) ([]byte, error)
	ExecuteScriptAsyncRaw(script string, args []interface{}) ([]byte, error)

	WaitWithTimeoutAndInterval(condition Condition, timeout, interval time.Duration) error
	WaitWithTimeout(condition Condition, timeout time.Duration) error
	Wait(condition Condition) error
}

// WebElement defines methods supported by web elements (Selenium 4, W3C).
type WebElement interface {
	Click() error
	SendKeys(keys string) error
	Submit() error
	Clear() error
	MoveTo(xOffset, yOffset int) error

	FindElement(by, value string) (WebElement, error)
	FindElements(by, value string) ([]WebElement, error)

	TagName() (string, error)
	Text() (string, error)
	IsSelected() (bool, error)
	IsEnabled() (bool, error)
	IsDisplayed() (bool, error)
	GetAttribute(name string) (string, error)
	GetProperty(name string) (string, error)
	Location() (*Point, error)
	Size() (*Size, error)
	CSSProperty(name string) (string, error)
	Screenshot(scroll bool) ([]byte, error)
}
