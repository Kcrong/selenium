package actions

import (
	"time"
)

// KeyInput represents a keyboard input device.
type KeyInput struct {
	name    string
	actions []Action
}

// NewKeyInput creates a new keyboard input device.
func NewKeyInput(name string) *KeyInput {
	return &KeyInput{
		name:    name,
		actions: make([]Action, 0),
	}
}

// GetName returns the name of the input device.
func (k *KeyInput) GetName() string {
	return k.name
}

// GetType returns the type of the input device.
func (k *KeyInput) GetType() string {
	return "key"
}

// ClearActions clears all stored actions.
func (k *KeyInput) ClearActions() {
	k.actions = make([]Action, 0)
}

// CreatePause creates a pause action.
func (k *KeyInput) CreatePause(duration time.Duration) Action {
	return NewPauseAction(duration)
}

// GetActions returns all stored actions.
func (k *KeyInput) GetActions() []Action {
	return k.actions
}

// AddAction adds an action to the device.
func (k *KeyInput) AddAction(action Action) {
	k.actions = append(k.actions, action)
}

// KeyAction represents a keyboard action.
type KeyAction struct {
	BaseAction
	Value string
}

// NewKeyAction creates a new keyboard action.
func NewKeyAction(actionType string, value string) *KeyAction {
	return &KeyAction{
		BaseAction: BaseAction{
			Type:     actionType,
			Duration: 0,
		},
		Value: value,
	}
}

// Encode encodes the key action into a format that can be sent to the WebDriver.
func (a *KeyAction) Encode() map[string]interface{} {
	encoded := map[string]interface{}{
		"type":  a.Type,
		"value": a.Value,
	}
	if a.Duration > 0 {
		encoded["duration"] = a.Duration
	}
	return encoded
}

// KeyDown creates a key down action.
func (k *KeyInput) KeyDown(value string) {
	k.actions = append(k.actions, NewKeyAction("keyDown", value))
}

// KeyUp creates a key up action.
func (k *KeyInput) KeyUp(value string) {
	k.actions = append(k.actions, NewKeyAction("keyUp", value))
}

// SendKeys creates a sequence of key down and key up actions.
func (k *KeyInput) SendKeys(value string) {
	for _, char := range value {
		k.KeyDown(string(char))
		k.KeyUp(string(char))
	}
}
