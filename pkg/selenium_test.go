package pkg_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/Kcrong/selenium/pkg"
	"github.com/Kcrong/selenium/pkg/chrome"
	"github.com/Kcrong/selenium/pkg/firefox"
	"github.com/Kcrong/selenium/pkg/log"
)

func TestCapabilities(t *testing.T) {
	caps := pkg.Capabilities{}
	// 1) Chrome options
	chromeOpts := chrome.Capabilities{
		Path: "/usr/bin/chrome",
		Args: []string{"--headless"},
	}
	caps.AddChrome(chromeOpts)

	if _, ok := caps[chrome.CapabilitiesKey]; !ok {
		t.Fatalf("expected caps to contain key %q for Chrome, but not found", chrome.CapabilitiesKey)
	}

	// 2) Firefox options
	ffOpts := firefox.Capabilities{
		Binary: "/usr/bin/firefox",
		Args:   []string{"--devtools"},
	}
	caps.AddFirefox(ffOpts)
	if _, ok := caps[firefox.CapabilitiesKey]; !ok {
		t.Fatalf("expected caps to contain key %q for Firefox, but not found", firefox.CapabilitiesKey)
	}

	// 3) Proxy
	px := pkg.Proxy{
		Type:          pkg.Manual,
		HTTP:          "http://myproxy.example.com",
		HTTPPort:      8080,
		AutoconfigURL: "http://someconfig.example.com",
	}
	caps.AddProxy(px)
	if _, ok := caps["proxy"]; !ok {
		t.Fatalf("expected caps to contain 'proxy', but not found")
	}

	// 4) Logging
	logCaps := log.Capabilities{
		log.Browser: log.Info,
	}
	caps.AddLogging(logCaps)
	if _, ok := caps[log.CapabilitiesKey]; !ok {
		t.Fatalf("expected caps to contain key %q for logging, but not found", log.CapabilitiesKey)
	}

	// 5) SetLogLevel
	caps.SetLogLevel(log.Driver, log.Debug)
	logs := caps[log.CapabilitiesKey].(log.Capabilities)
	if logs[log.Driver] != log.Debug {
		t.Errorf("expected log level for driver to be 'debug', got %v", logs[log.Driver])
	}
}

func TestCookie(t *testing.T) {
	// 단순 Cookie 구조체 직렬화/역직렬화 테스트
	c := pkg.Cookie{
		Name:     "testCookie",
		Value:    "someValue",
		Path:     "/",
		Domain:   "example.com",
		Secure:   true,
		Expiry:   1234567890,
		HTTPOnly: true,
		SameSite: pkg.SameSiteLax,
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("failed to marshal cookie: %v", err)
	}

	var c2 pkg.Cookie
	if err := json.Unmarshal(data, &c2); err != nil {
		t.Fatalf("failed to unmarshal cookie: %v", err)
	}

	if !reflect.DeepEqual(c, c2) {
		t.Errorf("cookie after marshal/unmarshal mismatch\n got:  %+v\n want: %+v", c2, c)
	}
}

func TestKeyConstants(t *testing.T) {
	// 예시: NullKey 가 \ue000 인지 확인
	if pkg.NullKey != "\ue000" {
		t.Errorf("NullKey expected '\\ue000', got %q", pkg.NullKey)
	}
	// 필요한 만큼 상수를 확인할 수도 있음
}

func TestCapabilitiesJSONMarshal(t *testing.T) {
	// Capabilities 자체를 JSON 직렬화해보기
	caps := pkg.Capabilities{
		"browserName": "chrome",
		"customKey":   "customValue",
	}
	data, err := json.Marshal(caps)
	if err != nil {
		t.Fatalf("failed to marshal Capabilities: %v", err)
	}
	// JSON 결과가 {"browserName":"chrome","customKey":"customValue"} 형태인지 부분확인
	s := string(data)
	if !contains(s, `"browserName":"chrome"`) || !contains(s, `"customKey":"customValue"`) {
		t.Errorf("unexpected JSON: %s", s)
	}
}

// contains is a simple substring check helper.
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (strings.Contains(s, sub))
}
