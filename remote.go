// Remote Selenium client implementation.
// See https://www.w3.org/TR/webdriver for the protocol.

package selenium

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/blang/semver"

	"github.com/Kcrong/selenium/log"
)

// chromeCapabilityNames 는 ChromeDriver에서 top-level에 배치하길 기대하는 추가 capability 목록
var chromeCapabilityNames = []string{
	"loggingPrefs",
}

// W3C WebDriver에서 공식적으로 인정되는 top-level capabilities 목록
// (https://www.w3.org/TR/webdriver/#capabilities 참조)
var w3cCapabilityNames = []string{
	"acceptInsecureCerts",
	"browserName",
	"browserVersion",
	"platformName",
	"pageLoadStrategy",
	"proxy",
	"setWindowRect",
	"timeouts",
	"unhandledPromptBehavior",
}

// remoteWD implements WebDriver, holding information about the session.
type remoteWD struct {
	id, urlPrefix  string
	capabilities   Capabilities
	browser        string
	browserVersion semver.Version

	// Actions API: queued KeyActions/PointerActions for Perform/Release.
	storedActions Actions
}

const (
	// DefaultWaitInterval 는 Wait 함수에서 조건을 재확인할 때의 기본 폴링 간격입니다.
	DefaultWaitInterval = 100 * time.Millisecond

	// DefaultWaitTimeout 은 Wait 함수에서 조건 만족을 대기하는 기본 최대 시간입니다.
	DefaultWaitTimeout = 60 * time.Second
)

// HTTPClient is the default client to use to communicate with the WebDriver server.
var HTTPClient = http.DefaultClient

const jsonContentType = "application/json"

// newRequest creates a basic HTTP request object with the given method, URL, and payload.
func newRequest(method, url string, data []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", jsonContentType)
	return req, nil
}

func (wd *remoteWD) requestURL(template string, args ...interface{}) string {
	return wd.urlPrefix + fmt.Sprintf(template, args...)
}

// serverReply is the top-level response from a WebDriver (Selenium 4, W3C).
type serverReply struct {
	// SessionID는 최상위에 있을 수도 있지만, W3C에서는 "value.sessionId"에만 존재하기도 함.
	SessionID *string

	// Value는 W3C 표준에서 실제 응답(또는 에러 정보)을 담는 필드.
	Value json.RawMessage

	// Error는 W3C WebDriver 스펙에서 정의된 에러(에러명, 메시지 등).
	Error
}

// Error는 W3C WebDriver 에러. https://www.w3.org/TR/webdriver/#handling-errors
type Error struct {
	Err        string `json:"error"`
	Message    string `json:"message"`
	Stacktrace string `json:"stacktrace"`
	HTTPCode   int    // 실제 HTTP 상태 코드 저장용
}

// Error()는 error 인터페이스 구현.
func (e *Error) Error() string {
	if e.Err == "" {
		return ""
	}
	return fmt.Sprintf("%s: %s", e.Err, e.Message)
}

// execute는 지정된 method/URL/data로 HTTP 요청을 보낸 뒤, WebDriver 응답(JSON)을 검사한다.
func (wd *remoteWD) execute(method, url string, data []byte) (json.RawMessage, error) {
	return executeCommand(method, url, data)
}

func executeCommand(method, url string, data []byte) (json.RawMessage, error) {
	debugLog("-> %s %s\n%s", method, filteredURL(url), data)

	request, err := newRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	resp, err := HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if debugFlag {
		if err == nil {
			// JSON pretty print
			var prettyBuf bytes.Buffer
			if e := json.Indent(&prettyBuf, buf, "", "    "); e == nil && prettyBuf.Len() > 0 {
				buf = prettyBuf.Bytes()
			}
		}
		debugLog("<- %s [%s]\n%s", resp.Status, resp.Header["Content-Type"], buf)
	}
	if err != nil {
		// Body를 제대로 읽지 못했다면, HTTP 상태 메시지를 반환
		return nil, errors.New(resp.Status)
	}

	fullCType := resp.Header.Get("Content-Type")
	cType, _, e := mime.ParseMediaType(fullCType)
	if e != nil || cType != jsonContentType {
		return nil, fmt.Errorf("got content type %q, expected %q", fullCType, jsonContentType)
	}

	reply := new(serverReply)
	if e := json.Unmarshal(buf, reply); e != nil {
		// HTTP 자체는 200이지만 JSON 파싱 실패 시
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("bad server reply status: %s", resp.Status)
		}
		return nil, e
	}

	// W3C 표준 에러 처리: reply.Err가 존재하면 에러
	if reply.Err != "" {
		reply.HTTPCode = resp.StatusCode
		return nil, &reply.Error
	}

	// 추가로 "value" 안에 W3C 에러 형식이 있을 수 있음
	if len(reply.Value) > 0 {
		possibleErr := new(Error)
		if e := json.Unmarshal(reply.Value, possibleErr); e == nil && possibleErr.Err != "" {
			possibleErr.HTTPCode = resp.StatusCode
			return nil, possibleErr
		}
	}

	// 정상 응답이면 raw JSON payload(buf) 반환
	return buf, nil
}

