//go:build integration_test

package test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Kcrong/selenium"
)

// TestIntegration_NewChromeDriverService verifies that ChromeDriver starts correctly
func TestIntegration_NewChromeDriverService(t *testing.T) {
	t.Parallel()

	// Fetch ChromeDriver path and port
	chromeDriverPath := os.Getenv("CHROMEDRIVER_PATH")
	if chromeDriverPath == "" {
		t.Skip("CHROMEDRIVER_PATH not set, skipping ChromeDriver service test")
	}
	// Random port from 49152 to 65535
	port := 49152 + time.Now().Nanosecond()%16383

	// Start ChromeDriver service
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port)
	require.NoError(t, err, "Failed to start ChromeDriver service")
	t.Cleanup(func() {
		assert.NoError(t, service.Stop(), "Failed to stop ChromeDriver service")
	})

	// Wait for ChromeDriver to be ready
	url := fmt.Sprintf("http://localhost:%d/status", port)
	assert.Eventually(t, func() bool {
		resp, err := http.Get(url)
		if err != nil {
			return false
		}
		t.Cleanup(func() {
			assert.NoError(t, resp.Body.Close(), "Failed to close response body")
		})

		return resp.StatusCode == http.StatusOK

	}, 10*time.Second, 500*time.Millisecond, "ChromeDriver did not start in time")

	// Verify that a WebDriver session can be created
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	require.NoError(t, err, "Failed to create WebDriver session")
	t.Cleanup(func() {
		assert.NoError(t, wd.Quit(), "Failed to quit WebDriver")
	})

	// Verify session is active
	assert.NotEmpty(t, wd.SessionID(), "WebDriver session ID should not be empty")
}
