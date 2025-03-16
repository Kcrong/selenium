package ie

import (
	"errors"
	"fmt"

	"github.com/Kcrong/selenigo"
)

const (
	// OptionsKey is the capability key for IE options.
	OptionsKey = "se:ieOptions"
	// SwitchesKey is the capability key for IE command line switches.
	SwitchesKey = "ie.browserCommandLineSwitches"
	// BrowserName is the name of the Internet Explorer browser.
	BrowserName = "internet explorer"
)

// ElementScrollBehavior represents the scroll behavior for elements.
type ElementScrollBehavior int

const (
	// ScrollToTop scrolls elements to the top.
	ScrollToTop ElementScrollBehavior = iota
	// ScrollToBottom scrolls elements to the bottom.
	ScrollToBottom
)

// ErrInvalidScrollBehavior is returned when an invalid scroll behavior is provided.
var ErrInvalidScrollBehavior = errors.New("invalid element scroll behavior")

// Options contains the options for Internet Explorer browser.
type Options struct {
	options map[string]interface{}
}

// NewOptions creates a new IE options instance.
func NewOptions() *Options {
	return &Options{
		options: make(map[string]interface{}),
	}
}

// SetBrowserAttachTimeout sets the timeout for browser attach.
func (o *Options) SetBrowserAttachTimeout(timeout int) {
	o.options["browserAttachTimeout"] = timeout
}

// GetBrowserAttachTimeout returns the browser attach timeout.
func (o *Options) GetBrowserAttachTimeout() int {
	if timeout, ok := o.options["browserAttachTimeout"].(int); ok {
		return timeout
	}

	return 0
}

// SetElementScrollBehavior sets the scroll behavior for elements.
func (o *Options) SetElementScrollBehavior(behavior ElementScrollBehavior) error {
	if behavior != ScrollToTop && behavior != ScrollToBottom {
		return fmt.Errorf("%w: %d", ErrInvalidScrollBehavior, behavior)
	}

	o.options["elementScrollBehavior"] = behavior

	return nil
}

// GetElementScrollBehavior returns the scroll behavior for elements.
func (o *Options) GetElementScrollBehavior() ElementScrollBehavior {
	if behavior, ok := o.options["elementScrollBehavior"].(ElementScrollBehavior); ok {
		return behavior
	}

	return ScrollToTop
}

// SetEnsureCleanSession sets whether to ensure a clean session.
func (o *Options) SetEnsureCleanSession(ensure bool) {
	o.options["ie.ensureCleanSession"] = ensure
}

// GetEnsureCleanSession returns whether to ensure a clean session.
func (o *Options) GetEnsureCleanSession() bool {
	if ensure, ok := o.options["ie.ensureCleanSession"].(bool); ok {
		return ensure
	}

	return false
}

// SetFileUploadDialogTimeout sets the timeout for file upload dialog.
func (o *Options) SetFileUploadDialogTimeout(timeout int) {
	o.options["ie.fileUploadDialogTimeout"] = timeout
}

// GetFileUploadDialogTimeout returns the file upload dialog timeout.
func (o *Options) GetFileUploadDialogTimeout() int {
	if timeout, ok := o.options["ie.fileUploadDialogTimeout"].(int); ok {
		return timeout
	}

	return 0
}

// SetForceCreateProcessAPI sets whether to force create process API.
func (o *Options) SetForceCreateProcessAPI(force bool) {
	o.options["ie.forceCreateProcessApi"] = force
}

// GetForceCreateProcessAPI returns whether to force create process API.
func (o *Options) GetForceCreateProcessAPI() bool {
	if force, ok := o.options["ie.forceCreateProcessApi"].(bool); ok {
		return force
	}

	return false
}

// SetForceShellWindowsAPI sets whether to force shell windows API.
func (o *Options) SetForceShellWindowsAPI(force bool) {
	o.options["ie.forceShellWindowsApi"] = force
}

// GetForceShellWindowsAPI returns whether to force shell windows API.
func (o *Options) GetForceShellWindowsAPI() bool {
	if force, ok := o.options["ie.forceShellWindowsApi"].(bool); ok {
		return force
	}

	return false
}

// SetFullPageScreenshot sets whether to enable full page screenshot.
func (o *Options) SetFullPageScreenshot(enable bool) {
	o.options["ie.enableFullPageScreenshot"] = enable
}

// GetFullPageScreenshot returns whether full page screenshot is enabled.
func (o *Options) GetFullPageScreenshot() bool {
	if enable, ok := o.options["ie.enableFullPageScreenshot"].(bool); ok {
		return enable
	}

	return false
}

// SetIgnoreProtectedModeSettings sets whether to ignore protected mode settings.
func (o *Options) SetIgnoreProtectedModeSettings(ignore bool) {
	o.options["ignoreProtectedModeSettings"] = ignore
}

