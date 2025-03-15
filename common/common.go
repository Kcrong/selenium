package common

import (
	"github.com/Kcrong/selenium-go"
)

// Re-export all error types from webdriver package
type (
	WebDriverError               = selenium.WebDriverError
	InvalidSwitchToTargetError   = selenium.InvalidSwitchToTargetError
	NoSuchFrameError             = selenium.NoSuchFrameError
	NoSuchWindowError            = selenium.NoSuchWindowError
	NoSuchElementError           = selenium.NoSuchElementError
	NoSuchAttributeError         = selenium.NoSuchAttributeError
	NoSuchDriverError            = selenium.NoSuchDriverError
	NoSuchShadowRootError        = selenium.NoSuchShadowRootError
	StaleElementReferenceError   = selenium.StaleElementReferenceError
	InvalidElementStateError     = selenium.InvalidElementStateError
	UnexpectedAlertPresentError  = selenium.UnexpectedAlertPresentError
	NoAlertPresentError          = selenium.NoAlertPresentError
	ElementNotVisibleError       = selenium.ElementNotVisibleError
	ElementNotInteractableError  = selenium.ElementNotInteractableError
	ElementNotSelectableError    = selenium.ElementNotSelectableError
	InvalidCookieDomainError     = selenium.InvalidCookieDomainError
	UnableToSetCookieError       = selenium.UnableToSetCookieError
	TimeoutError                 = selenium.TimeoutError
	MoveTargetOutOfBoundsError   = selenium.MoveTargetOutOfBoundsError
	UnexpectedTagNameError       = selenium.UnexpectedTagNameError
	InvalidSelectorError         = selenium.InvalidSelectorError
	ImeNotAvailableError         = selenium.ImeNotAvailableError
	ImeActivationFailedError     = selenium.ImeActivationFailedError
	InvalidArgumentError         = selenium.InvalidArgumentError
	JavascriptError              = selenium.JavascriptError
	NoSuchCookieError            = selenium.NoSuchCookieError
	ScreenshotError              = selenium.ScreenshotError
	ElementClickInterceptedError = selenium.ElementClickInterceptedError
	InsecureCertificateError     = selenium.InsecureCertificateError
	InvalidCoordinatesError      = selenium.InvalidCoordinatesError
	InvalidSessionIDError        = selenium.InvalidSessionIDError
	SessionNotCreatedError       = selenium.SessionNotCreatedError
	UnknownMethodError           = selenium.UnknownMethodError
	DetachedShadowRootError      = selenium.DetachedShadowRootError
)

// Re-export constructor functions
var (
	NewWebDriverError             = selenium.NewWebDriverError
	NewNoSuchElementError         = selenium.NewNoSuchElementError
	NewStaleElementReferenceError = selenium.NewStaleElementReferenceError
	NewInvalidSelectorError       = selenium.NewInvalidSelectorError
	NewNoSuchDriverError          = selenium.NewNoSuchDriverError
)
