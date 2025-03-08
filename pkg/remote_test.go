package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Kcrong/selenium/pkg/log"
)

type mockRoundTripper struct {
	handler func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.handler(req)
}

func TestFindElement(t *testing.T) {
	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			if req.Method == "POST" && strings.Contains(req.URL.Path, "/element") {
				body := `{"value":{"element-6066-11e4-a52e-4f735466cecf":"mockElementID"}}`
				return &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}
			return nil, errors.New("unexpected request")
		},
	}

	wd := &remoteWD{}
	wd.httpClient = &http.Client{Transport: mock}
	err := wd.SwitchSession("fakeSession")
	require.NoError(t, err)
	elem, err := wd.FindElement(ByID, "myID")
	if err != nil {
		t.Fatalf("FindElement error: %v", err)
	}
	if elem == nil {
		t.Fatal("FindElement returned nil element")
	}
}

func TestFindElements(t *testing.T) {
	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			if req.Method == "POST" && strings.HasSuffix(req.URL.Path, "/elements") {
				body := `{"value":[{"element-6066-11e4-a52e-4f735466cecf":"id1"},{"element-6066-11e4-a52e-4f735466cecf":"id2"}]}`
				return &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}
			return nil, errors.New("unexpected request")
		},
	}

	wd := &remoteWD{httpClient: &http.Client{Transport: mock}}

	require.NoError(t, wd.SwitchSession("fakeSession"))
	elems, err := wd.FindElements(ByName, "myName")
	if err != nil {
		t.Fatalf("FindElements error: %v", err)
	}
	if len(elems) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(elems))
	}
}

func TestAddCookie(t *testing.T) {
	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			if req.Method == "POST" && strings.Contains(req.URL.Path, "/cookie") {
				body := `{"value":{}}`
				return &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected path: %s", req.URL.Path)
		},
	}

	wd := &remoteWD{httpClient: &http.Client{Transport: mock}}
	require.NoError(t, wd.SwitchSession("fakeSession"))

	c := &Cookie{Name: "myCookie", Value: "val"}
	if err := wd.AddCookie(c); err != nil {
		t.Fatalf("AddCookie error: %v", err)
	}
}

func TestGetCookies(t *testing.T) {
	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			if req.Method == "GET" && strings.Contains(req.URL.Path, "/cookie") {
				body := `{"value":[{"name":"test","value":"123","path":"/","domain":"example.com"}]}`
				return &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}
			return nil, errors.New("unexpected request")
		},
	}

	wd := &remoteWD{httpClient: &http.Client{Transport: mock}}
	require.NoError(t, wd.SwitchSession("fakeSession"))

	cookies, err := wd.GetCookies()
	if err != nil {
		t.Fatalf("GetCookies error: %v", err)
	}
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
}

func TestDismissAlert(t *testing.T) {
	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			if req.Method == "POST" && strings.HasSuffix(req.URL.Path, "/alert/dismiss") {
				body := `{"value":{}}`
				return &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected path: %s", req.URL.Path)
		},
	}

	wd := &remoteWD{httpClient: &http.Client{Transport: mock}}
	require.NoError(t, wd.SwitchSession("fakeSession"))

	if err := wd.DismissAlert(); err != nil {
		t.Fatalf("DismissAlert error: %v", err)
	}
}

func TestResizeWindow(t *testing.T) {
	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			if req.Method == "POST" && strings.HasSuffix(req.URL.Path, "/window/rect") {
				body := `{"value":{"x":0,"y":0,"width":800,"height":600}}`
				return &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected path: %s", req.URL.Path)
		},
	}

	wd := &remoteWD{httpClient: &http.Client{Transport: mock}}
	require.NoError(t, wd.SwitchSession("fakeSession"))

	if err := wd.ResizeWindow("", 800, 600); err != nil {
		t.Fatalf("ResizeWindow error: %v", err)
	}
}

