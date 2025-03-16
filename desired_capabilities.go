package selenigo

/*
## Usage

// Create a new DesiredCapabilities instance
caps := NewDesiredCapabilities()

// Get the default capabilities for Firefox
firefoxCaps := caps.Firefox()

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
		BrowserName:         "firefox",
		AcceptInsecureCerts: true,
		// TODO: Add support for "moz:debuggerAddress"
		// "moz:debuggerAddress": true,
	}
}

// InternetExplorer returns the default capabilities for Internet Explorer
func (d *DesiredCapabilities) InternetExplorer() Capabilities {
	return Capabilities{
		BrowserName:  "internet explorer",
		PlatformName: "windows",
	}
}

// Edge returns the default capabilities for Microsoft Edge
func (d *DesiredCapabilities) Edge() Capabilities {
	return Capabilities{
		BrowserName: "MicrosoftEdge",
	}
}

// Chrome returns the default capabilities for Chrome
func (d *DesiredCapabilities) Chrome() Capabilities {
	return Capabilities{
		BrowserName: "chrome",
	}
}

// Safari returns the default capabilities for Safari
func (d *DesiredCapabilities) Safari() Capabilities {
	return Capabilities{
		BrowserName:  "safari",
		PlatformName: "mac",
	}
}

// HTMLUnit returns the default capabilities for HTMLUnit
func (d *DesiredCapabilities) HTMLUnit() Capabilities {
	return Capabilities{
		BrowserName:    "htmlunit",
		BrowserVersion: "",
		PlatformName:   PlatformANY,
	}
}

// HTMLUnitWithJS returns the default capabilities for HTMLUnit with JavaScript enabled
func (d *DesiredCapabilities) HTMLUnitWithJS() Capabilities {
	return Capabilities{
		BrowserName:         "htmlunit",
		BrowserVersion:      "firefox",
		PlatformName:        "ANY",
		IsJavaScriptEnabled: true,
	}
}

// IPhone returns the default capabilities for iPhone
func (d *DesiredCapabilities) IPhone() Capabilities {
	return Capabilities{
		BrowserName:    "iPhone",
		BrowserVersion: "",
		PlatformName:   "mac",
	}
}

// IPad returns the default capabilities for iPad
func (d *DesiredCapabilities) IPad() Capabilities {
	return Capabilities{
		BrowserName:    "iPad",
		BrowserVersion: "",
		PlatformName:   "mac",
	}
}

// WebKitGTK returns the default capabilities for WebKitGTK
func (d *DesiredCapabilities) WebKitGTK() Capabilities {
	return Capabilities{
		BrowserName: "MiniBrowser",
	}
}

// WPEWebKit returns the default capabilities for WPEWebKit
func (d *DesiredCapabilities) WPEWebKit() Capabilities {
	return Capabilities{
		BrowserName: "MiniBrowser",
	}
}
