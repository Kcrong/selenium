package selenigo

import (
	"errors"
	"fmt"
)

// Orientation represents the page orientation type.
type Orientation string

const (
	// Portrait orientation.
	Portrait Orientation = "portrait"
	// Landscape orientation.
	Landscape Orientation = "landscape"
)

// MarginOptions represents the margin settings for printing.
type MarginOptions struct {
	Left   *float64 `json:"left,omitempty"`
	Right  *float64 `json:"right,omitempty"`
	Top    *float64 `json:"top,omitempty"`
	Bottom *float64 `json:"bottom,omitempty"`
}

// PageOptions represents the page size settings for printing.
type PageOptions struct {
	Width  *float64 `json:"width,omitempty"`
	Height *float64 `json:"height,omitempty"`
}

// PrintOptions represents all available options for printing a page.
type PrintOptions struct {
	margin      MarginOptions
	page        PageOptions
	background  *bool
	orientation *Orientation
	scale       *float64
	shrinkToFit *bool
	pageRanges  []string
}

// Predefined page sizes (in centimeters).
//
//nolint:mnd // Magic numbers are used for page sizes.
var (
	A4      = PageOptions{Width: ptr(21.0), Height: ptr(29.7)}
	Legal   = PageOptions{Width: ptr(21.59), Height: ptr(35.56)}
	Letter  = PageOptions{Width: ptr(21.59), Height: ptr(27.94)}
	Tabloid = PageOptions{Width: ptr(43.18), Height: ptr(27.94)}
)

// NewPrintOptions creates a new PrintOptions instance.
func NewPrintOptions() *PrintOptions {
	//nolint:exhaustruct // No need to initialize all fields.
	return &PrintOptions{}
}

// SetPageSize sets the page size using predefined or custom dimensions.
func (p *PrintOptions) SetPageSize(pageSize PageOptions) {
	p.page = pageSize
}

// SetMarginTop sets the top margin.
func (p *PrintOptions) SetMarginTop(value float64) error {
	if err := validateNumProperty("margin top", value); err != nil {
		return err
	}

	p.margin.Top = &value

	return nil
}

// SetMarginBottom sets the bottom margin.
func (p *PrintOptions) SetMarginBottom(value float64) error {
	if err := validateNumProperty("margin bottom", value); err != nil {
		return err
	}

	p.margin.Bottom = &value

	return nil
}

// SetMarginLeft sets the left margin.
func (p *PrintOptions) SetMarginLeft(value float64) error {
	if err := validateNumProperty("margin left", value); err != nil {
		return err
	}

	p.margin.Left = &value

	return nil
}

// SetMarginRight sets the right margin.
func (p *PrintOptions) SetMarginRight(value float64) error {
	if err := validateNumProperty("margin right", value); err != nil {
		return err
	}

	p.margin.Right = &value

	return nil
}

var ErrInvalidScale = errors.New("scale must be between 0.1 and 2.0")

// SetScale sets the scale of the page.
func (p *PrintOptions) SetScale(value float64) error {
	if err := validateNumProperty("scale", value); err != nil {
		return err
	}

	if value < 0.1 || value > 2.0 {
		return fmt.Errorf("%w: %f", ErrInvalidScale, value)
	}

	p.scale = &value

	return nil
}

var ErrInvalidOrientation = errors.New("orientation must be either 'portrait' or 'landscape'")

// SetOrientation sets the page orientation.
func (p *PrintOptions) SetOrientation(value Orientation) error {
	if value != Portrait && value != Landscape {
		return fmt.Errorf("%w: %s", ErrInvalidOrientation, value)
	}

	p.orientation = &value

	return nil
}

// SetBackground sets whether to print background graphics.
func (p *PrintOptions) SetBackground(value bool) {
	p.background = &value
}

// SetShrinkToFit sets whether to shrink the page content to fit.
func (p *PrintOptions) SetShrinkToFit(value bool) {
	p.shrinkToFit = &value
}

// SetPageRanges sets the page ranges to print.
func (p *PrintOptions) SetPageRanges(ranges []string) {
	p.pageRanges = ranges
}

// ToMap converts the PrintOptions to a map for JSON serialization.
//
//nolint:cyclop,gocyclo,gocognit // Mapping is straightforward and doesn't need to be split.
func (p *PrintOptions) ToMap() map[string]interface{} {
	result := make(map[string]interface{})

	// Add margin options if any are set
	margin := make(map[string]interface{})
	if p.margin.Top != nil {
		margin["top"] = *p.margin.Top
	}

	if p.margin.Bottom != nil {
		margin["bottom"] = *p.margin.Bottom
	}

	if p.margin.Left != nil {
		margin["left"] = *p.margin.Left
	}

	if p.margin.Right != nil {
		margin["right"] = *p.margin.Right
	}

	if len(margin) > 0 {
		result["margin"] = margin
	}

	// Add page options if any are set
	page := make(map[string]interface{})
	if p.page.Width != nil {
		page["width"] = *p.page.Width
	}

	if p.page.Height != nil {
		page["height"] = *p.page.Height
	}

	if len(page) > 0 {
		result["page"] = page
	}

	if p.background != nil {
		result["background"] = *p.background
	}

	if p.orientation != nil {
		result["orientation"] = *p.orientation
	}

	if p.scale != nil {
		result["scale"] = *p.scale
	}

	if p.shrinkToFit != nil {
		result["shrinkToFit"] = *p.shrinkToFit
	}

	if len(p.pageRanges) > 0 {
		result["pageRanges"] = p.pageRanges
	}

	return result
}

var ErrNegativeValue = errors.New("value must be non-negative")

// validateNumProperty validates a numeric property value.
func validateNumProperty(propertyName string, value float64) error {
	if value < 0 {
		return fmt.Errorf("%w: property %s must be non-negative", ErrNegativeValue, propertyName)
	}

	return nil
}

// ptr returns a pointer to the given value.
func ptr[T any](v T) *T {
	return &v
}