// DefaultURLPrefix 는 기본 WebDriver endpoint.
const DefaultURLPrefix = "http://127.0.0.1:4444/wd/hub"

// NewRemote 는 새로운 원격(WebDriver) 세션을 생성한다.
func NewRemote(capabilities Capabilities, urlPrefix string) (WebDriver, error) {
	if urlPrefix == "" {
		urlPrefix = DefaultURLPrefix
	}
	wd := &remoteWD{
		urlPrefix:    urlPrefix,
		capabilities: capabilities,
	}

	if b := capabilities["browserName"]; b != nil {
		wd.browser = b.(string)
	}

	if _, err := wd.NewSession(); err != nil {
		return nil, err
	}
	return wd, nil
}

// DeleteSession 는 주어진 세션 ID를 삭제(종료).
func DeleteSession(urlPrefix, id string) error {
	u, err := url.Parse(urlPrefix)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "session", id)
	return voidCommand("DELETE", u.String(), nil)
}

// NewSession 는 wd에 설정된 capabilities 기반으로 새 세션을 만든다.
func (wd *remoteWD) NewSession() (string, error) {
	params := map[string]interface{}{
		"capabilities": map[string]interface{}{
			"alwaysMatch": newW3CCapabilities(wd.capabilities),
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("failed to marshal capabilities: %w", err)
	}

	response, err := wd.execute("POST", wd.requestURL("/session"), data)
	if err != nil {
		return "", fmt.Errorf("failed to create new session: %w", err)
	}

	reply := new(serverReply)
	if err := json.Unmarshal(response, reply); err != nil {
		return "", fmt.Errorf("failed to unmarshal server reply: %w", err)
	}

	// 세션ID 설정
	if reply.SessionID != nil {
		wd.id = *reply.SessionID
	}

	// 반환된 "value"에서 sessionId, capabilities 등을 추가 파싱
	if len(reply.Value) > 0 {
		type returnedCapabilities struct {
			BrowserVersion   string `json:"browserVersion"`
			Version          string `json:"version"`
			PageLoadStrategy string `json:"pageLoadStrategy"`
			Proxy            Proxy  `json:"proxy"`
			Timeouts         struct {
				Implicit float32 `json:"implicit"`
				PageLoad float32 `json:"pageLoad"`
				Script   float32 `json:"script"`
			} `json:"timeouts"`
		}

		value := struct {
			SessionID    string                `json:"sessionId"`
			Capabilities *returnedCapabilities `json:"capabilities"`
		}{}

		if e := json.Unmarshal(reply.Value, &value); e == nil {
			if value.SessionID != "" && wd.id == "" {
				wd.id = value.SessionID
			}
			if value.Capabilities != nil {
				caps := value.Capabilities
				// 브라우저 버전 설정
				for _, ver := range []string{caps.BrowserVersion, caps.Version} {
					if ver == "" {
						continue
					}
					parsed, e := parseVersion(ver)
					if e != nil {
						debugLog("error parsing version: %v\n", e)
						continue
					}
					wd.browserVersion = parsed
					break
				}
			}
		}
	}

	return wd.id, nil
}

// SessionID returns the current session's ID.
func (wd *remoteWD) SessionID() string {
	return wd.id
}

// SwitchSession 는 현재 사용 중인 session ID를 변경한다.
func (wd *remoteWD) SwitchSession(sessionID string) error {
	wd.id = sessionID
	return nil
}

// Capabilities 는 현재 세션의 capabilities를 가져온다.
func (wd *remoteWD) Capabilities() (Capabilities, error) {
	u := wd.requestURL("/session/%s", wd.id)
	resp, err := wd.execute("GET", u, nil)
	if err != nil {
		return nil, err
	}
	c := new(struct{ Value Capabilities })
	if err := json.Unmarshal(resp, c); err != nil {
		return nil, err
	}
	return c.Value, nil
}

// parseVersion 는 브라우저 버전 문자열(예: "91.0.4472.114")을 semver.Version 으로 파싱한다.
func parseVersion(v string) (semver.Version, error) {
	parts := strings.Split(v, ".")
	var err error
	for i := len(parts); i > 0; i-- {
		var ver semver.Version
		ver, err = semver.ParseTolerant(strings.Join(parts[:i], "."))
		if err == nil {
			return ver, nil
		}
	}
	return semver.Version{}, err
}

// W3C 표준에 맞는 capabilities 를 생성한다.
func newW3CCapabilities(caps Capabilities) Capabilities {
	isValidW3CCapability := map[string]bool{}
	for _, name := range w3cCapabilityNames {
		isValidW3CCapability[name] = true
	}
	if b, ok := caps["browserName"]; ok && b == "chrome" {
		for _, name := range chromeCapabilityNames {
			isValidW3CCapability[name] = true
		}
	}

	alwaysMatch := make(Capabilities)
	for name, value := range caps {
		// "moz:" / "goog:" 등 벤더 prefix는 그대로 포함
		if isValidW3CCapability[name] || strings.Contains(name, ":") {
			alwaysMatch[name] = value
		}
	}

	return Capabilities{
		"alwaysMatch": alwaysMatch,
	}
}

// voidCommand 는 요청을 보내고 응답은 무시하는(Body만 확인) 헬퍼.
func voidCommand(method, url string, params interface{}) error {
	if params == nil {
		params = make(map[string]interface{})
	}
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	_, err = executeCommand(method, url, data)
	return err
}

func (wd *remoteWD) voidCommand(urlTemplate string, params interface{}) error {
	return voidCommand("POST", wd.requestURL(urlTemplate, wd.id), params)
}

// stringCommand는 GET으로 요청한 뒤, 응답에서 Value(string)를 추출.
func (wd *remoteWD) stringCommand(urlTemplate string) (string, error) {
	u := wd.requestURL(urlTemplate, wd.id)
	resp, err := wd.execute("GET", u, nil)
	if err != nil {
		return "", err
	}
	reply := new(struct{ Value *string })
	if err := json.Unmarshal(resp, reply); err != nil {
		return "", err
	}
	if reply.Value == nil {
		return "", fmt.Errorf("nil return value")
	}
	return *reply.Value, nil
}

// stringsCommand는 GET으로 요청한 뒤, 응답에서 Value([]string]) 추출.
func (wd *remoteWD) stringsCommand(urlTemplate string) ([]string, error) {
	resp, err := wd.execute("GET", wd.requestURL(urlTemplate, wd.id), nil)
	if err != nil {
		return nil, err
	}
	reply := new(struct{ Value []string })
	if err := json.Unmarshal(resp, reply); err != nil {
		return nil, err
	}
	return reply.Value, nil
}

// SetAsyncScriptTimeout 는 W3C 표준 방식으로 async script timeout 설정.
func (wd *remoteWD) SetAsyncScriptTimeout(timeout time.Duration) error {
	return wd.voidCommand("/session/%s/timeouts", map[string]uint{
		"script": uint(timeout / time.Millisecond),
	})
}

// SetImplicitWaitTimeout 는 W3C 표준 방식으로 implicit wait timeout 설정.
func (wd *remoteWD) SetImplicitWaitTimeout(timeout time.Duration) error {
	return wd.voidCommand("/session/%s/timeouts", map[string]uint{
		"implicit": uint(timeout / time.Millisecond),
	})
}

// SetPageLoadTimeout 는 W3C 표준 방식으로 page load timeout 설정.
func (wd *remoteWD) SetPageLoadTimeout(timeout time.Duration) error {
	return wd.voidCommand("/session/%s/timeouts", map[string]uint{
		"pageLoad": uint(timeout / time.Millisecond),
	})
}

// Quit 는 현재 세션 종료.
func (wd *remoteWD) Quit() error {
	if wd.id == "" {
		return nil
	}
	_, err := wd.execute("DELETE", wd.requestURL("/session/%s", wd.id), nil)
	if err == nil {
		wd.id = ""
	}
	return err
}

func (wd *remoteWD) CurrentWindowHandle() (string, error) {
	return wd.stringCommand("/session/%s/window")
}

func (wd *remoteWD) WindowHandles() ([]string, error) {
	return wd.stringsCommand("/session/%s/window/handles")
}

func (wd *remoteWD) CurrentURL() (string, error) {
	u := wd.requestURL("/session/%s/url", wd.id)
	resp, err := wd.execute("GET", u, nil)
	if err != nil {
		return "", err
	}
	reply := new(struct{ Value *string })
	if err := json.Unmarshal(resp, reply); err != nil {
		return "", err
	}
	return *reply.Value, nil
}

func (wd *remoteWD) Get(url string) error {
	reqURL := wd.requestURL("/session/%s/url", wd.id)
	params := map[string]string{
		"url": url,
	}
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	_, err = wd.execute("POST", reqURL, data)
	return err
}

func (wd *remoteWD) Forward() error {
	return wd.voidCommand("/session/%s/forward", nil)
}

func (wd *remoteWD) Back() error {
	return wd.voidCommand("/session/%s/back", nil)
}

func (wd *remoteWD) Refresh() error {
	return wd.voidCommand("/session/%s/refresh", nil)
}

func (wd *remoteWD) Title() (string, error) {
	return wd.stringCommand("/session/%s/title")
}

func (wd *remoteWD) PageSource() (string, error) {
	return wd.stringCommand("/session/%s/source")
}

// Element-finding helper. W3C에서는 ByCSS, ByXPath, etc.를 사용.
func (wd *remoteWD) find(by, value, suffix, baseURL string) ([]byte, error) {
	// W3C에서는 ByID, ByName 등은 사라졌지만, 필요 시 CSS/XPath로 변환 가능.
	if by == ByID {
		by = ByCSSSelector
		value = "#" + value
	} else if by == ByName {
		by = ByCSSSelector
		value = fmt.Sprintf("input[name=%q]", value)
	}

	params := map[string]string{
		"using": by,
		"value": value,
	}
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if baseURL == "" {
		baseURL = "/session/%s/element"
	}
	return wd.execute("POST", wd.requestURL(baseURL+suffix, wd.id), data)
}

// W3C WebDriver 명세에서 정의된 Unique Element Identifier 키
// https://www.w3.org/TR/webdriver/#elements
const webElementIdentifier = "element-6066-11e4-a52e-4f735466cecf"

// DecodeElement 는 단일 element JSON 응답을 WebElement로 변환.
func (wd *remoteWD) DecodeElement(data []byte) (WebElement, error) {
	reply := new(struct{ Value map[string]string })
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, err
	}

	id := reply.Value[webElementIdentifier] // W3C spec key
	if id == "" {
		return nil, fmt.Errorf("invalid element returned: %+v", reply)
	}
	return &remoteWE{
		parent: wd,
		id:     id,
	}, nil
}

