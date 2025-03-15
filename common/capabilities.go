package common

/*
## Usage

// Create a new DesiredCapabilities instance
caps := NewDesiredCapabilities()

// Get the default capabilities for Firefox
firefoxCaps := caps.Firefox()

// Modify the capabilities
modifiedCaps := CopyCapabilities(firefoxCaps)
modifiedCaps["platform"] = "WINDOWS"
modifiedCaps["version"] = "10"
*/

// DesiredCapabilities provides a set of predefined capability sets for different browsers
type DesiredCapabilities struct{}

// NewDesiredCapabilities creates a new DesiredCapabilities instance
func NewDesiredCapabilities() *DesiredCapabilities {
	return &DesiredCapabilities{}
}

// Firefox returns the default capabilities for Firefox
func (d *DesiredCapabilities) Firefox() Capabilities {
	return Capabilities{
		"browserName":         "firefox",
		"acceptInsecureCerts": true,
		"moz:debuggerAddress": true,
	}
}

// InternetExplorer returns the default capabilities for Internet Explorer
func (d *DesiredCapabilities) InternetExplorer() Capabilities {
	return Capabilities{
		"browserName":  "internet explorer",
		"platformName": "windows",
	}
}

// Edge returns the default capabilities for Microsoft Edge
func (d *DesiredCapabilities) Edge() Capabilities {
	return Capabilities{
		"browserName": "MicrosoftEdge",
	}
}

// Chrome returns the default capabilities for Chrome
func (d *DesiredCapabilities) Chrome() Capabilities {
	return Capabilities{
		"browserName": "chrome",
	}
}

// Safari returns the default capabilities for Safari
func (d *DesiredCapabilities) Safari() Capabilities {
	return Capabilities{
		"browserName":  "safari",
		"platformName": "mac",
	}
}

// HTMLUnit returns the default capabilities for HTMLUnit
func (d *DesiredCapabilities) HTMLUnit() Capabilities {
	return Capabilities{
		"browserName": "htmlunit",
		"version":     "",
		"platform":    "ANY",
	}
}

// HTMLUnitWithJS returns the default capabilities for HTMLUnit with JavaScript enabled
func (d *DesiredCapabilities) HTMLUnitWithJS() Capabilities {
	return Capabilities{
		"browserName":       "htmlunit",
		"version":           "firefox",
		"platform":          "ANY",
		"javascriptEnabled": true,
	}
}

// IPhone returns the default capabilities for iPhone
func (d *DesiredCapabilities) IPhone() Capabilities {
	return Capabilities{
		"browserName": "iPhone",
		"version":     "",
		"platform":    "mac",
	}
}

// IPad returns the default capabilities for iPad
func (d *DesiredCapabilities) IPad() Capabilities {
	return Capabilities{
		"browserName": "iPad",
		"version":     "",
		"platform":    "mac",
	}
}

// WebKitGTK returns the default capabilities for WebKitGTK
func (d *DesiredCapabilities) WebKitGTK() Capabilities {
	return Capabilities{
		"browserName": "MiniBrowser",
	}
}

// WPEWebKit returns the default capabilities for WPEWebKit
func (d *DesiredCapabilities) WPEWebKit() Capabilities {
	return Capabilities{
		"browserName": "MiniBrowser",
	}
}

// CopyCapabilities creates a deep copy of the capabilities
func CopyCapabilities(caps Capabilities) Capabilities {
	newCaps := make(Capabilities)
	for k, v := range caps {
		newCaps[k] = v
	}
	return newCaps
}
