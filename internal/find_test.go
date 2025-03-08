//go:build integration_test

package internal

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Kcrong/selenium/pkg"
)

func TestIntegration_FindElement_And_FindElements(t *testing.T) {
	t.Parallel()

	testServerURL := os.Getenv("SELENIUM_TEST_SERVER_URL")
	if testServerURL == "" {
		t.Skip("SELENIUM_TEST_SERVER_URL not set, skipping test")
	}

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

			caps := pkg.Capabilities{"browserName": browser.name}
			wd, err := pkg.NewRemote(caps, browserURL)
			require.NoErrorf(t, err, "cannot create WebDriver (%s)", browser.name)
			t.Cleanup(func() {
				assert.NoError(t, wd.Quit())
			})

			path, err := url.JoinPath(testServerURL, "find.html")
			require.NoError(t, err)
			require.NoError(t, wd.Get(path))

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					t.Run("FindElement", func(t *testing.T) {
						t.Parallel()

						elem, err := wd.FindElement(tt.by, tt.value)
						if tt.expectElem {
							require.NoError(t, err)
							assert.NotNil(t, elem, "expected to find an element")
						} else {
							assert.Error(t, err, "expected an error when finding nonexistent element")
						}
					})

					t.Run("FindElements", func(t *testing.T) {
						t.Parallel()

						elements, err := wd.FindElements(tt.by, tt.value)
						require.NoError(t, err)

						if tt.expectElem {
							assert.Greater(t, len(elements), 0, "expected at least one element")
						} else {
							assert.Equal(t, 0, len(elements), "expected no elements")
						}
					})
				})
			}

			time.Sleep(1 * time.Second) // Allow time for browser cleanup
		})
	}
}