// DecodeElements 는 elements JSON 응답을 WebElement slice로 변환.
func (wd *remoteWD) DecodeElements(data []byte) ([]WebElement, error) {
	reply := new(struct{ Value []map[string]string })
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, err
	}

	elems := make([]WebElement, len(reply.Value))
	for i, elem := range reply.Value {
		id := elem[webElementIdentifier]
		if id == "" {
			return nil, fmt.Errorf("invalid element returned: %+v", elem)
		}
		elems[i] = &remoteWE{
			parent: wd,
			id:     id,
		}
	}
	return elems, nil
}

func (wd *remoteWD) FindElement(by, value string) (WebElement, error) {
	resp, err := wd.find(by, value, "", "")
	if err != nil {
		return nil, err
	}
	return wd.DecodeElement(resp)
}

func (wd *remoteWD) FindElements(by, value string) ([]WebElement, error) {
	resp, err := wd.find(by, value, "s", "")
	if err != nil {
		return nil, err
	}
	return wd.DecodeElements(resp)
}

// Close closes the current window.
func (wd *remoteWD) Close() error {
	url := wd.requestURL("/session/%s/window", wd.id)
	_, err := wd.execute("DELETE", url, nil)
	return err
}

// SwitchWindow 는 주어진 handle로 브라우저 창을 전환.
func (wd *remoteWD) SwitchWindow(handle string) error {
	return wd.voidCommand("/session/%s/window", map[string]string{
		"handle": handle,
	})
}

