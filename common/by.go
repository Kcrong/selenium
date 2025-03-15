package common

import "sync"

// By represents a mechanism to locate elements within a document.
// Example usage:
//
//	element := driver.FindElement(By.ID, "myElement")
//	element := driver.FindElement(By.XPATH, "//html/body/div")
//	element := driver.FindElement(By.LinkText, "myLink")
//	element := driver.FindElement(By.PartialLinkText, "my")
//	element := driver.FindElement(By.Name, "myElement")
//	element := driver.FindElement(By.TagName, "div")
//	element := driver.FindElement(By.ClassName, "myElement")
//	element := driver.FindElement(By.CSSSelector, "div.myElement")
type By struct {
	// Standard locator strategies
	ID              string
	XPath           string
	LinkText        string
	PartialLinkText string
	Name            string
	TagName         string
	ClassName       string
	CSSSelector     string

	customFindersMu sync.RWMutex
	customFinders   map[string]string
}

// NewBy creates a new By instance with initialized values
func NewBy() *By {
	return &By{
		ID:              "id",
		XPath:           "xpath",
		LinkText:        "link text",
		PartialLinkText: "partial link text",
		Name:            "name",
		TagName:         "tag name",
		ClassName:       "class name",
		CSSSelector:     "css selector",
		customFinders:   make(map[string]string),
	}
}

// RegisterCustomFinder registers a custom finder strategy
func (b *By) RegisterCustomFinder(name, strategy string) {
	b.customFindersMu.Lock()
	defer b.customFindersMu.Unlock()
	b.customFinders[name] = strategy
}

// GetFinder returns the strategy for the given finder name
func (b *By) GetFinder(name string) string {
	b.customFindersMu.RLock()
	defer b.customFindersMu.RUnlock()

	if strategy, ok := b.customFinders[name]; ok {
		return strategy
	}

	switch name {
	case "id":
		return b.ID
	case "xpath":
		return b.XPath
	case "link text":
		return b.LinkText
	case "partial link text":
		return b.PartialLinkText
	case "name":
		return b.Name
	case "tag name":
		return b.TagName
	case "class name":
		return b.ClassName
	case "css selector":
		return b.CSSSelector
	default:
		return ""
	}
}

// ClearCustomFinders removes all registered custom finders
func (b *By) ClearCustomFinders() {
	b.customFindersMu.Lock()
	defer b.customFindersMu.Unlock()
	b.customFinders = make(map[string]string)
}
