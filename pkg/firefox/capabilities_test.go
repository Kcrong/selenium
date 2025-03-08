package firefox

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestEmptyCapabilities 는 아무 옵션도 채우지 않은 Capabilities 가 직렬화되었을 때 "{}" 인지 확인합니다.
func TestEmptyCapabilities(t *testing.T) {
	caps := Capabilities{}
	data, err := json.Marshal(caps)
	if err != nil {
		t.Fatalf("json.Marshal(Capabilities{}) returned error: %v", err)
	}
	got, want := string(data), "{}"
	if got != want {
		t.Fatalf("json.Marshal(Capabilities{}) = %q, want %q", got, want)
	}
}

// TestPartialCapabilities 는 Binary, Args 등의 일부만 채운 경우 직렬화 결과를 확인합니다.
func TestPartialCapabilities(t *testing.T) {
	caps := Capabilities{
		Binary: "/usr/bin/firefox",
		Args:   []string{"--headless", "--devtools"},
	}
	data, err := json.Marshal(caps)
	if err != nil {
		t.Fatalf("json.Marshal(...) returned error: %v", err)
	}
	got := string(data)

	// 기대 JSON 구조: {"binary":"/usr/bin/firefox","args":["--headless","--devtools"]}
	// 순서가 다를 수 있으므로, 특정 필드가 들어있는지만 간단히 확인합니다.
	if !contains(got, `"binary":"/usr/bin/firefox"`) {
		t.Errorf("serialized JSON missing binary; got: %s", got)
	}
	if !contains(got, `"args":["--headless","--devtools"]`) {
		t.Errorf("serialized JSON missing args; got: %s", got)
	}
}

// TestSetProfile 는 임시 디렉토리를 생성하고 간단한 user.js 파일을 넣은 뒤,
// SetProfile 로 압축+base64 인코딩된 프로필이 Capabilities.Profile 에 설정되는지 확인합니다.
func TestSetProfile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "firefox-profile-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)

	// 간단한 user.js 파일 생성
	userJSPath := filepath.Join(tmpDir, "user.js")
	content := []byte(`user_pref("some.test.pref", true);`)
	if err := os.WriteFile(userJSPath, content, 0644); err != nil {
		t.Fatalf("failed to write user.js: %v", err)
	}

	caps := Capabilities{}
	if err := caps.SetProfile(tmpDir); err != nil {
		t.Fatalf("SetProfile returned error: %v", err)
	}

	if caps.Profile == "" {
		t.Fatalf("expected Profile to be a non-empty base64 string, got empty string")
	}

	// base64 디코딩해서 압축 결과가 user.js 파일을 포함하는지 대략 확인해봄
	decoded, err := base64.StdEncoding.DecodeString(caps.Profile)
	if err != nil {
		t.Fatalf("failed to decode base64 from Profile: %v", err)
	}

	// 간단히 "user.js"라는 문자열이 압축 파일 내부에 있는지 검사
	// (정확히 zip 해제하여 내용 확인까지 할 수도 있으나 여기서는 간단한 수준으로)
	if !bytes.Contains(decoded, []byte("user.js")) {
		t.Errorf("expected profile zip to contain 'user.js', but not found")
	}
}

// TestFullCapabilities 는 모든(혹은 대부분) 필드를 채워서 직렬화 결과가 예상대로 포함되는지 테스트합니다.
func TestFullCapabilities(t *testing.T) {
	// 아래 preferences 는 Firefox about:config 항목 설정 예시
	preferences := map[string]interface{}{
		"browser.startup.homepage": "https://example.com",
		"javascript.enabled":       false,
	}
	caps := Capabilities{
		Binary:  "/path/to/custom/firefox",
		Args:    []string{"--private-window", "--no-remote"},
		Profile: "someFakeBase64Profile==", // 실제로는 SetProfile 처럼 zip+base64 되어야 합니다.
		Log: &Log{
			Level: Debug,
		},
		Prefs: preferences,
	}
	data, err := json.Marshal(caps)
	if err != nil {
		t.Fatalf("json.Marshal(...) returned error: %v", err)
	}

	got := string(data)
	// 주어진 값들이 모두 JSON 에 들어있는지 검사
	if !contains(got, `"binary":"/path/to/custom/firefox"`) {
		t.Error(`serialized JSON missing or incorrect "binary"`)
	}
	if !contains(got, `"args":["--private-window","--no-remote"]`) {
		t.Error(`serialized JSON missing or incorrect "args"`)
	}
	if !contains(got, `"profile":"someFakeBase64Profile=="`) {
		t.Error(`serialized JSON missing or incorrect "profile"`)
	}
	if !contains(got, `"log":{"level":"debug"}`) {
		t.Error(`serialized JSON missing or incorrect "log"`)
	}
	// Preferences 는 순서가 일정하지 않을 수 있으니 특정 키를 체크
	if !contains(got, `"browser.startup.homepage":"https://example.com"`) {
		t.Error(`preferences missing "browser.startup.homepage"`)
	}
	if !contains(got, `"javascript.enabled":false`) {
		t.Error(`preferences missing "javascript.enabled":false`)
	}
}

// contains 는 문자열 s가 sub 를 포함하는지 단순히 확인하는 헬퍼 함수입니다.
func contains(s, sub string) bool {
	return bytes.Contains([]byte(s), []byte(sub))
}
