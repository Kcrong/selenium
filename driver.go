package selenium

import (
	"context"

	"github.com/Kcrong/selenium/remote/command"
)

// Cookie represents a browser cookie.
type Cookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Path     string  `json:"path,omitempty"`
	Domain   string  `json:"domain,omitempty"`
	SameSite string  `json:"sameSite,omitempty"`
	Expiry   float64 `json:"expiry,omitempty"`
	Secure   bool    `json:"secure,omitempty"`
	HTTPOnly bool    `json:"httpOnly,omitempty"`
}

// WebElement represents a DOM element.
type WebElement interface {
	GetID() string
	Click(ctx context.Context) error
	SendKeys(ctx context.Context, keys string) error
	Clear(ctx context.Context) error
	GetText(ctx context.Context) (string, error)
	GetAttribute(ctx context.Context, name string) (string, error)
	GetProperty(ctx context.Context, name string) (string, error)
	IsSelected(ctx context.Context) (bool, error)
	IsEnabled(ctx context.Context) (bool, error)
	IsDisplayed(ctx context.Context) (bool, error)
	Screenshot(ctx context.Context) ([]byte, error)
	FindElement(ctx context.Context, by *By, value string) (WebElement, error)
	FindElements(ctx context.Context, by *By, value string) ([]WebElement, error)
}

// WebDriver interface defines the operations that all WebDriver implementations must support.
type WebDriver interface {
	// Session management

	// NewSession creates a new session with the desired capabilities.
	NewSession(ctx context.Context, capabilities Convertible) error
	// DeleteSession deletes the current session.
	DeleteSession(ctx context.Context) error
	// GetSessionID returns the current session ID.
	GetSessionID() string
	// GetCapabilities returns the capabilities of the current session.
	GetCapabilities() Convertible

	// Navigation

	// Get navigates to the given URL.
	Get(ctx context.Context, url string) error
	// GetCurrentURL returns the current URL.
	GetCurrentURL(ctx context.Context) (string, error)
	// Back navigates previous page in the browser history.
	Back(ctx context.Context) error
	// Forward navigates next page in the browser history.
	Forward(ctx context.Context) error
	// Refresh refreshes the current page.
	Refresh(ctx context.Context) error

	// Window management

	// GetWindowHandle returns the current window handle.
	GetWindowHandle(ctx context.Context) (string, error)
	// GetWindowHandles returns the list of all window handles available to the session.
	GetWindowHandles(ctx context.Context) ([]string, error)
	// MaximizeWindow maximizes the current window.
	MaximizeWindow(ctx context.Context) error
	// MinimizeWindow minimizes the current window.
	MinimizeWindow(ctx context.Context) error
	// FullscreenWindow makes the current window fullscreen.
	FullscreenWindow(ctx context.Context) error
	// SetWindowRect sets the position and size of the current window.
	SetWindowRect(ctx context.Context, x, y, width, height int) error
	// GetWindowRect gets the position and size of the current window.
	GetWindowRect(ctx context.Context) (x, y, width, height int, err error)

	// Element interaction

	// FindElement finds the first element that matches the given selector.
	FindElement(ctx context.Context, by *By, value string) (WebElement, error)
	// FindElements finds all elements that match the given selector.
	FindElements(ctx context.Context, by *By, value string) ([]WebElement, error)
	// GetActiveElement returns the currently active element.
	GetActiveElement(ctx context.Context) (WebElement, error)

	// Timeouts

	// SetTimeouts sets the timeouts for the current session.
	SetTimeouts(ctx context.Context, timeouts *Timeouts) error
	// GetTimeouts returns the timeouts for the current session.
	GetTimeouts(ctx context.Context) (*Timeouts, error)

	// Script execution

	// ExecuteScript executes the given script.
	ExecuteScript(ctx context.Context, script string, args []interface{}) (interface{}, error)
	// ExecuteAsyncScript executes the given async script.
	ExecuteAsyncScript(ctx context.Context, script string, args []interface{}) (interface{}, error)

	// Cookie management

	// GetCookies returns all cookies for the current page.
	GetCookies(ctx context.Context) ([]Cookie, error)
	// GetCookie returns the cookie with the given name.
	GetCookie(ctx context.Context, name string) (*Cookie, error)
	// AddCookie adds a cookie.
	AddCookie(ctx context.Context, cookie *Cookie) error
	// DeleteCookie deletes the cookie with the given name.
	DeleteCookie(ctx context.Context, name string) error
	// DeleteAllCookies deletes all cookies.
	DeleteAllCookies(ctx context.Context) error

	// Alert handling

	// AcceptAlert accepts the currently displayed alert dialog.
	AcceptAlert(ctx context.Context) error
	// DismissAlert dismisses the currently displayed alert dialog.
	DismissAlert(ctx context.Context) error
	// GetAlertText returns the text of the currently displayed alert dialog.
	GetAlertText(ctx context.Context) (string, error)

	// Screenshot takes a screenshot of the current window.
	Screenshot(ctx context.Context) ([]byte, error)

	// Other utilities

	// GetTitle returns the current page title.
	GetTitle(ctx context.Context) (string, error)
	// GetPageSource returns the current page source.
	GetPageSource(ctx context.Context) (string, error)

	// Close closes the current window.
	Close(ctx context.Context) error

	// Quit closes the browser and shuts down the WebDriver server.
	Quit(ctx context.Context) error

	// Execute executes a command with the given parameters.
	Execute(ctx context.Context, cmd command.Command, params map[string]interface{}) (map[string]interface{}, error)
}
