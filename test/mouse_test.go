//go:build integration_test

package test

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Kcrong/selenium"
)

func TestIntegration_MouseActions(t *testing.T) {
	t.Parallel()

	testServerURL := os.Getenv("SELENIUM_TEST_SERVER_URL")
	if testServerURL == "" {
		t.Skip("SELENIUM_TEST_SERVER_URL not set, skipping test")
	}

	browsers := []struct {
		name string
		env  string
	}{
		{"chrome", "SELENIUM_CHROME_URL"},
		{"firefox", "SELENIUM_FIREFOX_URL"},
	}

	for _, browser := range browsers {
		t.Run(browser.name, func(t *testing.T) {
			t.Parallel()

			browserURL := os.Getenv(browser.env)
			if browserURL == "" {
				t.Skipf("%s not set, skipping test", browser.env)
			}

			caps := selenium.Capabilities{"browserName": browser.name}
			wd, err := selenium.NewRemote(caps, browserURL)
			require.NoError(t, err, "cannot create WebDriver (%s): %v", browser.name, err)
			t.Cleanup(func() {
				assert.NoError(t, wd.Quit(), "Failed to quit WebDriver")
			})

			// Load mouse test page
			mouseTestPage, err := url.JoinPath(testServerURL, "mouse.html")
			require.NoError(t, err)
			require.NoError(t, wd.Get(mouseTestPage))

			// Find button element
			btn, err := wd.FindElement(selenium.ByID, "testButton")
			require.NoError(t, err, "Failed to find testButton element")

			testCases := []struct {
				name          string
				mouseButton   selenium.MouseButton
				expectedClass string
			}{
				{"Left Click", selenium.LeftButton, "clicked-0"},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()

					require.NoError(t, btn.Click(), "Failed to click button")

					time.Sleep(500 * time.Millisecond) // Wait for JavaScript event processing

					classAttr, err := btn.GetAttribute("class")
					require.NoErrorf(t, err, "Failed to get class attribute using %s: %v", tc.name, err)
					assert.Equal(t, tc.expectedClass, classAttr, "Expected button class to match mouse button action")
				})
			}
		})
	}
}