func (wd *remoteWD) CloseWindow(handle string) error {
	// W3C: window command는 현재 윈도만 닫을 수 있음. handle 무시.
	return voidCommand("DELETE", wd.requestURL("/session/%s/window", wd.id), nil)
}

func (wd *remoteWD) MaximizeWindow(_ string) error {
	url := wd.requestURL("/session/%s/window/maximize", wd.id)
	_, err := wd.execute("POST", url, nil)
	return err
}

func (wd *remoteWD) MinimizeWindow(_ string) error {
	url := wd.requestURL("/session/%s/window/minimize", wd.id)
	_, err := wd.execute("POST", url, nil)
	return err
}

func (wd *remoteWD) ResizeWindow(_ string, width, height int) error {
	url := wd.requestURL("/session/%s/window/rect", wd.id)
	data, err := json.Marshal(map[string]float64{
		"width":  float64(width),
		"height": float64(height),
	})
	if err != nil {
		return err
	}
	_, err = wd.execute("POST", url, data)
	return err
}

// SwitchFrame는 frame 인자로 주어진 요소/인덱스로 프레임 전환.
func (wd *remoteWD) SwitchFrame(frame interface{}) error {
	params := map[string]interface{}{}
	switch f := frame.(type) {
	case nil:
		params["id"] = nil // top-level
	case WebElement:
		params["id"] = f
	case int:
		params["id"] = f
	case string:
		if f == "" {
			params["id"] = nil
		} else {
			// id로 웹엘리먼트 찾기 (임시), CSS 등으로 전환 가능
			elem, err := wd.FindElement(ByID, f)
			if err != nil {
				return err
			}
			params["id"] = elem
		}
	default:
		return fmt.Errorf("invalid frame type %T", frame)
	}
	return wd.voidCommand("/session/%s/frame", params)
}

