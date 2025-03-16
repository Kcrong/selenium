package webelement

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Kcrong/selenium"
	"github.com/Kcrong/selenium/remote/command"
	"github.com/Kcrong/selenium/remote/connection"
)

// webElement represents a remote DOM element.
type webElement struct {
	id      string
	conn    *connection.RemoteConnection
	session string
}

// NewElement creates a new webElement with the given ID.
func NewElement(id, session string, conn *connection.RemoteConnection) selenium.WebElement {
	return &webElement{
		id:      id,
		conn:    conn,
		session: session,
	}
}

// GetID returns the internal element ID used by WebDriver.
func (e *webElement) GetID() string {
	return e.id
}

// Click clicks on the element.
func (e *webElement) Click(ctx context.Context) error {
	_, err := e.conn.Execute(ctx, command.Command("clickElement"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
	})

	return err
}

// SendKeys types the given keys into the element.
func (e *webElement) SendKeys(ctx context.Context, keys string) error {
	_, err := e.conn.Execute(ctx, command.Command("sendKeysToElement"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
		"value":     []string{keys},
	})

	return err
}

// Clear clears the element's value.
func (e *webElement) Clear(ctx context.Context) error {
	_, err := e.conn.Execute(ctx, command.Command("clearElement"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
	})

	return err
}

// GetText returns the visible text of the element.
func (e *webElement) GetText(ctx context.Context) (string, error) {
	response, err := e.conn.Execute(ctx, command.Command("getElementText"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
	})
	if err != nil {
		return "", err
	}

	if text, ok := response["value"].(string); ok {
		return text, nil
	}

	return "", fmt.Errorf("failed to get element text: %v", response)
}

// GetAttribute returns the value of the given attribute.
func (e *webElement) GetAttribute(ctx context.Context, name string) (string, error) {
	response, err := e.conn.Execute(ctx, command.Command("getElementAttribute"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
		"name":      name,
	})
	if err != nil {
		return "", err
	}

	if value, ok := response["value"].(string); ok {
		return value, nil
	}

	return "", fmt.Errorf("failed to get element attribute: %v", response)
}

// GetProperty returns the value of the given property.
func (e *webElement) GetProperty(ctx context.Context, name string) (string, error) {
	response, err := e.conn.Execute(ctx, command.Command("getElementProperty"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
		"name":      name,
	})
	if err != nil {
		return "", err
	}

	if value, ok := response["value"].(string); ok {
		return value, nil
	}

	return "", fmt.Errorf("failed to get element property: %v", response)
}

// IsSelected returns whether the element is selected.
func (e *webElement) IsSelected(ctx context.Context) (bool, error) {
	response, err := e.conn.Execute(ctx, command.Command("isElementSelected"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
	})
	if err != nil {
		return false, err
	}

	if selected, ok := response["value"].(bool); ok {
		return selected, nil
	}

	return false, fmt.Errorf("failed to get element selected state: %v", response)
}

// IsEnabled returns whether the element is enabled.
func (e *webElement) IsEnabled(ctx context.Context) (bool, error) {
	response, err := e.conn.Execute(ctx, command.Command("isElementEnabled"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
	})
	if err != nil {
		return false, err
	}

	if enabled, ok := response["value"].(bool); ok {
		return enabled, nil
	}

	return false, fmt.Errorf("failed to get element enabled state: %v", response)
}

// IsDisplayed returns whether the element is displayed.
func (e *webElement) IsDisplayed(ctx context.Context) (bool, error) {
	response, err := e.conn.Execute(ctx, command.Command("isElementDisplayed"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
	})
	if err != nil {
		return false, err
	}

	if displayed, ok := response["value"].(bool); ok {
		return displayed, nil
	}

	return false, fmt.Errorf("failed to get element displayed state: %v", response)
}

// Screenshot takes a screenshot of the element.
func (e *webElement) Screenshot(ctx context.Context) ([]byte, error) {
	response, err := e.conn.Execute(ctx, command.Command("elementScreenshot"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
	})
	if err != nil {
		return nil, err
	}

	if screenshot, ok := response["value"].(string); ok {
		return base64.StdEncoding.DecodeString(screenshot)
	}

	return nil, fmt.Errorf("failed to get element screenshot: %v", response)
}

// FindElement finds a child element using the given locator.
func (e *webElement) FindElement(ctx context.Context, by *selenium.By, value string) (selenium.WebElement, error) {
	response, err := e.conn.Execute(ctx, command.Command("findChildElement"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
		"using":     by.GetFinder(value),
		"value":     value,
	})
	if err != nil {
		return nil, err
	}

	if element, ok := response["value"].(map[string]interface{}); ok {
		if elementID, ok := element["ELEMENT"].(string); ok {
			return NewElement(elementID, e.session, e.conn), nil
		}
	}

	return nil, fmt.Errorf("failed to find element: %v", response)
}

// FindElements finds child elements using the given locator.
func (e *webElement) FindElements(ctx context.Context, by *selenium.By, value string) ([]selenium.WebElement, error) {
	response, err := e.conn.Execute(ctx, command.Command("findChildElements"), map[string]interface{}{
		"sessionId": e.session,
		"id":        e.id,
		"using":     by.GetFinder(value),
		"value":     value,
	})
	if err != nil {
		return nil, err
	}

	if elements, ok := response["value"].([]interface{}); ok {
		result := make([]selenium.WebElement, len(elements))
		for i, element := range elements {
			if elementMap, ok := element.(map[string]interface{}); ok {
				if elementID, ok := elementMap["ELEMENT"].(string); ok {
					result[i] = NewElement(elementID, e.session, e.conn)

					continue
				}
			}

			return nil, fmt.Errorf("invalid element at index %d: %v", i, element)
		}

		return result, nil
	}

	return nil, fmt.Errorf("failed to find elements: %v", response)
}

// Ensure webElement implements selenium.WebElement.
var _ selenium.WebElement = (*webElement)(nil)
