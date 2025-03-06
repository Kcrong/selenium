//go:build integration_test

package internal

import (
	"net/url"
	"os"
	"testing"

	"github.com/Kcrong/selenium"
	"github.com/stretchr/testify/assert"
)

func TestFindElementsWithCustomHTML_Chrome(t *testing.T) {
	seleniumURL := os.Getenv("SELENIUM_CHROME_URL")
	if seleniumURL == "" {
		t.Skip("SELENIUM_CHROME_URL not set, skipping Chrome test")
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		t.Fatalf("Failed to create WebDriver (chrome): %v", err)
	}
	defer func() {
		assert.NoError(t, wd.Quit())
	}()

	html := `
	<!DOCTYPE html>
	<html>
	  <head><title>FindElements test</title></head>
	  <body>
	    <h1>Testing findElements</h1>
	    <div class="test-class">First Div</div>
	    <div class="test-class">Second Div</div>
	    <div id="unique-id">Third Div</div>
	  </body>
	</html>
	`
	escaped := url.PathEscape(html)
	dataURL := "data:text/html;charset=utf-8," + escaped

	if err := wd.Get(dataURL); err != nil {
		t.Fatalf("Failed to load data URL: %v", err)
	}

	// 예시: 특정 클래스의 요소를 여러 개 찾기
	elems, err := wd.FindElements(selenium.ByCSSSelector, ".test-class")
	assert.NoError(t, err)
	assert.Len(t, elems, 2)

	// 예시: ID로 유일한 요소 찾기
	elem, err := wd.FindElement(selenium.ByID, "unique-id")
	assert.NoError(t, err)
	text, err := elem.Text()
	assert.NoError(t, err)
	assert.Equal(t, "Third Div", text)
}

func TestFindElementsWithCustomHTML_Firefox(t *testing.T) {
	seleniumURL := os.Getenv("SELENIUM_FIREFOX_URL")
	if seleniumURL == "" {
		t.Skip("SELENIUM_FIREFOX_URL not set, skipping Firefox test")
	}

	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		t.Fatalf("Failed to create WebDriver (firefox): %v", err)
	}
	defer func() {
		assert.NoError(t, wd.Quit())
	}()

	html := `
	<!DOCTYPE html>
	<html>
	  <head><title>FindElements test</title></head>
	  <body>
	    <h2>Another HTML for findElements</h2>
	    <p class="test-paragraph">Hello</p>
	    <p class="test-paragraph">World</p>
	    <section id="unique-section">Section content</section>
	  </body>
	</html>
	`
	escaped := url.PathEscape(html)
	dataURL := "data:text/html;charset=utf-8," + escaped

	if err := wd.Get(dataURL); err != nil {
		t.Fatalf("Failed to load data URL: %v", err)
	}

	elems, err := wd.FindElements(selenium.ByCSSSelector, ".test-paragraph")
	assert.NoError(t, err)
	assert.Len(t, elems, 2)

	elem, err := wd.FindElement(selenium.ByID, "unique-section")
	assert.NoError(t, err)
	secText, err := elem.Text()
	assert.NoError(t, err)
	assert.Equal(t, "Section content", secText)
}
