//go:build integration_test

package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Kcrong/selenium"
)

func TestNewSession(t *testing.T) {
	t.Parallel()

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

			// i.e "http://localhost:55175"
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

			require.NoError(t, wd.Get("https://www.example.com"))
		})
	}
}