// ActiveElement returns the current active element.
func (wd *remoteWD) ActiveElement() (WebElement, error) {
	u := wd.requestURL("/session/%s/element/active", wd.id)
	response, err := wd.execute("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return wd.DecodeElement(response)
}

// Cookie / GetCookie / GetCookies / AddCookie / DeleteCookie

func (wd *remoteWD) GetCookie(name string) (Cookie, error) {
	if wd.browser == "chrome" {
		// 일부 드라이버는 GET /cookie/name 없이 /cookie 전체에서 찾는다.
		cs, err := wd.GetCookies()
		if err != nil {
			return Cookie{}, err
		}
		for _, c := range cs {
			if c.Name == name {
				return c, nil
			}
		}
		return Cookie{}, errors.New("cookie not found")
	}
	url := wd.requestURL("/session/%s/cookie/%s", wd.id, name)
	data, err := wd.execute("GET", url, nil)
	if err != nil {
		return Cookie{}, err
	}

	// geckodriver는 단일 cookie 또는 []cookie 반환 가능
	replySingle := new(struct{ Value cookie })
	if e := json.Unmarshal(data, replySingle); e == nil && replySingle.Value.Name != "" {
		return replySingle.Value.sanitize(), nil
	}
	replyList := new(struct{ Value []cookie })
	if e := json.Unmarshal(data, replyList); e != nil {
		return Cookie{}, e
	}
	if len(replyList.Value) == 0 {
		return Cookie{}, errors.New("no cookies returned")
	}
	return replyList.Value[0].sanitize(), nil
}

func (wd *remoteWD) GetCookies() ([]Cookie, error) {
	u := wd.requestURL("/session/%s/cookie", wd.id)
	data, err := wd.execute("GET", u, nil)
	if err != nil {
		return nil, err
	}
	reply := new(struct{ Value []cookie })
	if err := json.Unmarshal(data, reply); err != nil {
		return nil, err
	}
	cookies := make([]Cookie, len(reply.Value))
	for i, c := range reply.Value {
		cookies[i] = c.sanitize()
	}
	return cookies, nil
}

func (wd *remoteWD) AddCookie(c *Cookie) error {
	return wd.voidCommand("/session/%s/cookie", map[string]*Cookie{
		"cookie": c,
	})
}

func (wd *remoteWD) DeleteAllCookies() error {
	u := wd.requestURL("/session/%s/cookie", wd.id)
	_, err := wd.execute("DELETE", u, nil)
	return err
}

func (wd *remoteWD) DeleteCookie(name string) error {
	u := wd.requestURL("/session/%s/cookie/%s", wd.id, name)
	_, err := wd.execute("DELETE", u, nil)
	return err
}

// cookie, sanitize -> Cookie
type cookie struct {
	Name     string      `json:"name"`
	Value    string      `json:"value"`
	Path     string      `json:"path"`
	Domain   string      `json:"domain"`
	Secure   bool        `json:"secure"`
	Expiry   interface{} `json:"expiry"`
	HTTPOnly bool        `json:"httpOnly"`
	SameSite string      `json:"sameSite,omitempty"`
}

func (c cookie) sanitize() Cookie {
	parseExpiry := func(e interface{}) uint {
		switch val := e.(type) {
		case int:
			if val > 0 {
				return uint(val)
			}
		case float64:
			return uint(val)
		}
		return 0
	}
	parseSameSite := func(s string) SameSite {
		if s == "" {
			return ""
		}
		for _, v := range []SameSite{SameSiteNone, SameSiteLax, SameSiteStrict} {
			if strings.EqualFold(string(v), s) {
				return v
			}
		}
		return SameSiteLax
	}
	return Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		Secure:   c.Secure,
		Expiry:   parseExpiry(c.Expiry),
		HTTPOnly: c.HTTPOnly,
		SameSite: parseSameSite(c.SameSite),
	}
}

// Click / DoubleClick / ButtonDown / ButtonUp / (Actions)

// Actions API (W3C):
func (wd *remoteWD) Click(button int) error {
	return wd.voidCommand("/session/%s/click", map[string]int{"button": button})
}

func (wd *remoteWD) DoubleClick() error {
	return wd.voidCommand("/session/%s/doubleclick", nil)
}

