package selenium

import (
	"github.com/Kcrong/selenium-go/remote/command"
)

const Version = "4.30.0.202502201302"

// BrowserType represents the type of browser
type BrowserType string

const (
	// Chrome browser
	Chrome BrowserType = "chrome"
	// Firefox browser
	Firefox BrowserType = "firefox"
	// Edge browser
	Edge BrowserType = "edge"
	// Safari browser
	Safari BrowserType = "safari"
	// IE browser
	IE BrowserType = "ie"
	// WebKitGTK browser
	WebKitGTK BrowserType = "webkitgtk"
	// WPEWebKit browser
	WPEWebKit BrowserType = "wpewebkit"
)

// WebDriver interface defines the basic operations that all WebDriver implementations must support
type WebDriver interface {
	// Execute executes a WebDriver command
	Execute(cmd command.Command, params map[string]interface{}) (map[string]interface{}, error)
	// Quit closes the browser and ends the session
	Quit() error
	// GetCapabilities returns the current session's capabilities
	GetCapabilities() (*Capabilities, error)
	// GetSessionID returns the current session ID
	GetSessionID() string
}

// Service represents the basic service configuration for WebDriver implementations
type Service struct {
	// Executable path to the driver executable
	Executable string
	// Port on which the driver should listen
	Port int
	// Args additional command-line arguments to pass to the driver
	Args []string
	// Env environment variables to set for the driver process
	Env map[string]string
	// Start timeout in seconds
	StartTimeout int
}

// LogLevel represents the logging level
type LogLevel string

const (
	// LogOff disables logging
	LogOff LogLevel = "OFF"
	// LogSevere logs severe errors
	LogSevere LogLevel = "SEVERE"
	// LogWarning logs warnings
	LogWarning LogLevel = "WARNING"
	// LogInfo logs information
	LogInfo LogLevel = "INFO"
	// LogDebug logs debug information
	LogDebug LogLevel = "DEBUG"
)

// Options represents the basic options configuration for WebDriver implementations
type Options struct {
	// Arguments command-line arguments to pass to the browser
	Arguments []string
	// BinaryLocation path to the browser binary
	BinaryLocation string
	// Extensions list of extension files to install
	Extensions []string
	// LogLevel logging level
	LogLevel LogLevel
	// Proxy proxy configuration
	Proxy *Proxy
}

// Capabilities represents the basic desired capabilities for WebDriver implementations
type Capabilities struct {
	// BrowserName name of the browser
	BrowserName BrowserType
	// BrowserVersion version of the browser
	BrowserVersion string
	// PlatformName name of the platform
	PlatformName string
	// AcceptInsecureCerts whether to accept insecure certificates
	AcceptInsecureCerts bool
	// PageLoadStrategy page load strategy
	PageLoadStrategy string
	// Proxy proxy configuration
	Proxy *Proxy
	// SetWindowRect whether the session can set window rect
	SetWindowRect bool
	// Timeouts session timeouts
	Timeouts map[string]int
	// UnhandledPromptBehavior how to handle unexpected alerts
	UnhandledPromptBehavior string
}

// ProxyType represents the type of proxy
type ProxyType string

const (
	// DirectProxy direct connection (no proxy)
	DirectProxy ProxyType = "direct"
	// ManualProxy manual proxy settings
	ManualProxy ProxyType = "manual"
	// PacProxy proxy autoconfiguration
	PacProxy ProxyType = "pac"
	// AutodetectProxy auto detect proxy settings
	AutodetectProxy ProxyType = "autodetect"
	// SystemProxy use system proxy settings
	SystemProxy ProxyType = "system"
)

// Proxy represents proxy configuration for WebDriver
type Proxy struct {
	// ProxyType type of proxy
	ProxyType ProxyType
	// HTTPProxy HTTP proxy
	HTTPProxy string
	// HTTPSProxy HTTPS proxy
	HTTPSProxy string
	// FTPProxy FTP proxy
	FTPProxy string
	// SOCKSProxy SOCKS proxy
	SOCKSProxy string
	// SOCKSVersion SOCKS proxy version
	SOCKSVersion int
	// NoProxy addresses to bypass proxy
	NoProxy []string
}

// ActionChains represents a way to automate low level interactions
type ActionChains struct {
	// Driver WebDriver instance
	Driver WebDriver
	// Actions list of actions to perform
	Actions []Action
}

// ActionType represents the type of action
type ActionType string

const (
	// KeyDown key press down
	KeyDown ActionType = "keyDown"
	// KeyUp key release
	KeyUp ActionType = "keyUp"
	// MouseDown mouse button press
	MouseDown ActionType = "mouseDown"
	// MouseUp mouse button release
	MouseUp ActionType = "mouseUp"
	// MouseMove mouse movement
	MouseMove ActionType = "mouseMove"
	// MouseClick mouse click
	MouseClick ActionType = "click"
)

// Action represents a single action in an action chain
type Action struct {
	// Type type of action
	Type ActionType
	// Value value for the action (e.g., key to press, coordinates to move to)
	Value interface{}
	// Options additional options for the action
	Options map[string]interface{}
}
