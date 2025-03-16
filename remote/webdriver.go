// Package remote provides WebDriver implementation.
package remote

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/Kcrong/selenium"
	"github.com/Kcrong/selenium/remote/command"
	"github.com/Kcrong/selenium/remote/connection"
	"github.com/Kcrong/selenium/remote/webelement"
)

// Common errors for WebDriver operations.
var (
	ErrFailedToCreateSession       = errors.New("failed to create session")
	ErrFailedToGetURL              = errors.New("failed to get current URL")
	ErrFailedToGetTitle            = errors.New("failed to get page title")
	ErrFailedToGetPageSource       = errors.New("failed to get page source")
	ErrFailedToGetWindowHandle     = errors.New("failed to get window handle")
	ErrFailedToGetWindowHandles    = errors.New("failed to get window handles")
	ErrFailedToGetCookies          = errors.New("failed to get cookies")
	ErrFailedToGetCookie           = errors.New("failed to get cookie")
	ErrFailedToGetScreenshot       = errors.New("failed to get screenshot")
	ErrFailedToGetSessionID        = errors.New("failed to get session ID")
	ErrFailedToGetCapabilities     = errors.New("failed to get capabilities")
	ErrFailedToConvertCapabilities = errors.New("failed to convert capabilities")
)

// Session represents a WebDriver session.
type Session struct {
	Capabilities selenium.Convertible
	SessionID    string
}

// WebDriver implements the WebDriver interface.
type WebDriver struct {
	capabilities selenium.Convertible
	conn         *connection.RemoteConnection
	sessionID    string
}

func (d *WebDriver) SetWindowRect(ctx context.Context, x, y, width, height int) error {
	_, err := d.Execute(ctx, command.Command("setWindowRect"), map[string]interface{}{
		"x":      x,
		"y":      y,
		"width":  width,
		"height": height,
	})

	return err
}

func (d *WebDriver) GetWindowRect(ctx context.Context) (x, y, width, height int, err error) {
	response, err := d.Execute(ctx, command.Command("getWindowRect"), nil)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	rect, ok := response["value"].(map[string]interface{})
	if !ok {
		return 0, 0, 0, 0, errors.New("failed to get window rect: invalid response format")
	}

	x, _ = rect["x"].(int)
	y, _ = rect["y"].(int)
	width, _ = rect["width"].(int)
	height, _ = rect["height"].(int)

	return x, y, width, height, nil
}

func (d *WebDriver) FindElement(ctx context.Context, by *selenium.By, value string) (selenium.WebElement, error) {
	response, err := d.Execute(ctx, command.Command("findElement"), map[string]interface{}{
		"using": by.GetFinder(value),
		"value": value,
	})
	if err != nil {
		return nil, err
	}

	if element, ok := response["value"].(map[string]interface{}); ok {
		if elementID, ok := element["ELEMENT"].(string); ok {
			return webelement.NewElement(elementID, d.sessionID, d.conn), nil
		}
	}

	return nil, fmt.Errorf("failed to find element: %v", response)
}

func (d *WebDriver) FindElements(ctx context.Context, by *selenium.By, value string) ([]selenium.WebElement, error) {
	response, err := d.Execute(ctx, command.Command("findElements"), map[string]interface{}{
		"using": by.GetFinder(value),
		"value": value,
	})
	if err != nil {
		return nil, err
	}

	if elements, ok := response["value"].([]interface{}); ok {
		result := make([]selenium.WebElement, len(elements))
		for i, element := range elements {
			if elementMap, ok := element.(map[string]interface{}); ok {
				if elementID, ok := elementMap["ELEMENT"].(string); ok {
					result[i] = webelement.NewElement(elementID, d.sessionID, d.conn)

					continue
				}
			}

			return nil, fmt.Errorf("invalid element at index %d: %v", i, element)
		}

		return result, nil
	}

	return nil, fmt.Errorf("failed to find elements: %v", response)
}

func (d *WebDriver) GetActiveElement(ctx context.Context) (selenium.WebElement, error) {
	response, err := d.Execute(ctx, command.GetActiveElement, nil)
	if err != nil {
		return nil, err
	}

	if element, ok := response["value"].(map[string]interface{}); ok {
		if elementID, ok := element["ELEMENT"].(string); ok {
			return webelement.NewElement(elementID, d.sessionID, d.conn), nil
		}
	}

	return nil, fmt.Errorf("failed to get active element: %v", response)
}

