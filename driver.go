package selenium

import (
	"github.com/Kcrong/selenium/remote/command"
)

// WebDriver interface defines the basic operations that all WebDriver implementations must support.
type WebDriver interface {
	// Execute executes a WebDriver command
	Execute(cmd command.Command, params map[string]interface{}) (map[string]interface{}, error)
	// Quit closes the browser and ends the session
	Quit() error
	// GetCapabilities returns the current session's capabilities
	GetCapabilities() (Capabilities, error)
	// GetSessionID returns the current session ID
	GetSessionID() string
}
