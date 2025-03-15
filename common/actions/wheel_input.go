package actions

import (
	"time"

	"github.com/Kcrong/selenium-go/remote/webelement"
)

// ScrollOrigin represents the origin point for a scroll action.
type ScrollOrigin struct {
	Element  webelement.WebElement
	XOffset  int
	YOffset  int
	Viewport bool
}

// WheelInput represents a wheel input device.
type WheelInput struct {
	name    string
	actions []Action
}

// NewWheelInput creates a new wheel input device.
func NewWheelInput(name string) *WheelInput {
	return &WheelInput{
		name:    name,
		actions: make([]Action, 0),
	}
}

// GetName returns the name of the input device.
func (w *WheelInput) GetName() string {
	return w.name
}

// GetType returns the type of the input device.
func (w *WheelInput) GetType() string {
	return "wheel"
}

// ClearActions clears all stored actions.
func (w *WheelInput) ClearActions() {
	w.actions = make([]Action, 0)
}

// CreatePause creates a pause action.
func (w *WheelInput) CreatePause(duration time.Duration) Action {
	return NewPauseAction(duration)
}

// GetActions returns all stored actions.
func (w *WheelInput) GetActions() []Action {
	return w.actions
}

// AddAction adds an action to the device.
func (w *WheelInput) AddAction(action Action) {
	w.actions = append(w.actions, action)
}

// WheelAction represents a wheel action.
type WheelAction struct {
	BaseAction
	DeltaX int
	DeltaY int
	Origin interface{}
	X      int
	Y      int
}

// NewWheelAction creates a new wheel action.
func NewWheelAction(deltaX, deltaY int, origin interface{}, x, y int) *WheelAction {
	return &WheelAction{
		BaseAction: BaseAction{
			Type:     "wheel",
			Duration: 0,
		},
		DeltaX: deltaX,
		DeltaY: deltaY,
		Origin: origin,
		X:      x,
		Y:      y,
	}
}

// Encode encodes the wheel action into a format that can be sent to the WebDriver.
func (a *WheelAction) Encode() map[string]interface{} {
	encoded := map[string]interface{}{
		"type":   a.Type,
		"deltaX": a.DeltaX,
		"deltaY": a.DeltaY,
	}

	if a.Duration > 0 {
		encoded["duration"] = a.Duration
	}

	if a.Origin != nil {
		if element, ok := a.Origin.(webelement.WebElement); ok {
			encoded["origin"] = map[string]string{
				"element-6066-11e4-a52e-4f735466cecf": element.GetID(),
			}
		} else {
			encoded["origin"] = a.Origin
		}
	}

	if a.X != 0 {
		encoded["x"] = a.X
	}
	if a.Y != 0 {
		encoded["y"] = a.Y
	}

	return encoded
}

// Scroll performs a scroll action.
func (w *WheelInput) Scroll(deltaX, deltaY int, origin interface{}, x, y int) {
	action := NewWheelAction(deltaX, deltaY, origin, x, y)
	w.actions = append(w.actions, action)
}
