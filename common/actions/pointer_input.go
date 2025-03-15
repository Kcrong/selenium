package actions

import (
	"time"

	"github.com/Kcrong/selenium-go/remote/webelement"
)

const (
	DefaultMoveDuration = 250 * time.Millisecond
)

// PointerInput represents a mouse or touch input device.
type PointerInput struct {
	name     string
	kind     string
	actions  []Action
	duration time.Duration
}

// NewPointerInput creates a new pointer input device.
func NewPointerInput(kind, name string, duration time.Duration) *PointerInput {
	if duration == 0 {
		duration = DefaultMoveDuration
	}
	return &PointerInput{
		name:     name,
		kind:     kind,
		actions:  make([]Action, 0),
		duration: duration,
	}
}

// GetName returns the name of the input device.
func (p *PointerInput) GetName() string {
	return p.name
}

// GetType returns the type of the input device.
func (p *PointerInput) GetType() string {
	return "pointer"
}

// ClearActions clears all stored actions.
func (p *PointerInput) ClearActions() {
	p.actions = make([]Action, 0)
}

// CreatePause creates a pause action.
func (p *PointerInput) CreatePause(duration time.Duration) Action {
	return NewPauseAction(duration)
}

// GetActions returns all stored actions.
func (p *PointerInput) GetActions() []Action {
	return p.actions
}

// AddAction adds an action to the device.
func (p *PointerInput) AddAction(action Action) {
	p.actions = append(p.actions, action)
}

// PointerAction represents a pointer action.
type PointerAction struct {
	BaseAction
	Button   int
	Origin   interface{}
	X        int
	Y        int
	Width    int
	Height   int
	Pressure float64
	TangentX float64
	TangentY float64
	TiltX    int
	TiltY    int
	Twist    int
}

// NewPointerAction creates a new pointer action.
func NewPointerAction(actionType string, options map[string]interface{}) *PointerAction {
	action := &PointerAction{
		BaseAction: BaseAction{
			Type:     actionType,
			Duration: 0,
		},
	}

	if button, ok := options["button"].(int); ok {
		action.Button = button
	}
	if origin, ok := options["origin"]; ok {
		action.Origin = origin
	}
	if x, ok := options["x"].(int); ok {
		action.X = x
	}
	if y, ok := options["y"].(int); ok {
		action.Y = y
	}
	if width, ok := options["width"].(int); ok {
		action.Width = width
	}
	if height, ok := options["height"].(int); ok {
		action.Height = height
	}
	if pressure, ok := options["pressure"].(float64); ok {
		action.Pressure = pressure
	}
	if tangentX, ok := options["tangentX"].(float64); ok {
		action.TangentX = tangentX
	}
	if tangentY, ok := options["tangentY"].(float64); ok {
		action.TangentY = tangentY
	}
	if tiltX, ok := options["tiltX"].(int); ok {
		action.TiltX = tiltX
	}
	if tiltY, ok := options["tiltY"].(int); ok {
		action.TiltY = tiltY
	}
	if twist, ok := options["twist"].(int); ok {
		action.Twist = twist
	}

	return action
}

// Encode encodes the pointer action into a format that can be sent to the WebDriver.
func (a *PointerAction) Encode() map[string]interface{} {
	encoded := map[string]interface{}{
		"type": a.Type,
	}

	if a.Duration > 0 {
		encoded["duration"] = a.Duration
	}
	if a.Button != 0 {
		encoded["button"] = a.Button
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
	if a.Width != 0 {
		encoded["width"] = a.Width
	}
	if a.Height != 0 {
		encoded["height"] = a.Height
	}
	if a.Pressure != 0 {
		encoded["pressure"] = a.Pressure
	}
	if a.TangentX != 0 {
		encoded["tangentX"] = a.TangentX
	}
	if a.TangentY != 0 {
		encoded["tangentY"] = a.TangentY
	}
	if a.TiltX != 0 {
		encoded["tiltX"] = a.TiltX
	}
	if a.TiltY != 0 {
		encoded["tiltY"] = a.TiltY
	}
	if a.Twist != 0 {
		encoded["twist"] = a.Twist
	}

	return encoded
}

// Move moves the pointer to a specific location.
func (p *PointerInput) Move(x, y int, origin interface{}) {
	action := NewPointerAction("pointerMove", map[string]interface{}{
		"x":      x,
		"y":      y,
		"origin": origin,
	})
	action.Duration = p.duration
	p.actions = append(p.actions, action)
}

// Click performs a click action.
func (p *PointerInput) Click(button int) {
	p.actions = append(p.actions, NewPointerAction("pointerDown", map[string]interface{}{
		"button": button,
	}))
	p.actions = append(p.actions, NewPointerAction("pointerUp", map[string]interface{}{
		"button": button,
	}))
}

// ClickAndHold performs a click and hold action.
func (p *PointerInput) ClickAndHold(button int) {
	p.actions = append(p.actions, NewPointerAction("pointerDown", map[string]interface{}{
		"button": button,
	}))
}

// Release releases a held button.
func (p *PointerInput) Release(button int) {
	p.actions = append(p.actions, NewPointerAction("pointerUp", map[string]interface{}{
		"button": button,
	}))
}

// DoubleClick performs a double click action.
func (p *PointerInput) DoubleClick(button int) {
	p.Click(button)
	p.Click(button)
}
