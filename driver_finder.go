package selenium

import (
	"errors"
	"fmt"
	"os"
)

// DriverFinder is responsible for obtaining the correct driver and associated browser
type DriverFinder struct {
	service *Service
	options *BaseOptions
	paths   map[string]string
	manager *SeleniumManager
}

// NewDriverFinder creates a new DriverFinder instance
func NewDriverFinder(service *Service, options *BaseOptions) *DriverFinder {
	return &DriverFinder{
		service: service,
		options: options,
		paths: map[string]string{
			"driver_path":  "",
			"browser_path": "",
		},
		manager: NewSeleniumManager(),
	}
}

// GetBrowserPath returns the path to the browser binary
func (d *DriverFinder) GetBrowserPath() (string, error) {
	paths, err := d.getBinaryPaths()
	if err != nil {
		return "", err
	}
	return paths["browser_path"], nil
}

// GetDriverPath returns the path to the driver binary
func (d *DriverFinder) GetDriverPath() (string, error) {
	paths, err := d.getBinaryPaths()
	if err != nil {
		return "", err
	}
	return paths["driver_path"], nil
}

var ErrBrowserCapabilityNotFound = errors.New("browserName capability not found")

// getBinaryPaths returns the paths to both the driver and browser binaries
func (d *DriverFinder) getBinaryPaths() (map[string]string, error) {
	if d.paths["driver_path"] != "" {
		return d.paths, nil
	}

	browser := d.options.Capabilities.BrowserName
	if browser == "" {
		return nil, ErrBrowserCapabilityNotFound
	}

	// If service path is specified, use it
	if d.service.Path != "" {
		if _, err := os.Stat(d.service.Path); err != nil {
			return nil, err
		}
		d.paths["driver_path"] = d.service.Path
		return d.paths, nil
	}

	// Use Selenium Manager to find paths
	args := d.toArgs()
	output, err := d.manager.BinaryPaths(args)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, args)
	}

	// Validate driver path
	if driverPath, ok := output["driver_path"]; ok && driverPath != "" {
		if _, err := os.Stat(driverPath); err != nil {
			return nil, fmt.Errorf("%w: %s", err, driverPath)
		}
		d.paths["driver_path"] = driverPath
	}

	// Validate browser path if provided
	if browserPath, ok := output["browser_path"]; ok && browserPath != "" {
		if _, err := os.Stat(browserPath); err != nil {
			return nil, fmt.Errorf("%w: %s", err, browserPath)
		}
		d.paths["browser_path"] = browserPath
	}

	return d.paths, nil
}

// toArgs converts options to command line arguments for Selenium Manager
func (d *DriverFinder) toArgs() []string {
	args := []string{"--browser"}

	browser := d.options.Capabilities.BrowserName
	if browser != "" {
		args = append(args, string(browser))
	}

	if version := d.options.Capabilities.BrowserVersion; version != "" {
		args = append(args, "--browser-version", version)
	}

	// Check if options has binary location
	if caps := d.options.ToCapabilities(); caps != nil {
		if binaryLocation, ok := caps["binary"].(string); ok && binaryLocation != "" {
			args = append(args, "--browser-path", binaryLocation)
		}
	}

	// Handle proxy settings
	proxy := d.options.Proxy
	if proxy != nil {
		if httpProxy := proxy.GetHTTPProxy(); httpProxy != "" {
			args = append(args, "--proxy", httpProxy)
		} else if sslProxy := proxy.GetSSLProxy(); sslProxy != "" {
			args = append(args, "--proxy", sslProxy)
		}
	}

	return args
}

// isExecutable checks if a file exists and is executable
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if the file is regular and has execute permission
	// Note: This is a simplified check and might need to be adjusted
	// based on the operating system
	return !info.IsDir() && (info.Mode()&0o111 != 0)
}
