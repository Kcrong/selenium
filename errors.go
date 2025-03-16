package selenigo

import (
	"fmt"
	"strings"
)

const (
	// SupportMsg is the base message for error documentation
	SupportMsg = "For documentation on this error, please visit:"
	// ErrorURL is the base URL for error documentation
	ErrorURL = "https://www.selenium.dev/documentation/webdriver/troubleshooting/errors"
)

// WebDriverError represents a base webdriver error
type WebDriverError struct {
	Message    string
	Screen     string
	Stacktrace []string
}

// Error implements the error interface
func (e *WebDriverError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Message: %s\n", e.Message))
	if e.Screen != "" {
		sb.WriteString("Screenshot: available via screen\n")
	}
	if len(e.Stacktrace) > 0 {
		sb.WriteString(fmt.Sprintf("Stacktrace:\n%s", strings.Join(e.Stacktrace, "\n")))
	}
	return sb.String()
}

// NewWebDriverError creates a new WebDriverError
func NewWebDriverError(msg string, screen string, stacktrace []string) *WebDriverError {
	return &WebDriverError{
		Message:    msg,
		Screen:     screen,
		Stacktrace: stacktrace,
	}
}

// InvalidSwitchToTargetError is thrown when frame or window target to be switched doesn't exist
type InvalidSwitchToTargetError struct {
	*WebDriverError
}

// NoSuchFrameError is thrown when frame target to be switched doesn't exist
type NoSuchFrameError struct {
	*InvalidSwitchToTargetError
}

// NoSuchWindowError is thrown when window target to be switched doesn't exist
type NoSuchWindowError struct {
	*InvalidSwitchToTargetError
}

// NoSuchElementError is thrown when element could not be found
type NoSuchElementError struct {
	*WebDriverError
}

// NewNoSuchElementError creates a new NoSuchElementError
func NewNoSuchElementError(msg string, screen string, stacktrace []string) *NoSuchElementError {
	withSupport := fmt.Sprintf("%s; %s %s#no-such-element-exception", msg, SupportMsg, ErrorURL)
	return &NoSuchElementError{
		WebDriverError: NewWebDriverError(withSupport, screen, stacktrace),
	}
}

// NoSuchAttributeError is thrown when the attribute of element could not be found
type NoSuchAttributeError struct {
	*WebDriverError
}

// NoSuchShadowRootError is thrown when trying to access the shadow root of an element when it does not have a shadow root attached
type NoSuchShadowRootError struct {
	*WebDriverError
}

// StaleElementReferenceError is thrown when a reference to an element is now "stale"
type StaleElementReferenceError struct {
	*WebDriverError
}

// NewStaleElementReferenceError creates a new StaleElementReferenceError
func NewStaleElementReferenceError(msg string, screen string, stacktrace []string) *StaleElementReferenceError {
	withSupport := fmt.Sprintf("%s; %s %s#stale-element-reference-exception", msg, SupportMsg, ErrorURL)
	return &StaleElementReferenceError{
		WebDriverError: NewWebDriverError(withSupport, screen, stacktrace),
	}
}

// InvalidElementStateError is thrown when a command could not be completed because the element is in an invalid state
type InvalidElementStateError struct {
	*WebDriverError
}

// UnexpectedAlertPresentError is thrown when an unexpected alert has appeared
type UnexpectedAlertPresentError struct {
	*WebDriverError
	AlertText string
}

// Error implements the error interface for UnexpectedAlertPresentError
func (e *UnexpectedAlertPresentError) Error() string {
	return fmt.Sprintf("Alert Text: %s\n%s", e.AlertText, e.WebDriverError.Error())
}

// NoAlertPresentError is thrown when switching to no presented alert
type NoAlertPresentError struct {
	*WebDriverError
}

// ElementNotVisibleError is thrown when an element is present on the DOM, but it is not visible
type ElementNotVisibleError struct {
	*InvalidElementStateError
}

