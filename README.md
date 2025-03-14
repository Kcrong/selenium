# **selenium**

## **🔹 About**
This is a **forked version** of [tebeka/selenium](https://github.com/tebeka/selenium),  
**maintained for Selenium 4**, specifically for **Chrome**.

For more details on the original project, please refer to the [original repository](https://github.com/tebeka/selenium).

---

## **📦 Installation**
Run the following command to install:
```sh
go get github.com/Kcrong/selenium
```

---

## **🚀 Usage**
Below is an example demonstrating how to use **selenium** in Go:

```go
package main

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/Kcrong/selenium"
	"github.com/Kcrong/selenium/pkg/log"
)

func TestExampleUsage(t *testing.T) {
	// Initialize WebDriver
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, browserURL)
	require.NoError(t, err, "failed to create WebDriver")
	defer func() { assert.NoError(t, wd.Quit()) }()

	// Load test page
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
}
```

---

## **🔧 Integration Testing**
To run integration tests, use the following command:
```sh
make integration
```