func (d *WebDriver) SetTimeouts(ctx context.Context, timeouts *selenium.Timeouts) error {
	params := timeouts.ToCapabilities()
	_, err := d.Execute(ctx, command.SetTimeouts, params)

	return err
}

func (d *WebDriver) GetTimeouts(ctx context.Context) (*selenium.Timeouts, error) {
	response, err := d.Execute(ctx, command.GetTimeouts, nil)
	if err != nil {
		return nil, err
	}

	if timeouts, ok := response["value"].(map[string]interface{}); ok {
		implicit, _ := timeouts["implicit"].(float64)
		pageLoad, _ := timeouts["pageLoad"].(float64)
		script, _ := timeouts["script"].(float64)

		return selenium.NewTimeouts(implicit, pageLoad, script), nil
	}

	return nil, fmt.Errorf("failed to get timeouts: %v", response)
}

func (d *WebDriver) AcceptAlert(ctx context.Context) error {
	_, err := d.Execute(ctx, command.W3CAcceptAlert, nil)

	return err
}

func (d *WebDriver) GetAlertText(ctx context.Context) (string, error) {
	response, err := d.Execute(ctx, command.W3CGetAlertText, nil)
	if err != nil {
		return "", err
	}

	if text, ok := response["value"].(string); ok {
		return text, nil
	}

	return "", fmt.Errorf("failed to get alert text: %v", response)
}

func (d *WebDriver) DismissAlert(ctx context.Context) error {
	_, err := d.Execute(ctx, command.W3CDismissAlert, nil)

	return err
}

var _ selenium.WebDriver = (*WebDriver)(nil)

// newSession creates a new browser session.
func (d *WebDriver) newSession(ctx context.Context) (*Session, error) {
	params := map[string]interface{}{
		"capabilities": map[string]interface{}{
			"alwaysMatch": d.capabilities,
		},
	}

	response, err := d.conn.Execute(ctx, "newSession", params)
	if err != nil {
		return nil, err
	}

	sessionID, ok := response["sessionId"].(string)
	if !ok {
		return nil, ErrFailedToGetSessionID
	}

	capabilities, ok := response["capabilities"].(map[string]interface{})
	if !ok {
		return nil, ErrFailedToGetCapabilities
	}

	return &Session{
		SessionID:    sessionID,
		Capabilities: selenium.RawConvertible(capabilities),
	}, nil
}

func (d *WebDriver) NewSession(ctx context.Context, capabilities selenium.Convertible) error {
	if d.sessionID != "" {
		return errors.New("session already exists")
	}

	d.capabilities = capabilities

	session, err := d.newSession(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCreateSession, err)
	}

	d.sessionID = session.SessionID
	d.capabilities = session.Capabilities

	return nil
}

// Execute executes a WebDriver command.
func (d *WebDriver) Execute(ctx context.Context, cmd command.Command, params map[string]interface{}) (map[string]interface{}, error) {
	if d.sessionID == "" {
		return nil, errors.New("no active session")
	}

	return d.conn.Execute(ctx, cmd, params)
}

// DeleteSession deletes the current session.
func (d *WebDriver) DeleteSession(ctx context.Context) error {
	if d.sessionID == "" {
		return nil
	}

	_, err := d.Execute(ctx, "deleteSession", map[string]interface{}{
		"sessionId": d.sessionID,
	})

	return err
}

// Quit closes the browser and ends the session.
func (d *WebDriver) Quit(ctx context.Context) error {
	if d.sessionID == "" {
		return nil
	}

	err := d.DeleteSession(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	d.sessionID = ""

	return nil
}

// GetCapabilities returns the current session's capabilities.
func (d *WebDriver) GetCapabilities() selenium.Convertible {
	return d.capabilities
}

// GetSessionID returns the current session ID.
func (d *WebDriver) GetSessionID() string {
	return d.sessionID
}

// Get navigates to the given URL.
func (d *WebDriver) Get(ctx context.Context, url string) error {
	_, err := d.Execute(ctx, "get", map[string]interface{}{
		"url": url,
	})

	return err
}

// GetCurrentURL returns the current URL.
func (d *WebDriver) GetCurrentURL(ctx context.Context) (string, error) {
	response, err := d.Execute(ctx, "getCurrentUrl", nil)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGetURL, err)
	}

	if url, ok := response["value"].(string); ok {
		return url, nil
	}

	return "", fmt.Errorf("%w: %v", ErrFailedToGetURL, response)
}