// GetIgnoreProtectedModeSettings returns whether to ignore protected mode settings.
func (o *Options) GetIgnoreProtectedModeSettings() bool {
	if ignore, ok := o.options["ignoreProtectedModeSettings"].(bool); ok {
		return ignore
	}

	return false
}

// SetIgnoreZoomLevel sets whether to ignore zoom level.
func (o *Options) SetIgnoreZoomLevel(ignore bool) {
	o.options["ignoreZoomSetting"] = ignore
}

// GetIgnoreZoomLevel returns whether to ignore zoom level.
func (o *Options) GetIgnoreZoomLevel() bool {
	if ignore, ok := o.options["ignoreZoomSetting"].(bool); ok {
		return ignore
	}

	return false
}

// SetInitialBrowserURL sets the initial browser URL.
func (o *Options) SetInitialBrowserURL(url string) {
	o.options["initialBrowserUrl"] = url
}

// GetInitialBrowserURL returns the initial browser URL.
func (o *Options) GetInitialBrowserURL() string {
	if url, ok := o.options["initialBrowserUrl"].(string); ok {
		return url
	}

	return ""
}

// SetNativeEvents sets whether to use native events.
func (o *Options) SetNativeEvents(use bool) {
	o.options["nativeEvents"] = use
}

// GetNativeEvents returns whether to use native events.
func (o *Options) GetNativeEvents() bool {
	if use, ok := o.options["nativeEvents"].(bool); ok {
		return use
	}

	return false
}

// SetPersistentHover sets whether to enable persistent hover.
func (o *Options) SetPersistentHover(enable bool) {
	o.options["enablePersistentHover"] = enable
}

// GetPersistentHover returns whether persistent hover is enabled.
func (o *Options) GetPersistentHover() bool {
	if enable, ok := o.options["enablePersistentHover"].(bool); ok {
		return enable
	}

	return false
}

// SetRequireWindowFocus sets whether to require window focus.
func (o *Options) SetRequireWindowFocus(require bool) {
	o.options["requireWindowFocus"] = require
}

// GetRequireWindowFocus returns whether window focus is required.
func (o *Options) GetRequireWindowFocus() bool {
	if require, ok := o.options["requireWindowFocus"].(bool); ok {
		return require
	}

	return false
}

// SetUsePerProcessProxy sets whether to use per process proxy.
func (o *Options) SetUsePerProcessProxy(use bool) {
	o.options["ie.usePerProcessProxy"] = use
}

// GetUsePerProcessProxy returns whether to use per process proxy.
func (o *Options) GetUsePerProcessProxy() bool {
	if use, ok := o.options["ie.usePerProcessProxy"].(bool); ok {
		return use
	}

	return false
}

// SetUseLegacyFileUploadDialogHandling sets whether to use legacy file upload dialog handling.
func (o *Options) SetUseLegacyFileUploadDialogHandling(use bool) {
	o.options["ie.useLegacyFileUploadDialogHandling"] = use
}

// GetUseLegacyFileUploadDialogHandling returns whether to use legacy file upload dialog handling.
func (o *Options) GetUseLegacyFileUploadDialogHandling() bool {
	if use, ok := o.options["ie.useLegacyFileUploadDialogHandling"].(bool); ok {
		return use
	}

	return false
}

// SetAttachToEdgeChrome sets whether to attach to Edge Chrome.
func (o *Options) SetAttachToEdgeChrome(attach bool) {
	o.options["ie.edgechromium"] = attach
}

// GetAttachToEdgeChrome returns whether to attach to Edge Chrome.
func (o *Options) GetAttachToEdgeChrome() bool {
	if attach, ok := o.options["ie.edgechromium"].(bool); ok {
		return attach
	}

	return false
}

// SetEdgeExecutablePath sets the path to Edge executable.
func (o *Options) SetEdgeExecutablePath(path string) {
	o.options["ie.edgepath"] = path
}

// GetEdgeExecutablePath returns the path to Edge executable.
func (o *Options) GetEdgeExecutablePath() string {
	if path, ok := o.options["ie.edgepath"].(string); ok {
		return path
	}

	return ""
}

// SetIgnoreProcessMatch sets whether to ignore process match.
func (o *Options) SetIgnoreProcessMatch(ignore bool) {
	o.options["ie.ignoreprocessmatch"] = ignore
}

// GetIgnoreProcessMatch returns whether to ignore process match.
func (o *Options) GetIgnoreProcessMatch() bool {
	if ignore, ok := o.options["ie.ignoreprocessmatch"].(bool); ok {
		return ignore
	}

	return false
}

// ToCapabilities converts the options to a capabilities map.
func (o *Options) ToCapabilities() map[string]interface{} {
	caps := selenigo.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	if len(o.options) > 0 {
		caps.SetBrowserOptions(OptionsKey, o.options)
	}

	return caps.ToCapabilities()
}

// DefaultCapabilities returns the default capabilities for IE.
func (o *Options) DefaultCapabilities() map[string]interface{} {
	caps := selenigo.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	return caps.ToCapabilities()
}
