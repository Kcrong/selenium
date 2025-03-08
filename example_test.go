//go:build example

package selenium_test

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Kcrong/selenium"
	"github.com/Kcrong/selenium/log"
)

func TestExampleUsage_Chrome(t *testing.T) {
	/*
		Before running this test, you need to start the web mock server.
		Please refer docker-compose.yml for more details.
	*/

	// Set up WebDriver for Chrome
	browserURL := os.Getenv("SELENIUM_CHROME_URL")
	if browserURL == "" {
		t.Skip("SELENIUM_CHROME_URL not set, skipping test")
	}

	webMockServerURL := os.Getenv("SELENIUM_TEST_SERVER_URL")
	if webMockServerURL == "" {
		t.Skip("SELENIUM_TEST_SERVER_URL not set, skipping test")
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, browserURL)
	require.NoError(t, err, "failed to create WebDriver")
	defer func() { assert.NoError(t, wd.Quit()) }()

	path, err := url.JoinPath(webMockServerURL, "find.html")
	require.NoError(t, err)
	require.NoError(t, wd.Get(path))

	// Find elements
	elem, err := wd.FindElement(selenium.ByID, "test-id")
	require.NoError(t, err)
	assert.NotNil(t, elem, "expected to find an element by ID")

	elements, err := wd.FindElements(selenium.ByClassName, "test-class")
	require.NoError(t, err)
	assert.Greater(t, len(elements), 0, "expected to find elements by class name")

	// Click button
	btn, err := wd.FindElement(selenium.ByID, "testButton")
	require.NoError(t, err)
	require.NoError(t, btn.Click(), "failed to click button")

	// Perform Key Actions
	require.NoError(t, wd.KeyDown("A"))
	require.NoError(t, wd.KeyUp("A"))

	// Manage Cookies
	cookie := &selenium.Cookie{Name: "testCookie", Value: "12345"}
	require.NoError(t, wd.AddCookie(cookie))
	cookies, err := wd.GetCookies()
	require.NoError(t, err)
	assert.Greater(t, len(cookies), 0, "expected at least one cookie")

	// Take Screenshot
	screenshot, err := wd.Screenshot()
	require.NoError(t, err)
	assert.NotEmpty(t, screenshot, "expected non-empty screenshot")

	// Resize Window
	require.NoError(t, wd.ResizeWindow("", 1024, 768))

	// Dismiss Alert (if present)
	require.NoError(t, wd.DismissAlert())

	// Retrieve Logs
	logs, err := wd.Log(log.Browser)
	require.NoError(t, err)
	fmt.Println("Browser Logs:", logs)

	// Perform Mouse Action
	require.NoError(t, wd.PerformActions())

	time.Sleep(1 * time.Second) // Allow time for cleanup
}