// GetTitle returns the title of the current page.
func (d *WebDriver) GetTitle(ctx context.Context) (string, error) {
	response, err := d.Execute(ctx, "getTitle", nil)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGetTitle, err)
	}

	if title, ok := response["value"].(string); ok {
		return title, nil
	}

	return "", fmt.Errorf("%w: %v", ErrFailedToGetTitle, response)
}

// GetPageSource returns the source of the current page.
func (d *WebDriver) GetPageSource(ctx context.Context) (string, error) {
	response, err := d.Execute(ctx, "getPageSource", nil)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGetPageSource, err)
	}

	if source, ok := response["value"].(string); ok {
		return source, nil
	}

	return "", fmt.Errorf("%w: %v", ErrFailedToGetPageSource, response)
}

// Close closes the current window.
func (d *WebDriver) Close(ctx context.Context) error {
	_, err := d.Execute(ctx, "closeWindow", nil)

	return err
}

// GetWindowHandle returns the handle of the current window.
func (d *WebDriver) GetWindowHandle(ctx context.Context) (string, error) {
	response, err := d.Execute(ctx, "getWindowHandle", nil)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGetWindowHandle, err)
	}
	if handle, ok := response["value"].(string); ok {
		return handle, nil
	}

	return "", fmt.Errorf("%w: %v", ErrFailedToGetWindowHandle, response)
}

// GetWindowHandles returns the handles of all windows.
func (d *WebDriver) GetWindowHandles(ctx context.Context) ([]string, error) {
	response, err := d.Execute(ctx, "getWindowHandles", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetWindowHandles, err)
	}
	if handles, ok := response["value"].([]interface{}); ok {
		result := make([]string, len(handles))
		for i, handle := range handles {
			if str, ok := handle.(string); ok {
				result[i] = str
			}
		}

		return result, nil
	}

	return nil, fmt.Errorf("%w: %v", ErrFailedToGetWindowHandles, response)
}

// SwitchToWindow switches to the window with the given handle.
func (d *WebDriver) SwitchToWindow(ctx context.Context, handle string) error {
	_, err := d.Execute(ctx, "switchToWindow", map[string]interface{}{
		"handle": handle,
	})

	return err
}

// MaximizeWindow maximizes the current window.
func (d *WebDriver) MaximizeWindow(ctx context.Context) error {
	_, err := d.Execute(ctx, "maximizeWindow", nil)

	return err
}

// MinimizeWindow minimizes the current window.
func (d *WebDriver) MinimizeWindow(ctx context.Context) error {
	_, err := d.Execute(ctx, "minimizeWindow", nil)

	return err
}

// FullscreenWindow makes the current window fullscreen.
func (d *WebDriver) FullscreenWindow(ctx context.Context) error {
	_, err := d.Execute(ctx, "fullscreenWindow", nil)

	return err
}

// Back navigates to the previous page in the browser history.
func (d *WebDriver) Back(ctx context.Context) error {
	_, err := d.Execute(ctx, "back", nil)

	return err
}

// Forward navigates to the next page in the browser history.
func (d *WebDriver) Forward(ctx context.Context) error {
	_, err := d.Execute(ctx, "forward", nil)

	return err
}

// Refresh refreshes the current page.
func (d *WebDriver) Refresh(ctx context.Context) error {
	_, err := d.Execute(ctx, "refresh", nil)

	return err
}

// GetCookies returns all cookies.
func (d *WebDriver) GetCookies(ctx context.Context) ([]selenium.Cookie, error) {
	response, err := d.Execute(ctx, "getAllCookies", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetCookies, err)
	}
	if cookies, ok := response["value"].([]interface{}); ok {
		result := make([]selenium.Cookie, len(cookies))
		for i, cookie := range cookies {
			if cookieMap, ok := cookie.(map[string]interface{}); ok {
				result[i] = mapToCookie(cookieMap)
			}
		}

		return result, nil
	}

	return nil, fmt.Errorf("%w: %v", ErrFailedToGetCookies, response)
}

// GetCookie returns the cookie with the given name.
func (d *WebDriver) GetCookie(ctx context.Context, name string) (*selenium.Cookie, error) {
	response, err := d.Execute(ctx, "getCookie", map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetCookie, err)
	}
	if cookie, ok := response["value"].(map[string]interface{}); ok {
		result := mapToCookie(cookie)

		return &result, nil
	}

	return nil, fmt.Errorf("%w: %v", ErrFailedToGetCookie, response)
}

