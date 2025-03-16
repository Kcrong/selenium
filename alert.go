package selenium

import (
	"context"

	"github.com/Kcrong/selenium/remote/command"
)

// Alert represents a JavaScript alert, confirm, or prompt dialog.
type Alert struct {
	driver WebDriver
}

// NewAlert creates a new Alert instance.
func NewAlert(driver WebDriver) *Alert {
	return &Alert{
		driver: driver,
	}
}

// Text gets the text of the Alert.
//
// Example usage:
//
//	alert := NewAlert(driver)
//	text, err := alert.Text(ctx) // Get the text of an alert dialog.
func (a *Alert) Text(ctx context.Context) (string, error) {
	result, err := a.driver.Execute(ctx, command.W3CGetAlertText, nil)
	if err != nil {
		return "", err
	}
	return result["value"].(string), nil
}

// Dismiss dismisses the alert available.
//
// Example usage:
//
//	alert := NewAlert(driver)
//	err := alert.Dismiss(ctx) // Dismiss an alert dialog.
func (a *Alert) Dismiss(ctx context.Context) error {
	_, err := a.driver.Execute(ctx, command.W3CDismissAlert, nil)
	return err
}

// Accept accepts the alert available.
//
// Example usage:
//
//	alert := NewAlert(driver)
//	err := alert.Accept(ctx) // Confirm an alert dialog.
func (a *Alert) Accept(ctx context.Context) error {
	_, err := a.driver.Execute(ctx, command.W3CAcceptAlert, nil)
	return err
}

// SendKeys sends keys to the Alert.
//
// Example usage:
//
//	alert := NewAlert(driver)
//	err := alert.SendKeys(ctx, "Hello, World!") // Send keys to a prompt dialog.
func (a *Alert) SendKeys(ctx context.Context, keysToSend string) error {
	params := map[string]interface{}{
		"value": []string{keysToSend},
		"text":  keysToSend,
	}
	_, err := a.driver.Execute(ctx, command.W3CSetAlertValue, params)
	return err
}
