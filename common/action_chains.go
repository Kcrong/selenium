package common

import (
	"time"

	"github.com/Kcrong/selenium-go"
	"github.com/Kcrong/selenium-go/common/actions"
	"github.com/Kcrong/selenium-go/remote/webelement"
)

// ActionChains provides a way to automate low level interactions such as
// mouse movements, mouse button actions, key press, and context menu interactions.
type ActionChains struct {
	driver  selenium.WebDriver
	actions *actions.ActionBuilder
}

// NewActionChains creates a new ActionChains instance.
func NewActionChains(driver selenium.WebDriver, duration time.Duration) *ActionChains {
	return &ActionChains{
		driver:  driver,
		actions: actions.NewActionBuilder(driver, nil, nil, nil, duration),
	}
}

// Perform performs all stored actions.
func (a *ActionChains) Perform() error {
	return a.actions.Perform()
}

// ResetActions clears actions that are already stored on the remote end.
func (a *ActionChains) ResetActions() error {
	return a.actions.ClearActions()
}

// Click clicks an element.
func (a *ActionChains) Click(element webelement.WebElement) *ActionChains {
	if element != nil {
		a.MoveToElement(element)
	}
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "pointer" {
			if pointer, ok := device.(*actions.PointerInput); ok {
				pointer.Click(0)
			}
		}
	}
	return a
}

// ClickAndHold holds down the left mouse button on an element.
func (a *ActionChains) ClickAndHold(element webelement.WebElement) *ActionChains {
	if element != nil {
		a.MoveToElement(element)
	}
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "pointer" {
			if pointer, ok := device.(*actions.PointerInput); ok {
				pointer.ClickAndHold(0)
			}
		}
	}
	return a
}

// ContextClick performs a context-click (right click) on an element.
func (a *ActionChains) ContextClick(element webelement.WebElement) *ActionChains {
	if element != nil {
		a.MoveToElement(element)
	}
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "pointer" {
			if pointer, ok := device.(*actions.PointerInput); ok {
				pointer.Click(2) // Right button
			}
		}
	}
	return a
}

// DoubleClick double-clicks an element.
func (a *ActionChains) DoubleClick(element webelement.WebElement) *ActionChains {
	if element != nil {
		a.MoveToElement(element)
	}
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "pointer" {
			if pointer, ok := device.(*actions.PointerInput); ok {
				pointer.DoubleClick(0)
			}
		}
	}
	return a
}

// DragAndDrop drags and drops an element.
func (a *ActionChains) DragAndDrop(source, target webelement.WebElement) *ActionChains {
	return a.ClickAndHold(source).MoveToElement(target).Release(nil)
}

// DragAndDropByOffset drags and drops an element by an offset.
func (a *ActionChains) DragAndDropByOffset(source webelement.WebElement, xOffset, yOffset int) *ActionChains {
	return a.ClickAndHold(source).MoveByOffset(xOffset, yOffset).Release(nil)
}

// KeyDown sends a key press only, without releasing it.
func (a *ActionChains) KeyDown(key string, element webelement.WebElement) *ActionChains {
	if element != nil {
		a.Click(element)
	}
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "key" {
			if keyboard, ok := device.(*actions.KeyInput); ok {
				keyboard.KeyDown(key)
			}
		}
	}
	return a
}

// KeyUp releases a modifier key.
func (a *ActionChains) KeyUp(key string, element webelement.WebElement) *ActionChains {
	if element != nil {
		a.Click(element)
	}
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "key" {
			if keyboard, ok := device.(*actions.KeyInput); ok {
				keyboard.KeyUp(key)
			}
		}
	}
	return a
}

// MoveByOffset moves the mouse by an offset.
func (a *ActionChains) MoveByOffset(xOffset, yOffset int) *ActionChains {
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "pointer" {
			if pointer, ok := device.(*actions.PointerInput); ok {
				pointer.Move(xOffset, yOffset, nil)
			}
		}
	}
	return a
}

// MoveToElement moves the mouse to the middle of an element.
func (a *ActionChains) MoveToElement(element webelement.WebElement) *ActionChains {
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "pointer" {
			if pointer, ok := device.(*actions.PointerInput); ok {
				pointer.Move(0, 0, element)
			}
		}
	}
	return a
}

// MoveToElementWithOffset moves the mouse to an element with offset.
func (a *ActionChains) MoveToElementWithOffset(element webelement.WebElement, xOffset, yOffset int) *ActionChains {
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "pointer" {
			if pointer, ok := device.(*actions.PointerInput); ok {
				pointer.Move(xOffset, yOffset, element)
			}
		}
	}
	return a
}

// Pause pauses all inputs for the specified duration.
func (a *ActionChains) Pause(duration time.Duration) *ActionChains {
	for _, device := range a.actions.GetDevices() {
		device.AddAction(device.CreatePause(duration))
	}
	return a
}

// Release releases a held mouse button.
func (a *ActionChains) Release(element webelement.WebElement) *ActionChains {
	if element != nil {
		a.MoveToElement(element)
	}
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "pointer" {
			if pointer, ok := device.(*actions.PointerInput); ok {
				pointer.Release(0)
			}
		}
	}
	return a
}

// SendKeys sends keys to the active element.
func (a *ActionChains) SendKeys(keys string) *ActionChains {
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "key" {
			if keyboard, ok := device.(*actions.KeyInput); ok {
				keyboard.SendKeys(keys)
			}
		}
	}
	return a
}

// SendKeysToElement sends keys to a specific element.
func (a *ActionChains) SendKeysToElement(element webelement.WebElement, keys string) *ActionChains {
	return a.Click(element).SendKeys(keys)
}

// ScrollToElement scrolls to an element.
func (a *ActionChains) ScrollToElement(element webelement.WebElement) *ActionChains {
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "wheel" {
			if wheel, ok := device.(*actions.WheelInput); ok {
				wheel.Scroll(0, 0, element, 0, 0)
			}
		}
	}
	return a
}

// ScrollByAmount scrolls by the specified amount.
func (a *ActionChains) ScrollByAmount(deltaX, deltaY int) *ActionChains {
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "wheel" {
			if wheel, ok := device.(*actions.WheelInput); ok {
				wheel.Scroll(deltaX, deltaY, nil, 0, 0)
			}
		}
	}
	return a
}

// ScrollFromOrigin scrolls from a specific origin.
func (a *ActionChains) ScrollFromOrigin(origin *actions.ScrollOrigin, deltaX, deltaY int) *ActionChains {
	for _, device := range a.actions.GetDevices() {
		if device.GetType() == "wheel" {
			if wheel, ok := device.(*actions.WheelInput); ok {
				wheel.Scroll(deltaX, deltaY, origin.Element, origin.XOffset, origin.YOffset)
			}
		}
	}
	return a
}