// AddCookie adds a cookie.
func (d *WebDriver) AddCookie(ctx context.Context, cookie *selenium.Cookie) error {
	_, err := d.Execute(ctx, "addCookie", map[string]interface{}{
		"cookie": cookie,
	})

	return err
}

// DeleteCookie deletes the cookie with the given name.
func (d *WebDriver) DeleteCookie(ctx context.Context, name string) error {
	_, err := d.Execute(ctx, "deleteCookie", map[string]interface{}{
		"name": name,
	})

	return err
}

// DeleteAllCookies deletes all cookies.
func (d *WebDriver) DeleteAllCookies(ctx context.Context) error {
	_, err := d.Execute(ctx, "deleteAllCookies", nil)

	return err
}

// ExecuteScript executes JavaScript in the context of the currently selected frame or window.
func (d *WebDriver) ExecuteScript(ctx context.Context, script string, args []interface{}) (interface{}, error) {
	response, err := d.Execute(ctx, "executeScript", map[string]interface{}{
		"script": script,
		"args":   args,
	})
	if err != nil {
		return nil, err
	}

	return response["value"], nil
}

// ExecuteAsyncScript executes JavaScript asynchronously in the context of the currently selected frame or window.
func (d *WebDriver) ExecuteAsyncScript(ctx context.Context, script string, args []interface{}) (interface{}, error) {
	response, err := d.Execute(ctx, "executeAsyncScript", map[string]interface{}{
		"script": script,
		"args":   args,
	})
	if err != nil {
		return nil, err
	}

	return response["value"], nil
}

// Screenshot takes a screenshot of the current page.
func (d *WebDriver) Screenshot(ctx context.Context) ([]byte, error) {
	response, err := d.Execute(ctx, "screenshot", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetScreenshot, err)
	}
	if screenshot, ok := response["value"].(string); ok {
		return base64.StdEncoding.DecodeString(screenshot)
	}

	return nil, fmt.Errorf("%w: %v", ErrFailedToGetScreenshot, response)
}

// SetImplicitWaitTimeout sets the implicit wait timeout.
func (d *WebDriver) SetImplicitWaitTimeout(ctx context.Context, timeout int) error {
	_, err := d.Execute(ctx, "setTimeouts", map[string]interface{}{
		"implicit": timeout,
	})

	return err
}

// SetScriptTimeout sets the script timeout.
func (d *WebDriver) SetScriptTimeout(ctx context.Context, timeout int) error {
	_, err := d.Execute(ctx, "setTimeouts", map[string]interface{}{
		"script": timeout,
	})

	return err
}

// SetPageLoadTimeout sets the page load timeout.
func (d *WebDriver) SetPageLoadTimeout(ctx context.Context, timeout int) error {
	_, err := d.Execute(ctx, "setTimeouts", map[string]interface{}{
		"pageLoad": timeout,
	})

	return err
}

// mapToCookie converts a map to a Cookie struct.
func mapToCookie(m map[string]interface{}) selenium.Cookie {
	cookie := selenium.Cookie{
		Name:  m["name"].(string),
		Value: m["value"].(string),
	}
	if path, ok := m["path"].(string); ok {
		cookie.Path = path
	}
	if domain, ok := m["domain"].(string); ok {
		cookie.Domain = domain
	}
	if secure, ok := m["secure"].(bool); ok {
		cookie.Secure = secure
	}
	if httpOnly, ok := m["httpOnly"].(bool); ok {
		cookie.HTTPOnly = httpOnly
	}
	if expiry, ok := m["expiry"].(float64); ok {
		cookie.Expiry = expiry
	}
	if sameSite, ok := m["sameSite"].(string); ok {
		cookie.SameSite = sameSite
	}

	return cookie
}

// New creates a new browser session with the given capabilities.
func New(ctx context.Context, conn *connection.RemoteConnection, caps selenium.Convertible) (*WebDriver, error) {
	if conn == nil {
		return nil, errors.New("connection cannot be nil")
	}

	driver := &WebDriver{
		conn:         conn,
		capabilities: caps,
	}

	session, err := driver.newSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateSession, err)
	}

	driver.sessionID = session.SessionID
	driver.capabilities = session.Capabilities

	return driver, nil
}
