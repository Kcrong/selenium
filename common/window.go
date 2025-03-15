package common

// WindowType represents a type of browser window
type WindowType string

const (
	// TabWindow represents a browser tab
	TabWindow WindowType = "tab"
	// NormalWindow represents a browser window
	NormalWindow WindowType = "window"
)
