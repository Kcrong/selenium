package firefox

import (
	"github.com/Kcrong/selenium"
)

const (
	// OptionsKey is the capability key for firefox options.
	OptionsKey = "moz:firefoxOptions"
	// BrowserName is the name of the Firefox browser.
	BrowserName = "firefox"
)

// Log represents Firefox logging options.
type Log struct {
	Level string
}

// ToCapabilities converts the log options to a capabilities map.
func (l *Log) ToCapabilities() map[string]interface{} {
	if l.Level == "" {
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		"log": map[string]interface{}{
			"level": l.Level,
		},
	}
}

// Options contains the options for Firefox browser.
type Options struct {
	binaryLocation string
	preferences    map[string]interface{}
	profile        string
	log            *Log
	arguments      []string
}

// NewOptions creates a new Firefox options instance.
func NewOptions() *Options {
	return &Options{
		binaryLocation: "",
		preferences:    make(map[string]interface{}),
		profile:        "",
		log:            &Log{Level: ""},
		arguments:      make([]string, 0),
	}
}

// SetBinaryLocation sets the path to Firefox binary.
func (o *Options) SetBinaryLocation(path string) {
	o.binaryLocation = path
}

// GetBinaryLocation returns the path to Firefox binary.
func (o *Options) GetBinaryLocation() string {
	return o.binaryLocation
}

// SetPreference sets a Firefox preference.
func (o *Options) SetPreference(name string, value interface{}) {
	o.preferences[name] = value
}

// GetPreferences returns all Firefox preferences.
func (o *Options) GetPreferences() map[string]interface{} {
	return o.preferences
}

// SetProfile sets the Firefox profile to use.
func (o *Options) SetProfile(profile string) {
	o.profile = profile
}

// GetProfile returns the Firefox profile being used.
func (o *Options) GetProfile() string {
	return o.profile
}

// AddArgument adds a command-line argument.
func (o *Options) AddArgument(arg string) {
	o.arguments = append(o.arguments, arg)
}

// SetLogLevel sets the logging level.
func (o *Options) SetLogLevel(level string) {
	o.log.Level = level
}

// ToCapabilities converts the options to a capabilities map.
func (o *Options) ToCapabilities() map[string]interface{} {
	caps := selenium.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	opts := make(map[string]interface{})

	if o.binaryLocation != "" {
		opts["binary"] = o.binaryLocation
	}

	prefs := o.preferences
	if len(prefs) > 0 {
		opts["prefs"] = prefs
	}

	if o.profile != "" {
		opts["profile"] = o.profile
	}

	args := o.arguments
	if len(args) > 0 {
		opts["args"] = args
	}

	// Add log options
	for k, v := range o.log.ToCapabilities() {
		opts[k] = v
	}

	if len(opts) > 0 {
		caps.SetBrowserOptions(OptionsKey, opts)
	}

	return caps.ToCapabilities()
}

// DefaultCapabilities returns the default capabilities for Firefox.
func (o *Options) DefaultCapabilities() map[string]interface{} {
	caps := selenium.NewCapabilities()
	caps.Capabilities.BrowserName = BrowserName

	return caps.ToCapabilities()
}
