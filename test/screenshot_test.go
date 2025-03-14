//go:build integration_test

package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Kcrong/selenium"
)

// TestIntegration_Screenshot captures a screenshot and verifies the image data
func runScreenshotTest(t *testing.T, browser, url string) {
	t.Helper()

	// Create WebDriver session
	caps := selenium.Capabilities{"browserName": browser}
	wd, err := selenium.NewRemote(caps, url)
	require.NoError(t, err, "Failed to create WebDriver")
	t.Cleanup(func() {
		assert.NoError(t, wd.Quit(), "Failed to quit WebDriver")
	})

	// Navigate to test page (e.g., example.com)
	testPage := "https://example.com"
	err = wd.Get(testPage)
	require.NoError(t, err, "Failed to load test page")

	// Take a full-page screenshot
	imgData, err := wd.Screenshot()
	require.NoError(t, err, "Failed to capture screenshot")

	// Validate screenshot data
	assert.NotEmpty(t, imgData, "Screenshot should not be empty")
	assert.Greater(t, len(imgData), 1000, "Screenshot data should be reasonably sized")
}

// TestIntegration_Screenshot_Chrome runs the test on Chrome
func TestIntegration_Screenshot_Chrome(t *testing.T) {
	t.Parallel()

	url := os.Getenv("SELENIUM_CHROME_URL")
	if url == "" {
		t.Skip("SELENIUM_CHROME_URL not set, skipping Chrome screenshot test")
	}
	runScreenshotTest(t, "chrome", url)
}

// TestIntegration_Screenshot_Firefox runs the test on Firefox
func TestIntegration_Screenshot_Firefox(t *testing.T) {
	t.Parallel()

	url := os.Getenv("SELENIUM_FIREFOX_URL")
	if url == "" {
		t.Skip("SELENIUM_FIREFOX_URL not set, skipping Firefox screenshot test")
	}
	runScreenshotTest(t, "firefox", url)
}
