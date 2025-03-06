package chrome

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"testing"
)

// TestEmptyCapabilities 는 빈 Capabilities 구조체를 json.Marshal 했을 때 "{}" 가 나오는지 확인합니다.
func TestEmptyCapabilities(t *testing.T) {
	data, err := json.Marshal(Capabilities{})
	if err != nil {
		t.Fatalf("json.Marshal(Capabilities{}) returned error: %v", err)
	}
	got, want := string(data), "{}"
	if got != want {
		t.Fatalf("json.Marshal(Capabilities{}) = %q, want %q", got, want)
	}
}

// TestPartialCapabilities 는 몇 가지 필드를 채워넣은 경우 직렬화 결과를 테스트합니다.
func TestPartialCapabilities(t *testing.T) {
	caps := Capabilities{
		Path: "/path/to/chrome",
		Args: []string{"--headless", "--disable-gpu"},
	}
	data, err := json.Marshal(caps)
	if err != nil {
		t.Fatalf("json.Marshal(...) returned error: %v", err)
	}

	// omitempty 옵션이 적용된 필드는 비어 있으면 직렬화되지 않습니다.
	// 따라서 아래처럼 Path(binary), Args만 들어간 JSON이 생성될 것으로 기대합니다.
	got := string(data)
	want := `{"binary":"/path/to/chrome","args":["--headless","--disable-gpu"]}`
	if got != want {
		t.Fatalf("json.Marshal(...) = %q, want %q", got, want)
	}
}

// TestAddExtension 는 Capabilities.AddExtension 로 확장(.crx) 파일을 추가했을 때,
// base64 인코딩된 문자열이 Extensions 필드에 잘 들어가는지 확인합니다.
func TestAddExtension(t *testing.T) {
	// 가상의 확장 데이터 (일반적으로 .crx 파일)
	fakeData := []byte("this is a fake extension data")

	caps := Capabilities{}
	// addExtension는 io.Reader를 받아 base64 인코딩하므로, bytes.NewReader로 가상 데이터를 제공합니다.
	err := caps.addExtension(bytes.NewReader(fakeData))
	if err != nil {
		t.Fatalf("addExtension returned error: %v", err)
	}

	if len(caps.Extensions) != 1 {
		t.Fatalf("expected 1 extension in caps.Extensions, got %d", len(caps.Extensions))
	}

	encoded := caps.Extensions[0]
	want := base64.StdEncoding.EncodeToString(fakeData)
	if encoded != want {
		t.Errorf("caps.Extensions[0] = %q, want %q", encoded, want)
	}
}

// TestFullCapabilities 는 가능한 많은 필드를 채운 후 JSON 직렬화가 잘 되는지,
// 그리고 omitempty로 인해 비어있지 않은 필드가 모두 포함되는지 확인하는 예시입니다.
func TestFullCapabilities(t *testing.T) {
	detach := true
	mobile := &MobileEmulation{
		DeviceName: "Pixel 5",
		DeviceMetrics: &DeviceMetrics{
			Width:      1080,
			Height:     1920,
			PixelRatio: 3.0,
		},
		UserAgent: "CustomUserAgent",
	}
	prefs := map[string]interface{}{
		"homepage": "https://example.com",
	}
	perf := &PerfLoggingPreferences{
		EnableNetwork:  boolPtr(true),
		EnablePage:     boolPtr(false),
		EnableTimeline: boolPtr(true),
	}

	caps := Capabilities{
		Path:             "/usr/bin/google-chrome",
		Args:             []string{"--incognito", "--disable-plugins"},
		ExcludeSwitches:  []string{"enable-automation"},
		Detach:           &detach,
		DebuggerAddr:     "127.0.0.1:9222",
		MinidumpPath:     "/tmp/minidumps",
		MobileEmulation:  mobile,
		Prefs:            prefs,
		PerfLoggingPrefs: perf,
		WindowTypes:      []string{"webview", "window"},
		AndroidPackage:   "com.android.chrome",
	}

	data, err := json.Marshal(caps)
	if err != nil {
		t.Fatalf("json.Marshal(...) returned error: %v", err)
	}

	got := string(data)
	// 간단히 특정 필드가 있는지만 확인합니다(정확한 문자열 비교 대신 Contains 등을 써도 됨).
	if !containsAll(got,
		`"binary":"/usr/bin/google-chrome"`,
		`"args":["--incognito","--disable-plugins"]`,
		`"excludeSwitches":["enable-automation"]`,
		`"detach":true`,
		`"debuggerAddress":"127.0.0.1:9222"`,
		`"minidumpPath":"/tmp/minidumps"`,
		`"mobileEmulation":`,
		`"prefs":{"homepage":"https://example.com"}`,
		`"perfLoggingPrefs":`,
		`"windowTypes":["webview","window"]`,
		`"androidPackage":"com.android.chrome"`,
	) {
		t.Errorf("Serialized JSON is missing some expected fields.\nGot: %s", got)
	}
}

// boolPtr 는 bool 값을 포인터로 변환해주는 헬퍼 함수입니다.
func boolPtr(b bool) *bool {
	return &b
}

// containsAll 는 문자열 s가 targets 안의 모든 서브스트링을 포함하는지 검사합니다.
func containsAll(s string, targets ...string) bool {
	for _, t := range targets {
		if !contains(s, t) {
			return false
		}
	}
	return true
}

// 간단한 문자열 검색 헬퍼 (Go 버전별로 strings.Contains 등을 활용)
func contains(s, sub string) bool {
	return bytes.Contains([]byte(s), []byte(sub))
}
