package webelement

// WebElement represents a DOM element.
type WebElement interface {
	// GetID returns the internal element ID used by WebDriver.
	GetID() string
}

// webElement represents a remote DOM element.
type webElement struct {
	id string
}

// NewElement creates a new webElement with the given ID.
func NewElement(id string) WebElement {
	return &webElement{
		id: id,
	}
}

// GetID returns the internal element ID used by WebDriver.
func (e *webElement) GetID() string {
	return e.id
}
