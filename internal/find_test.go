//go:build integration_test

package internal

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Kcrong/selenium/pkg"
)

func TestIntegration_FindElement(t *testing.T) {
	tests := []struct {
		name       string
		by         string
		value      string
		expectElem bool
	}{
		{"FindByID", pkg.ByID, "test-id", true},
		{"FindByXPath", pkg.ByXPATH, "//*[@id='test-id']", true},
		{"FindByLinkText", pkg.ByLinkText, "Click me", true},
		{"FindByPartialLinkText", pkg.ByPartialLinkText, "Click", true},
		{"FindByName", pkg.ByName, "test-name", true},
		{"FindByTagName", pkg.ByTagName, "p", true},
		{"FindByClassName", pkg.ByClassName, "test-class", true},
		{"FindByCSSSelector", pkg.ByCSSSelector, ".test-class", true},
		{"FindByInvalid", pkg.ByCSSSelector, ".invalid-class", false},
	}

	browsers := []struct {
		name string
		env  string
	}{
		{"Chrome", "SELENIUM_CHROME_URL"},
		{"Firefox", "SELENIUM_FIREFOX_URL"},
	}

	for _, browser := range browsers {
		t.Run(browser.name, func(t *testing.T) {
			// i.e "http://localhost:55175"
			browserURL := os.Getenv(browser.env)
			if browserURL == "" {
				t.Skipf("%s not set, skipping test", browser.env)
			}

			caps := pkg.Capabilities{"browserName": browser.name}
			wd, err := pkg.NewRemote(caps, browserURL)
			require.NoError(t, err, "cannot create WebDriver (%s): %v", browser.name, err)
			defer func(wd pkg.WebDriver) {
				err := wd.Quit()
				assert.NoError(t, err)
			}(wd)

			// Load test page with predefined HTML elements
			htmlContent := `
			<!DOCTYPE html>
			<html>
			<head><title>Test Page</title></head>
			<body>
				<p id="test-id">Hello</p>
				<p name="test-name">Test Name</p>
				<p class="test-class">Test Class</p>
				<a href="#" id="test-link">Click me</a>
			</body>
			</html>`
			tempFile, err := os.CreateTemp("", "testpage-*.html")
			assert.NoError(t, err)
			defer func(name string) {
				err := os.Remove(name)
				assert.NoError(t, err)
			}(tempFile.Name())

			_, err = tempFile.WriteString(htmlContent)
			assert.NoError(t, err)
			assert.NoError(t, tempFile.Close())

			fileURL := "file://" + tempFile.Name()
			err = wd.Get(fileURL)
			assert.NoError(t, err)

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					elem, err := wd.FindElement(tt.by, tt.value)
					if tt.expectElem {
						require.NoError(t, err)
						assert.NotNil(t, elem, "expected to find an element")
					} else {
						assert.Error(t, err, "expected an error when finding nonexistent element")
					}
				})
			}

			time.Sleep(1 * time.Second) // Allow time for browser cleanup
		})
	}
}

func TestIntegration_FindElements(t *testing.T) {
	tests := []struct {
		name       string
		by         string
		value      string
		expectElem bool
	}{
		{"FindByID", pkg.ByID, "test-id", true},
		{"FindByXPath", pkg.ByXPATH, "//*[@id='test-id']", true},
		{"FindByLinkText", pkg.ByLinkText, "Click me", true},
		{"FindByPartialLinkText", pkg.ByPartialLinkText, "Click", true},
		{"FindByName", pkg.ByName, "test-name", true},
		{"FindByTagName", pkg.ByTagName, "p", true},
		{"FindByClassName", pkg.ByClassName, "test-class", true},
		{"FindByCSSSelector", pkg.ByCSSSelector, ".test-class", true},
		{"FindByInvalid", pkg.ByCSSSelector, ".invalid-class", false},
	}

	browsers := []struct {
		name string
		env  string
	}{
		{"Chrome", "SELENIUM_CHROME_URL"},
		{"Firefox", "SELENIUM_FIREFOX_URL"},
	}

	for _, browser := range browsers {
		t.Run(browser.name, func(t *testing.T) {
			// i.e "http://localhost:55175"
			browserURL := os.Getenv(browser.env)
			if browserURL == "" {
				t.Skipf("%s not set, skipping test", browser.env)
			}

			caps := pkg.Capabilities{"browserName": browser.name}
			wd, err := pkg.NewRemote(caps, browserURL)
			require.NoErrorf(t, err, "cannot create WebDriver (%s)", browser.name)
			defer func(wd pkg.WebDriver) {
				err := wd.Quit()
				assert.NoError(t, err)
			}(wd)

			// Load test page with predefined HTML elements
			htmlContent := `
			<!DOCTYPE html>
			<html>
			<head><title>Test Page</title></head>
			<body>
				<p id="test-id">Hello</p>
				<p name="test-name">Test Name</p>
				<p class="test-class">Test Class</p>
				<a href="#" id="test-link">Click me</a>
			</body>
			</html>`
			tempFile, err := os.CreateTemp("", "testpage-*.html")
			assert.NoError(t, err)
			defer func(name string) {
				err := os.Remove(name)
				assert.NoError(t, err)
			}(tempFile.Name())

			_, err = tempFile.WriteString(htmlContent)
			assert.NoError(t, err)
			assert.NoError(t, tempFile.Close())

			fileURL := "file://" + tempFile.Name()
			err = wd.Get(fileURL)
			assert.NoError(t, err)

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					elements, err := wd.FindElements(tt.by, tt.value)
					require.NoError(t, err)
					if tt.expectElem {
						assert.Greater(t, len(elements), 0, "expected at least one element")
					} else {
						assert.Equal(t, 0, len(elements), "expected no elements")
					}
				})
			}

			time.Sleep(1 * time.Second) // Allow time for browser cleanup
		})
	}
}
