package actions

import (
	"context"
	"time"

	"github.com/Kcrong/selenium"
	"github.com/Kcrong/selenium/remote/command"
)

// ActionBuilder provides a way to create action sequences.
type ActionBuilder struct {
	driver   selenium.WebDriver
	devices  []InputDevice
	duration time.Duration
}

// NewActionBuilder creates a new ActionBuilder instance.
func NewActionBuilder(
	driver selenium.WebDriver, mouse *PointerInput, keyboard *KeyInput, wheel *WheelInput, duration time.Duration,
) *ActionBuilder {
	devices := make([]InputDevice, 0)

	if mouse == nil {
		mouse = NewPointerInput("mouse", "default mouse", duration)
	}
	if keyboard == nil {
		keyboard = NewKeyInput("default keyboard")
	}
	if wheel == nil {
		wheel = NewWheelInput("default wheel")
	}

	devices = append(devices, mouse, keyboard, wheel)

	return &ActionBuilder{
		driver:   driver,
		devices:  devices,
		duration: duration,
	}
}

// ClearActions clears all actions that are already stored on the remote end.
func (a *ActionBuilder) ClearActions(ctx context.Context) error {
	_, err := a.driver.Execute(ctx, command.W3CClearActions, nil)
	return err
}

// Perform performs all stored actions.
func (a *ActionBuilder) Perform(ctx context.Context) error {
	params := make(map[string]interface{})
	actions := make([]map[string]interface{}, 0)

	for _, device := range a.devices {
		deviceActions := make([]map[string]interface{}, 0)
		for _, action := range device.GetActions() {
			deviceActions = append(deviceActions, action.Encode())
		}

		if len(deviceActions) > 0 {
			actions = append(actions, map[string]interface{}{
				"type":    device.GetType(),
				"id":      device.GetName(),
				"actions": deviceActions,
			})
		}
	}

	params["actions"] = actions
	_, err := a.driver.Execute(ctx, command.W3CActions, params)
	return err
}

// GetDevices returns all input devices.
func (a *ActionBuilder) GetDevices() []InputDevice {
	return a.devices
}

// AddAction adds an action to a specific device.
func (a *ActionBuilder) AddAction(deviceType string, action Action) {
	for _, device := range a.devices {
		if device.GetType() == deviceType {
			device.AddAction(action)
			break
		}
	}
}

// AddPause adds a pause action to all devices.
func (a *ActionBuilder) AddPause(duration time.Duration) {
	for _, device := range a.devices {
		device.AddAction(device.CreatePause(duration))
	}
}