// ElementNotInteractableError is thrown when an element is present in the DOM but interactions with that element will hit another element
type ElementNotInteractableError struct {
	*InvalidElementStateError
}

// ElementNotSelectableError is thrown when trying to select an unselectable element
type ElementNotSelectableError struct {
	*InvalidElementStateError
}

// InvalidCookieDomainError is thrown when attempting to add a cookie under a different domain than the current URL
type InvalidCookieDomainError struct {
	*WebDriverError
}

// UnableToSetCookieError is thrown when a driver fails to set a cookie
type UnableToSetCookieError struct {
	*WebDriverError
}

// TimeoutError is thrown when a command does not complete in enough time
type TimeoutError struct {
	*WebDriverError
}

// MoveTargetOutOfBoundsError is thrown when the target provided to the ActionsChains move() method is invalid
type MoveTargetOutOfBoundsError struct {
	*WebDriverError
}

// UnexpectedTagNameError is thrown when a support class did not get an expected web element
type UnexpectedTagNameError struct {
	*WebDriverError
}

// InvalidSelectorError is thrown when the selector which is used to find an element does not return a WebElement
type InvalidSelectorError struct {
	*WebDriverError
}

// NewInvalidSelectorError creates a new InvalidSelectorError
func NewInvalidSelectorError(msg string, screen string, stacktrace []string) *InvalidSelectorError {
	withSupport := fmt.Sprintf("%s; %s %s#invalid-selector-exception", msg, SupportMsg, ErrorURL)
	return &InvalidSelectorError{
		WebDriverError: NewWebDriverError(withSupport, screen, stacktrace),
	}
}

// ImeNotAvailableError is thrown when IME support is not available
type ImeNotAvailableError struct {
	*WebDriverError
}

// ImeActivationFailedError is thrown when activating an IME engine has failed
type ImeActivationFailedError struct {
	*WebDriverError
}

// InvalidArgumentError is thrown when the arguments passed to a command are either invalid or malformed
type InvalidArgumentError struct {
	*WebDriverError
}

// JavascriptError is thrown when an error occurred while executing JavaScript supplied by the user
type JavascriptError struct {
	*WebDriverError
}

// NoSuchCookieError is thrown when no cookie matching the given path name was found
type NoSuchCookieError struct {
	*WebDriverError
}

// ScreenshotError is thrown when a screen capture was made impossible
type ScreenshotError struct {
	*WebDriverError
}

// ElementClickInterceptedError is thrown when the element click was intercepted by another element
type ElementClickInterceptedError struct {
	*WebDriverError
}

// InsecureCertificateError is thrown when there is an insecure certificate
type InsecureCertificateError struct {
	*WebDriverError
}

// InvalidCoordinatesError is thrown when the coordinates provided are invalid
type InvalidCoordinatesError struct {
	*WebDriverError
}

// InvalidSessionIDError is thrown when the session ID is invalid
type InvalidSessionIDError struct {
	*WebDriverError
}

// SessionNotCreatedError is thrown when a new session could not be created
type SessionNotCreatedError struct {
	*WebDriverError
}

// UnknownMethodError is thrown when the requested command matches no known command
type UnknownMethodError struct {
	*WebDriverError
}

// NoSuchDriverError is thrown when no driver matching the requirements could be found
type NoSuchDriverError struct {
	*WebDriverError
}

// NewNoSuchDriverError creates a new NoSuchDriverError
func NewNoSuchDriverError(msg string, screen string, stacktrace []string) *NoSuchDriverError {
	withSupport := fmt.Sprintf("%s; %s %s#no-such-driver-exception", msg, SupportMsg, ErrorURL)
	return &NoSuchDriverError{
		WebDriverError: NewWebDriverError(withSupport, screen, stacktrace),
	}
}

// DetachedShadowRootError is thrown when the shadow root is detached from the element
type DetachedShadowRootError struct {
	*WebDriverError
}