func TestPerformActions(t *testing.T) {
	testCases := []struct {
		name   string
		button MouseButton
	}{
		{"Left Button", LeftButton},
		{"Middle Button", MiddleButton},
		{"Right Button", RightButton},
	}

	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			if req.Method == "POST" && strings.HasSuffix(req.URL.Path, "/actions") {
				body := `{"value":{}}`
				return &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected request to %s", req.URL.Path)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wd := &remoteWD{httpClient: &http.Client{Transport: mock}}
			require.NoError(t, wd.SwitchSession("fakeSession"))

			wd.StoreKeyActions("testKey", KeyDownAction("A"), KeyUpAction("A"))
			wd.StorePointerActions("testPointer", MousePointer, PointerDownAction(tc.button), PointerUpAction(tc.button))

			require.NoError(t, wd.PerformActions(), "PerformActions should not return an error")
		})
	}
}

func TestReleaseActions(t *testing.T) {
	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			if req.Method == "DELETE" && strings.HasSuffix(req.URL.Path, "/actions") {
				body := `{"value":{}}`
				return &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
		},
	}

	wd := &remoteWD{httpClient: &http.Client{Transport: mock}}
	require.NoError(t, wd.SwitchSession("fakeSession"))

	assert.NoError(t, wd.ReleaseActions(), "ReleaseActions should not return an error")
}

func TestWaitWithTimeoutAndInterval(t *testing.T) {
	wd := &remoteWD{}
	cond := func(WebDriver) (bool, error) {
		return true, nil
	}
	err := wd.WaitWithTimeoutAndInterval(cond, 2*time.Second, 10*time.Millisecond)
	assert.NoError(t, err, "WaitWithTimeoutAndInterval should not return an error")
}

func TestKeyDownUp(t *testing.T) {
	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			if req.Method == "POST" && strings.HasSuffix(req.URL.Path, "/actions") {
				body := `{"value":{}}`
				return &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected actions call: %s", req.URL.Path)
		},
	}

	wd := &remoteWD{httpClient: &http.Client{Transport: mock}}

	assert.NoError(t, wd.KeyDown("abc"))
	assert.NoError(t, wd.KeyUp("abc"))
}

func TestLog(t *testing.T) {
	testCases := []struct {
		name     string
		logType  log.Type
		response string
		wantErr  bool
		wantLogs []log.Message
	}{
		{
			name:     "Browser Log",
			logType:  log.Browser,
			response: `{"value":[{"timestamp":1690000000000,"level":"INFO","message":"browser log"}]}`,
			wantErr:  false,
			wantLogs: []log.Message{
				{Timestamp: time.Unix(1690000000, 0), Level: log.Level("INFO"), Message: "browser log"},
			},
		},
		{
			name:     "Server Log",
			logType:  log.Server,
			response: `{"value":[{"timestamp":1690000001000,"level":"WARNING","message":"server log"}]}`,
			wantErr:  false,
			wantLogs: []log.Message{
				{Timestamp: time.Unix(1690000001, 0), Level: log.Level("WARNING"), Message: "server log"},
			},
		},
		{
			name:     "Client Log",
			logType:  log.Client,
			response: `{"value":[{"timestamp":1690000002000,"level":"ERROR","message":"client log"}]}`,
			wantErr:  false,
			wantLogs: []log.Message{
				{Timestamp: time.Unix(1690000002, 0), Level: log.Level("ERROR"), Message: "client log"},
			},
		},
		{
			name:     "Driver Log",
			logType:  log.Driver,
			response: `{"value":[{"timestamp":1690000003000,"level":"DEBUG","message":"driver log"}]}`,
			wantErr:  false,
			wantLogs: []log.Message{
				{Timestamp: time.Unix(1690000003, 0), Level: log.Level("DEBUG"), Message: "driver log"},
			},
		},
		{
			name:     "Performance Log",
			logType:  log.Performance,
			response: `{"value":[{"timestamp":1690000004000,"level":"INFO","message":"performance log"}]}`,
			wantErr:  false,
			wantLogs: []log.Message{
				{Timestamp: time.Unix(1690000004, 0), Level: log.Level("INFO"), Message: "performance log"},
			},
		},
		{
			name:     "Profiler Log",
			logType:  log.Profiler,
			response: `{"value":[{"timestamp":1690000005000,"level":"INFO","message":"profiler log"}]}`,
			wantErr:  false,
			wantLogs: []log.Message{
				{Timestamp: time.Unix(1690000005, 0), Level: log.Level("INFO"), Message: "profiler log"},
			},
		},
		{
			name:     "Invalid JSON Response",
			logType:  log.Browser,
			response: `{invalid": "data"}`, // 잘못된 JSON 구조
			wantErr:  true,
			wantLogs: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockRoundTripper{
				handler: func(req *http.Request) (*http.Response, error) {
					if req.Method == "POST" && strings.HasSuffix(req.URL.Path, "/log") {
						return &http.Response{
							StatusCode: 200,
							Header:     map[string][]string{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(bytes.NewBufferString(tc.response)),
						}, nil
					}
					return nil, fmt.Errorf("unexpected call: %s %s", req.Method, req.URL.Path)
				},
			}

			wd := &remoteWD{httpClient: &http.Client{Transport: mock}}
			require.NoError(t, wd.SwitchSession("fakeSession"))

			logs, err := wd.Log(tc.logType)
			if tc.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error in Log()")
				assert.Equal(t, len(tc.wantLogs), len(logs), "unexpected number of log entries")
				for i, wantLog := range tc.wantLogs {
					assert.Equal(t, wantLog.Message, logs[i].Message, "unexpected log message")
					assert.Equal(t, wantLog.Level, logs[i].Level, "unexpected log level")
					assert.Equal(t, wantLog.Timestamp.Unix(), logs[i].Timestamp.Unix(), "unexpected timestamp")
				}
			}
		})
	}
}

