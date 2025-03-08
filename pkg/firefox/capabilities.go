// Package firefox provides Firefox-specific types for WebDriver (Selenium 4, W3C).
package firefox

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/Kcrong/selenium/pkg/zip"
)

// CapabilitiesKey is the name of the Firefox-specific key in the W3C
// capabilities object. (Selenium 4 기준)
const CapabilitiesKey = "moz:firefoxOptions"

// Capabilities provides Firefox-specific options for WebDriver (W3C).
// 이 구조체는 "moz:firefoxOptions" 필드에 넣어 사용합니다.
//
// 예:
//
//	firefoxCaps := firefox.Capabilities{
//	    Binary: "/usr/bin/firefox",
//	    Args: []string{"--devtools"},
//	}
//	capabilities := map[string]interface{}{
//	    "alwaysMatch": map[string]interface{}{
//	        "browserName": "firefox",
//	        firefox.CapabilitiesKey: firefoxCaps,
//	    },
//	}
type Capabilities struct {
	// Binary 는 Firefox 실행 파일의 절대 경로입니다.
	// 지정되지 않으면 geckodriver가 시스템에 설치된 Firefox를 자동으로 찾습니다.
	Binary string `json:"binary,omitempty"`

	// Args 는 Firefox 실행 시에 전달할 커맨드 라인 인수 목록입니다.
	// 예) ["--devtools", "--headless"]
	Args []string `json:"args,omitempty"`

	// Profile 은 Base64로 인코딩된 프로필(zip 파일 형태) 데이터입니다.
	// SetProfile 메서드를 통해 기존 디렉토리에서 생성된 프로필을 설정할 수 있습니다.
	Profile string `json:"profile,omitempty"`

	// Log 는 Firefox가 남길 로그 레벨(verbosity)을 지정합니다.
	Log *Log `json:"log,omitempty"`

	// Prefs 는 Firefox의 about:config 항목(설정값)을 키/값으로 넘길 수 있습니다.
	// 값은 string, bool, int 등이 가능합니다.
	Prefs map[string]interface{} `json:"prefs,omitempty"`
}

// SetProfile 는 basePath 디렉토리를 zip으로 묶고, base64 인코딩하여 Profile 필드에 설정합니다.
func (c *Capabilities) SetProfile(basePath string) error {
	buf, err := zip.New(basePath)
	if err != nil {
		return fmt.Errorf("failed to create zip from profile directory: %w", err)
	}
	encoded := new(bytes.Buffer)
	encoded.Grow(buf.Len())
	encoder := base64.NewEncoder(base64.StdEncoding, encoded)
	if _, err := buf.WriteTo(encoder); err != nil {
		return fmt.Errorf("failed to write zip data to base64 encoder: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close base64 encoder: %w", err)
	}

	c.Profile = encoded.String()
	return nil
}

// LogLevel 은 Firefox (geckodriver)의 로깅 레벨을 정의합니다.
// 아래 값들은 geckodriver 에서 지원하는 주요 레벨입니다.
type LogLevel string

const (
	Trace  LogLevel = "trace"
	Debug  LogLevel = "debug"
	Config LogLevel = "config"
	Info   LogLevel = "info"
	Warn   LogLevel = "warn"
	Error  LogLevel = "error"
	Fatal  LogLevel = "fatal"
)

// Log 는 Firefox(geckodriver)가 남길 로그 설정을 지정합니다.
type Log struct {
	// Level 은 로그 레벨을 결정합니다.
	Level LogLevel `json:"level"`
}
