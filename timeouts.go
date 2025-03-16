package selenigo

import (
	"fmt"
	"time"
)

// Timeouts represents the various timeout settings for WebDriver operations.
type Timeouts struct {
	ImplicitWait time.Duration `json:"implicit"`
	PageLoad     time.Duration `json:"pageLoad"`
	Script       time.Duration `json:"script"`
}

var _ Convertible = (*Timeouts)(nil)

// NewTimeouts creates a new Timeouts instance with the specified durations.
func NewTimeouts(implicitWait, pageLoad, script float64) *Timeouts {
	return &Timeouts{
		ImplicitWait: secondsToDuration(implicitWait),
		PageLoad:     secondsToDuration(pageLoad),
		Script:       secondsToDuration(script),
	}
}

// SetImplicitWait sets how many seconds to wait when searching for elements.
func (t *Timeouts) SetImplicitWait(seconds float64) {
	t.ImplicitWait = secondsToDuration(seconds)
}

// GetImplicitWait returns how many seconds to wait when searching for elements.
func (t *Timeouts) GetImplicitWait() float64 {
	return durationToSeconds(t.ImplicitWait)
}

// SetPageLoad sets how many seconds to wait for the page to load.
func (t *Timeouts) SetPageLoad(seconds float64) {
	t.PageLoad = secondsToDuration(seconds)
}

// GetPageLoad returns how many seconds to wait for the page to load.
func (t *Timeouts) GetPageLoad() float64 {
	return durationToSeconds(t.PageLoad)
}

// SetScript sets how many seconds to wait for an asynchronous script to finish execution.
func (t *Timeouts) SetScript(seconds float64) {
	t.Script = secondsToDuration(seconds)
}

// GetScript returns how many seconds to wait for an asynchronous script to finish execution.
func (t *Timeouts) GetScript() float64 {
	return durationToSeconds(t.Script)
}

// ToCapabilities converts the Timeouts to a map suitable for use in a capabilities object.
func (t *Timeouts) ToCapabilities() map[string]interface{} {
	const fieldCount = 3
	timeouts := make(map[string]interface{}, fieldCount)

	if t.ImplicitWait > 0 {
		timeouts["implicit"] = int(t.ImplicitWait.Milliseconds())
	}

	if t.PageLoad > 0 {
		timeouts["pageLoad"] = int(t.PageLoad.Milliseconds())
	}

	if t.Script > 0 {
		timeouts["script"] = int(t.Script.Milliseconds())
	}

	return timeouts
}

// secondsToDuration converts seconds to time.Duration.
func secondsToDuration(seconds float64) time.Duration {
	if seconds < 0 {
		panic(fmt.Sprintf("timeout value must be non-negative: %f", seconds))
	}

	return time.Duration(seconds * float64(time.Second))
}

// durationToSeconds converts time.Duration to seconds.
func durationToSeconds(d time.Duration) float64 {
	return d.Seconds()
}