func (wd *remoteWD) ButtonDown() error {
	return wd.voidCommand("/session/%s/buttondown", nil)
}

func (wd *remoteWD) ButtonUp() error {
	return wd.voidCommand("/session/%s/buttonup", nil)
}

func (wd *remoteWD) SendModifier(modifier string, isDown bool) error {
	if isDown {
		return wd.KeyDown(modifier)
	}
	return wd.KeyUp(modifier)
}

func (wd *remoteWD) KeyDown(keys string) error {
	return wd.keyAction("keyDown", keys)
}

func (wd *remoteWD) KeyUp(keys string) error {
	return wd.keyAction("keyUp", keys)
}

func (wd *remoteWD) keyAction(action, keys string) error {
	type keyAction struct {
		Type string `json:"type"`
		Key  string `json:"value"`
	}
	actions := make([]keyAction, 0, len(keys))
	for _, k := range keys {
		actions = append(actions, keyAction{
			Type: action,
			Key:  string(k),
		})
	}
	return wd.voidCommand("/session/%s/actions", map[string]interface{}{
		"actions": []interface{}{
			map[string]interface{}{
				"type":    "key",
				"id":      "default keyboard",
				"actions": actions,
			},
		},
	})
}

// KeyPauseAction, KeyUpAction, etc.:
func KeyPauseAction(d time.Duration) KeyAction {
	return KeyAction{"type": "pause", "duration": uint(d / time.Millisecond)}
}
func KeyUpAction(k string) KeyAction {
	return KeyAction{"type": "keyUp", "value": k}
}
func KeyDownAction(k string) KeyAction {
	return KeyAction{"type": "keyDown", "value": k}
}

// PointerPauseAction, PointerMoveAction, etc.:
func PointerPauseAction(d time.Duration) PointerAction {
	return PointerAction{"type": "pause", "duration": uint(d / time.Millisecond)}
}
func PointerMoveAction(d time.Duration, offset Point, origin PointerMoveOrigin) PointerAction {
	return PointerAction{
		"type":     "pointerMove",
		"duration": uint(d / time.Millisecond),
		"origin":   origin,
		"x":        offset.X,
		"y":        offset.Y,
	}
}
func PointerUpAction(button MouseButton) PointerAction {
	return PointerAction{"type": "pointerUp", "button": button}
}
func PointerDownAction(button MouseButton) PointerAction {
	return PointerAction{"type": "pointerDown", "button": button}
}

// StoreKeyActions / StorePointerActions / PerformActions / ReleaseActions

func (wd *remoteWD) StoreKeyActions(inputID string, actions ...KeyAction) {
	raw := []map[string]interface{}{}
	for _, a := range actions {
		raw = append(raw, a)
	}
	wd.storedActions = append(wd.storedActions, map[string]interface{}{
		"type":    "key",
		"id":      inputID,
		"actions": raw,
	})
}

func (wd *remoteWD) StorePointerActions(inputID string, pointer PointerType, actions ...PointerAction) {
	raw := []map[string]interface{}{}
	for _, a := range actions {
		raw = append(raw, a)
	}
	wd.storedActions = append(wd.storedActions, map[string]interface{}{
		"type":       "pointer",
		"id":         inputID,
		"parameters": map[string]string{"pointerType": string(pointer)},
		"actions":    raw,
	})
}

func (wd *remoteWD) PerformActions() error {
	err := wd.voidCommand("/session/%s/actions", map[string]interface{}{
		"actions": wd.storedActions,
	})
	wd.storedActions = nil
	return err
}

func (wd *remoteWD) ReleaseActions() error {
	return voidCommand("DELETE", wd.requestURL("/session/%s/actions", wd.id), nil)
}

// Alerts
func (wd *remoteWD) DismissAlert() error {
	return wd.voidCommand("/session/%s/alert/dismiss", nil)
}
func (wd *remoteWD) AcceptAlert() error {
	return wd.voidCommand("/session/%s/alert/accept", nil)
}
func (wd *remoteWD) AlertText() (string, error) {
	return wd.stringCommand("/session/%s/alert/text")
}
func (wd *remoteWD) SetAlertText(text string) error {
	return wd.voidCommand("/session/%s/alert/text", map[string]string{"text": text})
}

// execScript / execScriptAsync / raw variants
func (wd *remoteWD) execScriptRaw(script string, args []interface{}, suffix string) ([]byte, error) {
	if args == nil {
		args = []interface{}{}
	}
	data, err := json.Marshal(map[string]interface{}{
		"script": script,
		"args":   args,
	})
	if err != nil {
		return nil, err
	}
	return wd.execute("POST", wd.requestURL("/session/%s/execute"+suffix, wd.id), data)
}

