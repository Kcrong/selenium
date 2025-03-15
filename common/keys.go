package common

// KeyboardKey represents a keyboard key
type KeyboardKey string

const (
	// NULL null key
	NULL KeyboardKey = "\ue000"
	// CANCEL cancel key
	CANCEL KeyboardKey = "\ue001"
	// HELP help key
	HELP KeyboardKey = "\ue002"
	// BACKSPACE backspace key
	BACKSPACE KeyboardKey = "\ue003"
	// TAB tab key
	TAB KeyboardKey = "\ue004"
	// CLEAR clear key
	CLEAR KeyboardKey = "\ue005"
	// RETURN return key
	RETURN KeyboardKey = "\ue006"
	// ENTER enter key
	ENTER KeyboardKey = "\ue007"
	// SHIFT shift key
	SHIFT KeyboardKey = "\ue008"
	// CONTROL control key
	CONTROL KeyboardKey = "\ue009"
	// ALT alt key
	ALT KeyboardKey = "\ue00a"
	// PAUSE pause key
	PAUSE KeyboardKey = "\ue00b"
	// ESCAPE escape key
	ESCAPE KeyboardKey = "\ue00c"
	// SPACE space key
	SPACE KeyboardKey = "\ue00d"
	// PAGE_UP page up key
	PAGE_UP KeyboardKey = "\ue00e"
	// PAGE_DOWN page down key
	PAGE_DOWN KeyboardKey = "\ue00f"
	// END end key
	END KeyboardKey = "\ue010"
	// HOME home key
	HOME KeyboardKey = "\ue011"
	// LEFT left arrow key
	LEFT KeyboardKey = "\ue012"
	// UP up arrow key
	UP KeyboardKey = "\ue013"
	// RIGHT right arrow key
	RIGHT KeyboardKey = "\ue014"
	// DOWN down arrow key
	DOWN KeyboardKey = "\ue015"
	// INSERT insert key
	INSERT KeyboardKey = "\ue016"
	// DELETE delete key
	DELETE KeyboardKey = "\ue017"
)

const (
	// Semicolon is SEMICOLON key
	Semicolon KeyboardKey = "\ue018"
	// Equals equals key
	Equals KeyboardKey = "\ue019"
)

// Numpad keys
const (
	// Numpad0 numpad 0 key
	Numpad0 KeyboardKey = "\ue01a"
	// Numpad1 numpad 1 key
	Numpad1 KeyboardKey = "\ue01b"
	// Numpad2 numpad 2 key
	Numpad2 KeyboardKey = "\ue01c"
	// Numpad3 numpad 3 key
	Numpad3 KeyboardKey = "\ue01d"
	// Numpad4 numpad 4 key
	Numpad4 KeyboardKey = "\ue01e"
	// Numpad5 numpad 5 key
	Numpad5 KeyboardKey = "\ue01f"
	// Numpad6 numpad 6 key
	Numpad6 KeyboardKey = "\ue020"
	// Numpad7 numpad 7 key
	Numpad7 KeyboardKey = "\ue021"
	// Numpad8 numpad 8 key
	Numpad8 KeyboardKey = "\ue022"
	// Numpad9 numpad 9 key
	Numpad9 KeyboardKey = "\ue023"
)

// Function keys
const (
	// F1 F1 key
	F1 KeyboardKey = "\ue031"
	// F2 F2 key
	F2 KeyboardKey = "\ue032"
	// F3 F3 key
	F3 KeyboardKey = "\ue033"
	// F4 F4 key
	F4 KeyboardKey = "\ue034"
	// F5 F5 key
	F5 KeyboardKey = "\ue035"
	// F6 F6 key
	F6 KeyboardKey = "\ue036"
	// F7 F7 key
	F7 KeyboardKey = "\ue037"
	// F8 F8 key
	F8 KeyboardKey = "\ue038"
	// F9 F9 key
	F9 KeyboardKey = "\ue039"
	// F10 F10 key
	F10 KeyboardKey = "\ue03a"
	// F11 F11 key
	F11 KeyboardKey = "\ue03b"
	// F12 F12 key
	F12 KeyboardKey = "\ue03c"
)

// Meta keys
const (
	// Meta key
	Meta KeyboardKey = "\ue03d"
	// Command key
	Command KeyboardKey = "\ue03d"
	// ZenkakuHankaku zenkaku hankaku key
	ZenkakuHankaku KeyboardKey = "\ue040"
)
