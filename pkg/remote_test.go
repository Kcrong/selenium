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

	"github.com/Kcrong/selenium/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper struct {
	handler func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.handler(req)
}

func TestFindElement(t *testing.T) {
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

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
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
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
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

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
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
	err := wd.SwitchSession("fakeSession")
	require.NoError(t, err)
	elems, err := wd.FindElements(ByName, "myName")
	if err != nil {
		t.Fatalf("FindElements error: %v", err)
	}
	if len(elems) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(elems))
	}
}

func TestAddCookie(t *testing.T) {
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

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
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
	err := wd.SwitchSession("fakeSession")
	require.NoError(t, err)

	c := &Cookie{Name: "myCookie", Value: "val"}
	if err := wd.AddCookie(c); err != nil {
		t.Fatalf("AddCookie error: %v", err)
	}
}

func TestGetCookies(t *testing.T) {
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

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
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
	err := wd.SwitchSession("fakeSession")
	require.NoError(t, err)

	cookies, err := wd.GetCookies()
	if err != nil {
		t.Fatalf("GetCookies error: %v", err)
	}
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
}

func TestDismissAlert(t *testing.T) {
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

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
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
	err := wd.SwitchSession("fakeSession")
	require.NoError(t, err)

	if err := wd.DismissAlert(); err != nil {
		t.Fatalf("DismissAlert error: %v", err)
	}
}

func TestResizeWindow(t *testing.T) {
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

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
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
	err := wd.SwitchSession("fakeSession")
	require.NoError(t, err)

	if err := wd.ResizeWindow("", 800, 600); err != nil {
		t.Fatalf("ResizeWindow error: %v", err)
	}
}

func TestPerformActions(t *testing.T) {
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

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
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
	err := wd.SwitchSession("fakeSession")
	require.NoError(t, err)

	wd.StoreKeyActions("testKey", KeyDownAction("A"), KeyUpAction("A"))
	wd.StorePointerActions("testPointer", MousePointer, PointerDownAction(LeftButton))
	if err := wd.PerformActions(); err != nil {
		t.Fatalf("PerformActions error: %v", err)
	}
}

func TestReleaseActions(t *testing.T) {
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

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
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
	err := wd.SwitchSession("fakeSession")
	require.NoError(t, err)

	if err := wd.ReleaseActions(); err != nil {
		t.Fatalf("ReleaseActions error: %v", err)
	}
}

func TestWaitWithTimeoutAndInterval(t *testing.T) {
	wd := &remoteWD{}
	cond := func(WebDriver) (bool, error) {
		return true, nil
	}
	err := wd.WaitWithTimeoutAndInterval(cond, 2*time.Second, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("WaitWithTimeoutAndInterval error: %v", err)
	}
}

func TestKeyDownUp(t *testing.T) {
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

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
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
	if err := wd.KeyDown("abc"); err != nil {
		t.Fatalf("KeyDown error: %v", err)
	}
	if err := wd.KeyUp("abc"); err != nil {
		t.Fatalf("KeyUp error: %v", err)
	}
}

func TestLog(t *testing.T) {
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

	mock := &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			if req.Method == "POST" && strings.HasSuffix(req.URL.Path, "/log") {
				body := `{"value":[{"timestamp":1690000000000,"level":"INFO","message":"test log"}]}`
				return &http.Response{
					StatusCode: 200,
					Header:     map[string][]string{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(body)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected call: %s %s", req.Method, req.URL.Path)
		},
	}
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
	err := wd.SwitchSession("fakeSession")
	require.NoError(t, err)

	logs, err := wd.Log(log.Browser)
	if err != nil {
		t.Fatalf("Log(Browser) error: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(logs))
	}
	if logs[0].Message != "test log" {
		t.Errorf("expected log message 'test log', got %q", logs[0].Message)
	}
}

func TestRemoteWE_Click(t *testing.T) {
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

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
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
	err := wd.SwitchSession("fakeSession")
	require.NoError(t, err)

	elem := &remoteWE{
		parent: wd,
		id:     "mockElemID",
	}
	if err := elem.Click(); err != nil {
		t.Fatalf("elem.Click() error: %v", err)
	}
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
		if got != tt.want {
			t.Errorf("round(%f) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestScreenshot(t *testing.T) {
	oldClient := HTTPClient
	defer func() { HTTPClient = oldClient }()

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
	HTTPClient = &http.Client{Transport: mock}

	wd := &remoteWD{}
	err := wd.SwitchSession("fakeSession")
	require.NoError(t, err)

	img, err := wd.Screenshot()
	if err != nil {
		t.Fatalf("Screenshot error: %v", err)
	}
	if string(img) != "hello world" {
		t.Errorf("expected screenshot data 'hello world', got %s", string(img))
	}
}

func TestElemMarshalJSON(t *testing.T) {
	elem := &remoteWE{
		id:     "element123",
		parent: &remoteWD{},
	}
	b, err := json.Marshal(elem)
	if err != nil {
		t.Fatalf("Marshal(elem) error: %v", err)
	}
	if !strings.Contains(string(b), `"element-6066-11e4-a52e-4f735466cecf":"element123"`) {
		t.Errorf("expected element-6066...=element123 in JSON, got %s", string(b))
	}
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