func (wd *remoteWD) execScript(script string, args []interface{}, suffix string) (interface{}, error) {
	resp, err := wd.execScriptRaw(script, args, suffix)
	if err != nil {
		return nil, err
	}
	reply := new(struct{ Value interface{} })
	if err = json.Unmarshal(resp, reply); err != nil {
		return nil, err
	}
	return reply.Value, nil
}

func (wd *remoteWD) ExecuteScript(script string, args []interface{}) (interface{}, error) {
	return wd.execScript(script, args, "/sync")
}

func (wd *remoteWD) ExecuteScriptAsync(script string, args []interface{}) (interface{}, error) {
	return wd.execScript(script, args, "/async")
}

func (wd *remoteWD) ExecuteScriptRaw(script string, args []interface{}) ([]byte, error) {
	return wd.execScriptRaw(script, args, "/sync")
}

func (wd *remoteWD) ExecuteScriptAsyncRaw(script string, args []interface{}) ([]byte, error) {
	return wd.execScriptRaw(script, args, "/async")
}

// Screenshot captures a screenshot of the current page.
func (wd *remoteWD) Screenshot() ([]byte, error) {
	data, err := wd.stringCommand("/session/%s/screenshot")
	if err != nil {
		return nil, err
	}
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(data))
	return io.ReadAll(decoder)
}

// Condition 은 WebDriver 상태를 검사하는 함수 타입입니다.
// 반환값 (true, nil)이면 조건을 만족했고 대기(Wait) 함수를 종료할 수 있음을 의미합니다.
// (false, nil)이면 아직 조건이 만족되지 않아 재시도가 필요함을 의미합니다.
// 에러를 반환하면 대기를 즉시 중단하고 해당 에러가 반환됩니다.
type Condition func(wd WebDriver) (bool, error)

func (wd *remoteWD) WaitWithTimeoutAndInterval(cond Condition, timeout, interval time.Duration) error {
	startTime := time.Now()
	for {
		ok, err := cond(wd)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		if time.Since(startTime) > timeout {
			return fmt.Errorf("timeout after %v", time.Since(startTime))
		}
		time.Sleep(interval)
	}
}

func (wd *remoteWD) WaitWithTimeout(cond Condition, timeout time.Duration) error {
	return wd.WaitWithTimeoutAndInterval(cond, timeout, DefaultWaitInterval)
}

func (wd *remoteWD) Wait(cond Condition) error {
	return wd.WaitWithTimeoutAndInterval(cond, DefaultWaitTimeout, DefaultWaitInterval)
}

// Log retrieve logs (e.g., browser, driver, performance, etc.)
func (wd *remoteWD) Log(typ log.Type) ([]log.Message, error) {
	u := wd.requestURL("/session/%s/log", wd.id)
	params := map[string]log.Type{"type": typ}
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	resp, err := wd.execute("POST", u, data)
	if err != nil {
		return nil, err
	}
	c := new(struct {
		Value []struct {
			Timestamp int64  `json:"timestamp"`
			Level     string `json:"level"`
			Message   string `json:"message"`
		} `json:"value"`
	})
	if err = json.Unmarshal(resp, c); err != nil {
		return nil, err
	}
	res := make([]log.Message, len(c.Value))
	for i, v := range c.Value {
		res[i] = log.Message{
			Timestamp: time.Unix(0, v.Timestamp*int64(time.Millisecond)),
			Level:     log.Level(v.Level),
			Message:   v.Message,
		}
	}
	return res, nil
}

// remoteWE represents a WebElement, storing parent session and element id.
type remoteWE struct {
	parent *remoteWD
	id     string
}

// Click / SendKeys / etc.

func (elem *remoteWE) Click() error {
	return elem.parent.voidCommand(fmt.Sprintf("/session/%%s/element/%s/click", elem.id), nil)
}

func (elem *remoteWE) SendKeys(keys string) error {
	params := map[string]string{"text": keys}
	urlTemplate := fmt.Sprintf("/session/%%s/element/%s/value", elem.id)
	return elem.parent.voidCommand(urlTemplate, params)
}

func (elem *remoteWE) TagName() (string, error) {
	urlTemplate := fmt.Sprintf("/session/%%s/element/%s/name", elem.id)
	return elem.parent.stringCommand(urlTemplate)
}

func (elem *remoteWE) Text() (string, error) {
	urlTemplate := fmt.Sprintf("/session/%%s/element/%s/text", elem.id)
	return elem.parent.stringCommand(urlTemplate)
}

func (elem *remoteWE) Submit() error {
	urlTemplate := fmt.Sprintf("/session/%%s/element/%s/submit", elem.id)
	return elem.parent.voidCommand(urlTemplate, nil)
}

