//go:build integration_test

package internal

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Kcrong/selenium/pkg"
)

const testHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Mouse Action Test</title>
    <script>
        function onClick(event) {
            let btn = event.button;
            let target = document.getElementById("testButton");
            target.className = "clicked-" + btn;
        }
    </script>
</head>
<body>
    <button id="testButton" onclick="onClick(event)">Click Me</button>
</body>
</html>
`

func serveTestHTML(t *testing.T) string {
	t.Helper()

	file, err := os.CreateTemp("", "mouse_action_test_*.html")
	if err != nil {
		panic(fmt.Sprintf("failed to create temp HTML file: %v", err))
	}
	defer func(file *os.File) {
		assert.NoError(t, file.Close())
	}(file)

	_, err = file.WriteString(testHTML)
	if err != nil {
		panic(fmt.Sprintf("failed to write test HTML content: %v", err))
	}

	return "file://" + file.Name()
}

func runMouseActionTest(t *testing.T, browser, url string) {
	t.Helper()

	caps := pkg.Capabilities{"browserName": browser}
	wd, err := pkg.NewRemote(caps, url)
	if err != nil {
		t.Fatalf("Failed to create WebDriver (%s): %v", browser, err)
	}
	defer func(wd pkg.WebDriver) {
		err := wd.Quit()
		assert.NoError(t, err)
	}(wd)

	// 테스트 HTML 페이지 로드
	htmlFile := serveTestHTML(t)
	err = wd.Get(htmlFile)
	require.NoError(t, err, "Failed to open HTML file")

	// 버튼 찾기
	btn, err := wd.FindElement(pkg.ByID, "testButton")
	require.NoError(t, err, "Failed to find testButton element")

	testCases := []struct {
		name          string
		mouseButton   pkg.MouseButton
		expectedClass string
	}{
		{"Left Click", pkg.LeftButton, "clicked-0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = btn.Click()
			require.NoError(t, err, "Failed to click button")

			time.Sleep(500 * time.Millisecond) // JavaScript 실행 대기

			classAttr, err := btn.GetAttribute("class")
			require.NoErrorf(t, err, "Failed to get class attribute using %s: %v", tc.name, err)
			assert.Equal(t, tc.expectedClass, classAttr, "Expected button class to match mouse button action")
		})
	}
}

func TestIntegration_MouseActions_Chrome(t *testing.T) {
	url := os.Getenv("SELENIUM_CHROME_URL")
	if url == "" {
		t.Skip("SELENIUM_CHROME_URL not set, skipping Chrome test")
	}
	runMouseActionTest(t, "chrome", url)
}

func TestIntegration_MouseActions_Firefox(t *testing.T) {
	url := os.Getenv("SELENIUM_FIREFOX_URL")
	if url == "" {
		t.Skip("SELENIUM_FIREFOX_URL not set, skipping Firefox test")
	}
	runMouseActionTest(t, "firefox", url)
}
