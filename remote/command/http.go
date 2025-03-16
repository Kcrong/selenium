package command

import (
	"net/http"
)

// Endpoint represents the HTTP method and path for a command
type Endpoint struct {
	Method string
	Path   string
}

type EndPointMapType map[Command]Endpoint

// EndpointMap maps commands to their HTTP method and path
var EndpointMap = EndPointMapType{
	NewSession:    {http.MethodPost, "/session"},
	DeleteSession: {http.MethodDelete, "/session/$sessionId"},
	Quit:          {http.MethodDelete, "/session/$sessionId"},

	GetCurrentURL: {http.MethodGet, "/session/$sessionId/url"},
	Get:           {http.MethodPost, "/session/$sessionId/url"},
	GoBack:        {http.MethodPost, "/session/$sessionId/back"},
	GoForward:     {http.MethodPost, "/session/$sessionId/forward"},
	Refresh:       {http.MethodPost, "/session/$sessionId/refresh"},
	GetTitle:      {http.MethodGet, "/session/$sessionId/title"},
	GetPageSource: {http.MethodGet, "/session/$sessionId/source"},

	FindElement:       {http.MethodPost, "/session/$sessionId/element"},
	FindElements:      {http.MethodPost, "/session/$sessionId/elements"},
	FindChildElement:  {http.MethodPost, "/session/$sessionId/element/$id/element"},
	FindChildElements: {http.MethodPost, "/session/$sessionId/element/$id/elements"},
	GetActiveElement:  {http.MethodGet, "/session/$sessionId/element/active"},

	ClickElement:                 {http.MethodPost, "/session/$sessionId/element/$id/click"},
	ClearElement:                 {http.MethodPost, "/session/$sessionId/element/$id/clear"},
	SendKeysToElement:            {http.MethodPost, "/session/$sessionId/element/$id/value"},
	SubmitElement:                {http.MethodPost, "/session/$sessionId/element/$id/submit"},
	GetElementText:               {http.MethodGet, "/session/$sessionId/element/$id/text"},
	GetElementTagName:            {http.MethodGet, "/session/$sessionId/element/$id/name"},
	IsElementSelected:            {http.MethodGet, "/session/$sessionId/element/$id/selected"},
	IsElementEnabled:             {http.MethodGet, "/session/$sessionId/element/$id/enabled"},
	IsElementDisplayed:           {http.MethodGet, "/session/$sessionId/element/$id/displayed"},
	GetElementLocation:           {http.MethodGet, "/session/$sessionId/element/$id/location"},
	GetElementSize:               {http.MethodGet, "/session/$sessionId/element/$id/size"},
	GetElementRect:               {http.MethodGet, "/session/$sessionId/element/$id/rect"},
	GetElementAttribute:          {http.MethodGet, "/session/$sessionId/element/$id/attribute/$name"},
	GetElementProperty:           {http.MethodGet, "/session/$sessionId/element/$id/property/$name"},
	GetElementCSSValue:           {http.MethodGet, "/session/$sessionId/element/$id/css/$propertyName"},
	Screenshot:                   {http.MethodGet, "/session/$sessionId/screenshot"},
	ElementScreenshot:            {http.MethodGet, "/session/$sessionId/element/$id/screenshot"},
	GetElementValueOfCssProperty: {http.MethodGet, "/session/$sessionId/element/$id/css/$propertyName"},
	GetElementAriaRole:           {http.MethodGet, "/session/$sessionId/element/$id/computedrole"},
	GetElementAriaLabel:          {http.MethodGet, "/session/$sessionId/element/$id/computedlabel"},
	W3CGetActiveElement:          {http.MethodGet, "/session/$sessionId/element/active"},

	ExecuteScript:         {http.MethodPost, "/session/$sessionId/execute/sync"},
	ExecuteAsyncScript:    {http.MethodPost, "/session/$sessionId/execute/async"},
	W3CExecuteScript:      {http.MethodPost, "/session/$sessionId/execute/sync"},
	W3CExecuteScriptAsync: {http.MethodPost, "/session/$sessionId/execute/async"},

	GetAllCookies:    {http.MethodGet, "/session/$sessionId/cookie"},
	GetCookie:        {http.MethodGet, "/session/$sessionId/cookie/$name"},
	AddCookie:        {http.MethodPost, "/session/$sessionId/cookie"},
	DeleteCookie:     {http.MethodDelete, "/session/$sessionId/cookie/$name"},
	DeleteAllCookies: {http.MethodDelete, "/session/$sessionId/cookie"},

	SwitchToFrame:             {http.MethodPost, "/session/$sessionId/frame"},
	SwitchToParentFrame:       {http.MethodPost, "/session/$sessionId/frame/parent"},
	SwitchToWindow:            {http.MethodPost, "/session/$sessionId/window"},
	GetWindowHandle:           {http.MethodGet, "/session/$sessionId/window"},
	GetWindowHandles:          {http.MethodGet, "/session/$sessionId/window/handles"},
	NewWindow:                 {http.MethodPost, "/session/$sessionId/window/new"},
	CloseWindow:               {http.MethodDelete, "/session/$sessionId/window"},
	PrintPage:                 {http.MethodPost, "/session/$sessionId/print"},
	W3CGetCurrentWindowHandle: {http.MethodGet, "/session/$sessionId/window"},
	W3CGetWindowHandles:       {http.MethodGet, "/session/$sessionId/window/handles"},

	SetScreenOrientation: {http.MethodPost, "/session/$sessionId/orientation"},
	GetScreenOrientation: {http.MethodGet, "/session/$sessionId/orientation"},

	GetTimeouts: {http.MethodGet, "/session/$sessionId/timeouts"},
	SetTimeouts: {http.MethodPost, "/session/$sessionId/timeouts"},

	W3CDismissAlert:  {http.MethodPost, "/session/$sessionId/alert/dismiss"},
	W3CAcceptAlert:   {http.MethodPost, "/session/$sessionId/alert/accept"},
	W3CSetAlertValue: {http.MethodPost, "/session/$sessionId/alert/text"},
	W3CGetAlertText:  {http.MethodGet, "/session/$sessionId/alert/text"},

	TakeScreenshot:        {http.MethodGet, "/session/$sessionId/screenshot"},
	TakeElementScreenshot: {http.MethodGet, "/session/$sessionId/element/$id/screenshot"},

	GetWindowRect:    {http.MethodGet, "/session/$sessionId/window/rect"},
	SetWindowRect:    {http.MethodPost, "/session/$sessionId/window/rect"},
	MaximizeWindow:   {http.MethodPost, "/session/$sessionId/window/maximize"},
	MinimizeWindow:   {http.MethodPost, "/session/$sessionId/window/minimize"},
	FullscreenWindow: {http.MethodPost, "/session/$sessionId/window/fullscreen"},

	W3CActions:      {http.MethodPost, "/session/$sessionId/actions"},
	W3CClearActions: {http.MethodDelete, "/session/$sessionId/actions"},

	GetShadowRoot:              {http.MethodGet, "/session/$sessionId/root"},
	FindElementFromShadowRoot:  {http.MethodPost, "/session/$sessionId/root/$shadowId/element"},
	FindElementsFromShadowRoot: {http.MethodPost, "/session/$sessionId/root/$shadowId/elements"},

	AddVirtualAuthenticator:    {http.MethodPost, "/session/$sessionId/authenticators"},
	RemoveVirtualAuthenticator: {http.MethodDelete, "/session/$sessionId/authenticators/$authenticatorId"},
	AddCredential:              {http.MethodPost, "/session/$sessionId/authenticators/$authenticatorId/credentials"},
	GetCredentials:             {http.MethodGet, "/session/$sessionId/authenticators/$authenticatorId/credentials"},
	RemoveCredential:           {http.MethodDelete, "/session/$sessionId/authenticators/$authenticatorId/credentials/$credentialId"},
	RemoveAllCredentials:       {http.MethodDelete, "/session/$sessionId/authenticators/$authenticatorId/credentials"},
	SetUserVerified:            {http.MethodPost, "/session/$sessionId/authenticators/$authenticatorId/uv"},

	UploadFile:              {http.MethodPost, "/session/$sessionId/se/file"},
	GetDownloadableFiles:    {http.MethodGet, "/session/$sessionId/se/files"},
	DownloadFile:            {http.MethodPost, "/session/$sessionId/se/files"},
	DeleteDownloadedFile:    {http.MethodDelete, "/session/$sessionId/se/files"},
	DeleteDownloadableFiles: {http.MethodDelete, "/session/$sessionId/se/files"},

	GetFedCMTitle:          {http.MethodGet, "/session/$sessionId/fedcm/gettitle"},
	GetFedCMDialogType:     {http.MethodGet, "/session/$sessionId/fedcm/getdialogtype"},
	GetFedCMAccountList:    {http.MethodGet, "/session/$sessionId/fedcm/accountlist"},
	SelectFedCMAccount:     {http.MethodPost, "/session/$sessionId/fedcm/selectaccount"},
	CancelFedCMDialog:      {http.MethodDelete, "/session/$sessionId/fedcm/canceldialog"},
	SetFedCMDelay:          {http.MethodPost, "/session/$sessionId/fedcm/setdelayenabled"},
	ClickFedCMDialogButton: {http.MethodPost, "/session/$sessionId/fedcm/clickdialogbutton"},
	ResetFedCMCooldown:     {http.MethodPost, "/session/$sessionId/fedcm/resetcooldown"},

	GetNetworkConnection: {http.MethodGet, "/session/$sessionId/network_connection"},
	SetNetworkConnection: {http.MethodPost, "/session/$sessionId/network_connection"},
	CurrentContextHandle: {http.MethodGet, "/session/$sessionId/context"},
	ContextHandles:       {http.MethodGet, "/session/$sessionId/contexts"},
	SwitchToContext:      {http.MethodPost, "/session/$sessionId/context"},
}