func (elem *remoteWE) Clear() error {
	urlTemplate := fmt.Sprintf("/session/%%s/element/%s/clear", elem.id)
	return elem.parent.voidCommand(urlTemplate, nil)
}

func (elem *remoteWE) MoveTo(xOffset, yOffset int) error {
	return elem.parent.voidCommand("/session/%s/moveto", map[string]interface{}{
		"element": elem.id,
		"xoffset": xOffset,
		"yoffset": yOffset,
	})
}

// FindElement / FindElements (within element)
func (elem *remoteWE) FindElement(by, value string) (WebElement, error) {
	url := fmt.Sprintf("/session/%%s/element/%s/element", elem.id)
	resp, err := elem.parent.find(by, value, "", url)
	if err != nil {
		return nil, err
	}
	return elem.parent.DecodeElement(resp)
}

func (elem *remoteWE) FindElements(by, value string) ([]WebElement, error) {
	url := fmt.Sprintf("/session/%%s/element/%s/element", elem.id)
	resp, err := elem.parent.find(by, value, "s", url)
	if err != nil {
		return nil, err
	}
	return elem.parent.DecodeElements(resp)
}

// boolQuery utility
func (elem *remoteWE) boolQuery(urlTemplate string) (bool, error) {
	return elem.parent.boolCommand(fmt.Sprintf(urlTemplate, elem.id))
}

// boolCommand fetches Value(bool)
func (wd *remoteWD) boolCommand(urlTemplate string) (bool, error) {
	resp, err := wd.execute("GET", wd.requestURL(urlTemplate, wd.id), nil)
	if err != nil {
		return false, err
	}
	reply := new(struct{ Value bool })
	if err := json.Unmarshal(resp, reply); err != nil {
		return false, err
	}
	return reply.Value, nil
}

// IsSelected / IsEnabled / IsDisplayed
func (elem *remoteWE) IsSelected() (bool, error) {
	return elem.boolQuery("/session/%%s/element/%s/selected")
}
func (elem *remoteWE) IsEnabled() (bool, error) {
	return elem.boolQuery("/session/%%s/element/%s/enabled")
}
func (elem *remoteWE) IsDisplayed() (bool, error) {
	return elem.boolQuery("/session/%%s/element/%s/displayed")
}

// GetProperty / GetAttribute
func (elem *remoteWE) GetProperty(name string) (string, error) {
	urlT := fmt.Sprintf("/session/%%s/element/%s/property/%s", elem.id, name)
	return elem.parent.stringCommand(urlT)
}

func (elem *remoteWE) GetAttribute(name string) (string, error) {
	urlT := fmt.Sprintf("/session/%%s/element/%s/attribute/%s", elem.id, name)
	return elem.parent.stringCommand(urlT)
}

// Location / LocationInView / Size
func (elem *remoteWE) Location() (*Point, error) {
	return elem.location("")
}
func (elem *remoteWE) LocationInView() (*Point, error) {
	return elem.location("_in_view")
}

func (elem *remoteWE) location(suffix string) (*Point, error) {
	// W3C: element rect 사용
	r, err := elem.rect()
	if err != nil {
		return nil, err
	}
	return &Point{round(r.X), round(r.Y)}, nil
}

func (elem *remoteWE) Size() (*Size, error) {
	r, err := elem.rect()
	if err != nil {
		return nil, err
	}
	return &Size{round(r.Width), round(r.Height)}, nil
}

type rect struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// rect implements "Get Element Rect" W3C.
func (elem *remoteWE) rect() (*rect, error) {
	u := elem.parent.requestURL("/session/%s/element/%s/rect", elem.parent.id, elem.id)
	resp, err := elem.parent.execute("GET", u, nil)
	if err != nil {
		return nil, err
	}
	r := new(struct{ Value rect })
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}
	return &r.Value, nil
}

func round(f float64) int {
	if f < 0 {
		return int(f - 0.5)
	}
	return int(f + 0.5)
}

// CSSProperty returns the computed style property of the element.
func (elem *remoteWE) CSSProperty(name string) (string, error) {
	urlT := fmt.Sprintf("/session/%%s/element/%s/css/%s", elem.id, name)
	return elem.parent.stringCommand(urlT)
}

// MarshalJSON W3C 에서는 "element-6066-11e4-a52e-4f735466cecf" 키만 사용.
func (elem *remoteWE) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		webElementIdentifier: elem.id,
	})
}

// Screenshot captures element-specific screenshot (W3C).
func (elem *remoteWE) Screenshot(_ bool) ([]byte, error) {
	urlT := fmt.Sprintf("/session/%%s/element/%s/screenshot", elem.id)
	data, err := elem.parent.stringCommand(urlT)
	if err != nil {
		return nil, err
	}
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(data))
	return io.ReadAll(decoder)
}
