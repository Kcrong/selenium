// Package command has WebDriver Commands
package command

// Command represents a WebDriver command
type Command string

// Session commands
const (
	NewSession    Command = "newSession"
	DeleteSession Command = "deleteSession"
	Quit          Command = "quit"
)

// Navigation commands
const (
	GetCurrentURL Command = "getCurrentUrl"
	Get           Command = "get"
	GoBack        Command = "goBack"
	GoForward     Command = "goForward"
	Refresh       Command = "refresh"
	GetTitle      Command = "getTitle"
	GetPageSource Command = "getPageSource"
)

// Element location commands
const (
	FindElement       Command = "findElement"
	FindElements      Command = "findElements"
	FindChildElement  Command = "findChildElement"
	FindChildElements Command = "findChildElements"
	GetActiveElement  Command = "getActiveElement"
)

// Element interaction commands
const (
	ClickElement                 Command = "clickElement"
	ClearElement                 Command = "clearElement"
	SendKeysToElement            Command = "sendKeysToElement"
	SubmitElement                Command = "submitElement"
	GetElementText               Command = "getElementText"
	GetElementTagName            Command = "getElementTagName"
	IsElementSelected            Command = "isElementSelected"
	IsElementEnabled             Command = "isElementEnabled"
	IsElementDisplayed           Command = "isElementDisplayed"
	GetElementLocation           Command = "getElementLocation"
	GetElementSize               Command = "getElementSize"
	GetElementRect               Command = "getElementRect"
	GetElementAttribute          Command = "getElementAttribute"
	GetElementProperty           Command = "getElementProperty"
	GetElementCSSValue           Command = "getElementCSSValue"
	Screenshot                   Command = "screenshot"
	ElementScreenshot            Command = "elementScreenshot"
	GetElementValueOfCssProperty Command = "getElementValueOfCssProperty"
	GetElementAriaRole           Command = "getElementAriaRole"
	GetElementAriaLabel          Command = "getElementAriaLabel"
	W3CGetActiveElement          Command = "w3cGetActiveElement"
)

// Script execution commands
const (
	ExecuteScript         Command = "executeScript"
	ExecuteAsyncScript    Command = "executeAsyncScript"
	W3CExecuteScript      Command = "w3cExecuteScript"
	W3CExecuteScriptAsync Command = "w3cExecuteScriptAsync"
)

// Cookie commands
const (
	GetAllCookies    Command = "getAllCookies"
	GetCookie        Command = "getCookie"
	AddCookie        Command = "addCookie"
	DeleteCookie     Command = "deleteCookie"
	DeleteAllCookies Command = "deleteAllCookies"
)

// Window and frame commands
const (
	SwitchToFrame             Command = "switchToFrame"
	SwitchToParentFrame       Command = "switchToParentFrame"
	SwitchToWindow            Command = "switchToWindow"
	GetWindowHandle           Command = "getWindowHandle"
	GetWindowHandles          Command = "getWindowHandles"
	NewWindow                 Command = "newWindow"
	CloseWindow               Command = "closeWindow"
	PrintPage                 Command = "printPage"
	W3CGetCurrentWindowHandle Command = "w3cGetCurrentWindowHandle"
	W3CGetWindowHandles       Command = "w3cGetWindowHandles"
)

// Screen Orientation commands
const (
	SetScreenOrientation Command = "setScreenOrientation"
	GetScreenOrientation Command = "getScreenOrientation"
)

// Timeout commands
const (
	GetTimeouts Command = "getTimeouts"
	SetTimeouts Command = "setTimeouts"
)

// Alert commands
const (
	AcceptAlert      Command = "acceptAlert"
	DismissAlert     Command = "dismissAlert"
	GetAlertText     Command = "getAlertText"
	SendAlertText    Command = "sendAlertText"
	W3CDismissAlert  Command = "w3cDismissAlert"
	W3CAcceptAlert   Command = "w3cAcceptAlert"
	W3CSetAlertValue Command = "w3cSetAlertValue"
	W3CGetAlertText  Command = "w3cGetAlertText"
)

// Screenshot commands
const (
	TakeScreenshot        Command = "takeScreenshot"
	TakeElementScreenshot Command = "takeElementScreenshot"
)

// Window state commands
const (
	GetWindowRect    Command = "getWindowRect"
	SetWindowRect    Command = "setWindowRect"
	MaximizeWindow   Command = "maximizeWindow"
	MinimizeWindow   Command = "minimizeWindow"
	FullscreenWindow Command = "fullscreenWindow"
)

// Action commands
const (
	W3CActions      Command = "actions"
	W3CClearActions Command = "clearActions"
)

// Web Components commands
const (
	GetShadowRoot              Command = "getShadowRoot"
	FindElementFromShadowRoot  Command = "findElementFromShadowRoot"
	FindElementsFromShadowRoot Command = "findElementsFromShadowRoot"
)

// Virtual Authenticator commands
const (
	AddVirtualAuthenticator    Command = "addVirtualAuthenticator"
	RemoveVirtualAuthenticator Command = "removeVirtualAuthenticator"
	AddCredential              Command = "addCredential"
	GetCredentials             Command = "getCredentials"
	RemoveCredential           Command = "removeCredential"
	RemoveAllCredentials       Command = "removeAllCredentials"
	SetUserVerified            Command = "setUserVerified"
)

// Remote File Management commands
const (
	UploadFile              Command = "uploadFile"
	GetDownloadableFiles    Command = "getDownloadableFiles"
	DownloadFile            Command = "downloadFile"
	DeleteDownloadedFile    Command = "deleteDownloadedFile"
	DeleteDownloadableFiles Command = "deleteDownloadableFiles"
)

// Federated Credential Management (FedCM) commands
const (
	GetFedCMTitle          Command = "getFedcmTitle"
	GetFedCMDialogType     Command = "getFedcmDialogType"
	GetFedCMAccountList    Command = "getFedcmAccountList"
	SelectFedCMAccount     Command = "selectFedcmAccount"
	CancelFedCMDialog      Command = "cancelFedcmDialog"
	SetFedCMDelay          Command = "setFedcmDelay"
	ClickFedCMDialogButton Command = "clickFedcmDialogButton"
	ResetFedCMCooldown     Command = "resetFedcmCooldown"
)

// Mobile commands
const (
	GetNetworkConnection Command = "getNetworkConnection"
	SetNetworkConnection Command = "setNetworkConnection"
	CurrentContextHandle Command = "getCurrentContextHandle"
	ContextHandles       Command = "getContextHandles"
	SwitchToContext      Command = "switchToContext"
)
