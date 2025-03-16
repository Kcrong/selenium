package bidi

// ConsoleType represents the type of console message.
type ConsoleType string

const (
	// ConsoleAll represents all console message types.
	ConsoleAll ConsoleType = "all"
	// ConsoleLog represents log level console messages.
	ConsoleLog ConsoleType = "log"
	// ConsoleError represents error level console messages.
	ConsoleError ConsoleType = "error"
)

// ConsoleMessage represents a console message from the browser.
type ConsoleMessage struct {
	Type    ConsoleType `json:"type"`
	Message string      `json:"message"`
}
