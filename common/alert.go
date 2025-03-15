package common

import (
	"github.com/Kcrong/selenium-go"
	"github.com/Kcrong/selenium-go/remote/command"
)

// Alert represents a JavaScript alert, confirm, or prompt dialog.
type Alert struct {
	driver selenium.WebDriver
}

// NewAlert creates a new Alert instance.
func NewAlert(driver selenium.WebDriver) *Alert {
	return &Alert{
		driver: driver,
	}
}

// Text gets the text of the Alert.
func (a *Alert) Text() (string, error) {
	result, err := a.driver.Execute(command.W3CGetAlertText, nil)
	if err != nil {
		return "", err
	}
	return result["value"].(string), nil
}

// Dismiss dismisses the alert available.
func (a *Alert) Dismiss() error {
	_, err := a.driver.Execute(command.W3CDismissAlert, nil)
	return err
}

// Accept accepts the alert available.
//
// Example usage:
//
//	alert := NewAlert(driver)
//	err := alert.Accept() // Confirm an alert dialog.
func (a *Alert) Accept() error {
	_, err := a.driver.Execute(command.W3CAcceptAlert, nil)
	return err
}

// SendKeys sends keys to the Alert.
func (a *Alert) SendKeys(keysToSend string) error {
	params := map[string]interface{}{
		"value": []string{keysToSend},
		"text":  keysToSend,
	}
	_, err := a.driver.Execute(command.W3CSetAlertValue, params)
	return err
}