func TestRemoteWE_Click(t *testing.T) {
	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			if req.Method == "POST" && strings.Contains(req.URL.Path, "/element/mockElemID/click") {
				body := `{"value":{}}`
				return &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
		},
	}

	wd := &remoteWD{httpClient: &http.Client{Transport: mock}}
	require.NoError(t, wd.SwitchSession("fakeSession"))

	elem := &remoteWE{
		parent: wd,
		id:     "mockElemID",
	}
	assert.NoError(t, elem.Click(), "Click() should not return an error")
}

func TestRound(t *testing.T) {
	tests := []struct {
		in   float64
		want int
	}{
		{0.4, 0},
		{0.5, 1},
		{1.49, 1},
		{-0.4, 0},
		{-0.6, -1},
	}
	for _, tt := range tests {
		got := round(tt.in)
		assert.Equal(t, tt.want, got)
	}
}

func TestScreenshot(t *testing.T) {
	dataBase64 := "aGVsbG8gd29ybGQ=" // "hello world"
	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":"%s"}`, dataBase64)
			return &http.Response{
				StatusCode: 200,
				Header:     map[string][]string{"Content-Type": {"application/json"}},
				Body:       io.NopCloser(bytes.NewBufferString(body)),
			}, nil
		},
	}

	wd := &remoteWD{httpClient: &http.Client{Transport: mock}}
	require.NoError(t, wd.SwitchSession("fakeSession"))

	img, err := wd.Screenshot()
	require.NoError(t, err, "Screenshot() should not return an error")
	assert.EqualValues(t, "hello world", string(img))
}

func TestElemMarshalJSON(t *testing.T) {
	elem := &remoteWE{
		id:     "element123",
		parent: &remoteWD{},
	}
	b, err := json.Marshal(elem)
	require.NoError(t, err, "MarshalJSON() should not return an error")
	assert.Contains(t, string(b), `"element-6066-11e4-a52e-4f735466cecf":"element123"`)
}

func TestVersionParse(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
		major uint64
	}{
		{"99.0", true, 99},
		{"101.10.55", true, 101},
		{"broken", false, 0},
	}
	for _, tt := range tests {
		v, err := parseVersion(tt.input)
		require.Equal(t, tt.ok, err == nil)
		assert.Equal(t, tt.major, v.Major)
	}
}
