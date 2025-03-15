package actions

import (
	"time"
)

// InputDevice represents a virtual input device for WebDriver actions.
type InputDevice interface {
	// GetName returns the name of the input device
	GetName() string
	// GetType returns the type of the input device
	GetType() string
	// ClearActions clears all stored actions
	ClearActions()
	// CreatePause creates a pause action
	CreatePause(duration time.Duration) Action
	// GetActions returns all stored actions
	GetActions() []Action
	// AddAction adds an action to the device
	AddAction(action Action)
}

// Action represents a single action that can be performed by an input device.
type Action interface {
	// GetType returns the type of the action
	GetType() string
	// GetDuration returns the duration of the action
	GetDuration() time.Duration
	// Encode encodes the action into a format that can be sent to the WebDriver
	Encode() map[string]interface{}
}

// BaseAction provides common functionality for actions.
type BaseAction struct {
	Type     string
	Duration time.Duration
}

// GetType returns the type of the action.
func (a *BaseAction) GetType() string {
	return a.Type
}

// GetDuration returns the duration of the action
func (a *BaseAction) GetDuration() time.Duration {
	return a.Duration
}

// PauseAction represents a pause in the action sequence.
type PauseAction struct {
	BaseAction
}

// NewPauseAction creates a new pause action.
func NewPauseAction(duration time.Duration) *PauseAction {
	return &PauseAction{
		BaseAction: BaseAction{
			Type:     "pause",
			Duration: duration,
		},
	}
}

// Encode encodes the pause action into a format that can be sent to the WebDriver.
func (a *PauseAction) Encode() map[string]interface{} {
	return map[string]interface{}{
		"type":     a.Type,
		"duration": a.Duration,
	}
}
