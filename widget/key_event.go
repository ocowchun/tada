package widget

import "github.com/gdamore/tcell"

const (
	KeyRune           = "Rune"
	KeyEnter          = "Enter"
	KeyBackspace      = "Backspace"
	KeyTab            = "Tab"
	KeyBacktab        = "Backtab"
	KeyEsc            = "Esc"
	KeyBackspace2     = "Backspace2"
	KeyDelete         = "Delete"
	KeyInsert         = "Insert"
	KeyUp             = "Up"
	KeyDown           = "Down"
	KeyLeft           = "Left"
	KeyRight          = "Right"
	KeyHome           = "Home"
	KeyEnd            = "End"
	KeyUpLeft         = "UpLeft"
	KeyUpRight        = "UpRight"
	KeyDownLeft       = "DownLeft"
	KeyDownRight      = "DownRight"
	KeyCenter         = "Center"
	KeyPgDn           = "PgDn"
	KeyPgUp           = "PgUp"
	KeyClear          = "Clear"
	KeyExit           = "Exit"
	KeyCancel         = "Cancel"
	KeyPause          = "Pause"
	KeyPrint          = "Print"
	KeyF1             = "F1"
	KeyF2             = "F2"
	KeyF3             = "F3"
	KeyF4             = "F4"
	KeyF5             = "F5"
	KeyF6             = "F6"
	KeyF7             = "F7"
	KeyF8             = "F8"
	KeyF9             = "F9"
	KeyF10            = "F10"
	KeyF11            = "F11"
	KeyF12            = "F12"
	KeyF13            = "F13"
	KeyF14            = "F14"
	KeyF15            = "F15"
	KeyF16            = "F16"
	KeyF17            = "F17"
	KeyF18            = "F18"
	KeyF19            = "F19"
	KeyF20            = "F20"
	KeyF21            = "F21"
	KeyF22            = "F22"
	KeyF23            = "F23"
	KeyF24            = "F24"
	KeyF25            = "F25"
	KeyF26            = "F26"
	KeyF27            = "F27"
	KeyF28            = "F28"
	KeyF29            = "F29"
	KeyF30            = "F30"
	KeyF31            = "F31"
	KeyF32            = "F32"
	KeyF33            = "F33"
	KeyF34            = "F34"
	KeyF35            = "F35"
	KeyF36            = "F36"
	KeyF37            = "F37"
	KeyF38            = "F38"
	KeyF39            = "F39"
	KeyF40            = "F40"
	KeyF41            = "F41"
	KeyF42            = "F42"
	KeyF43            = "F43"
	KeyF44            = "F44"
	KeyF45            = "F45"
	KeyF46            = "F46"
	KeyF47            = "F47"
	KeyF48            = "F48"
	KeyF49            = "F49"
	KeyF50            = "F50"
	KeyF51            = "F51"
	KeyF52            = "F52"
	KeyF53            = "F53"
	KeyF54            = "F54"
	KeyF55            = "F55"
	KeyF56            = "F56"
	KeyF57            = "F57"
	KeyF58            = "F58"
	KeyF59            = "F59"
	KeyF60            = "F60"
	KeyF61            = "F61"
	KeyF62            = "F62"
	KeyF63            = "F63"
	KeyF64            = "F64"
	KeyCtrlA          = "Ctrl-A"
	KeyCtrlB          = "Ctrl-B"
	KeyCtrlC          = "Ctrl-C"
	KeyCtrlD          = "Ctrl-D"
	KeyCtrlE          = "Ctrl-E"
	KeyCtrlF          = "Ctrl-F"
	KeyCtrlG          = "Ctrl-G"
	KeyCtrlJ          = "Ctrl-J"
	KeyCtrlK          = "Ctrl-K"
	KeyCtrlL          = "Ctrl-L"
	KeyCtrlN          = "Ctrl-N"
	KeyCtrlO          = "Ctrl-O"
	KeyCtrlP          = "Ctrl-P"
	KeyCtrlQ          = "Ctrl-Q"
	KeyCtrlR          = "Ctrl-R"
	KeyCtrlS          = "Ctrl-S"
	KeyCtrlT          = "Ctrl-T"
	KeyCtrlU          = "Ctrl-U"
	KeyCtrlV          = "Ctrl-V"
	KeyCtrlW          = "Ctrl-W"
	KeyCtrlX          = "Ctrl-X"
	KeyCtrlY          = "Ctrl-Y"
	KeyCtrlZ          = "Ctrl-Z"
	KeyCtrlSpace      = "Ctrl-Space"
	KeyCtrlUnderscore = "Ctrl-_"
	KeyCtrlRightSq    = "Ctrl-]"
	KeyCtrlBackslash  = "Ctrl-\\"
	KeyCtrlCarat      = "Ctrl-^"
)

var keyNames = map[tcell.Key]string{
	tcell.KeyRune:           KeyRune,
	tcell.KeyEnter:          KeyEnter,
	tcell.KeyBackspace:      KeyBackspace,
	tcell.KeyTab:            KeyTab,
	tcell.KeyBacktab:        KeyBacktab,
	tcell.KeyEsc:            KeyEsc,
	tcell.KeyBackspace2:     KeyBackspace2,
	tcell.KeyDelete:         KeyDelete,
	tcell.KeyInsert:         KeyInsert,
	tcell.KeyUp:             KeyUp,
	tcell.KeyDown:           KeyDown,
	tcell.KeyLeft:           KeyLeft,
	tcell.KeyRight:          KeyRight,
	tcell.KeyHome:           KeyHome,
	tcell.KeyEnd:            KeyEnd,
	tcell.KeyUpLeft:         KeyUpLeft,
	tcell.KeyUpRight:        KeyUpRight,
	tcell.KeyDownLeft:       KeyDownLeft,
	tcell.KeyDownRight:      KeyDownRight,
	tcell.KeyCenter:         KeyCenter,
	tcell.KeyPgDn:           KeyPgDn,
	tcell.KeyPgUp:           KeyPgUp,
	tcell.KeyClear:          KeyClear,
	tcell.KeyExit:           KeyExit,
	tcell.KeyCancel:         KeyCancel,
	tcell.KeyPause:          KeyPause,
	tcell.KeyPrint:          KeyPrint,
	tcell.KeyF1:             KeyF1,
	tcell.KeyF2:             KeyF2,
	tcell.KeyF3:             KeyF3,
	tcell.KeyF4:             KeyF4,
	tcell.KeyF5:             KeyF5,
	tcell.KeyF6:             KeyF6,
	tcell.KeyF7:             KeyF7,
	tcell.KeyF8:             KeyF8,
	tcell.KeyF9:             KeyF9,
	tcell.KeyF10:            KeyF10,
	tcell.KeyF11:            KeyF11,
	tcell.KeyF12:            KeyF12,
	tcell.KeyF13:            KeyF13,
	tcell.KeyF14:            KeyF14,
	tcell.KeyF15:            KeyF15,
	tcell.KeyF16:            KeyF16,
	tcell.KeyF17:            KeyF17,
	tcell.KeyF18:            KeyF18,
	tcell.KeyF19:            KeyF19,
	tcell.KeyF20:            KeyF20,
	tcell.KeyF21:            KeyF21,
	tcell.KeyF22:            KeyF22,
	tcell.KeyF23:            KeyF23,
	tcell.KeyF24:            KeyF24,
	tcell.KeyF25:            KeyF25,
	tcell.KeyF26:            KeyF26,
	tcell.KeyF27:            KeyF27,
	tcell.KeyF28:            KeyF28,
	tcell.KeyF29:            KeyF29,
	tcell.KeyF30:            KeyF30,
	tcell.KeyF31:            KeyF31,
	tcell.KeyF32:            KeyF32,
	tcell.KeyF33:            KeyF33,
	tcell.KeyF34:            KeyF34,
	tcell.KeyF35:            KeyF35,
	tcell.KeyF36:            KeyF36,
	tcell.KeyF37:            KeyF37,
	tcell.KeyF38:            KeyF38,
	tcell.KeyF39:            KeyF39,
	tcell.KeyF40:            KeyF40,
	tcell.KeyF41:            KeyF41,
	tcell.KeyF42:            KeyF42,
	tcell.KeyF43:            KeyF43,
	tcell.KeyF44:            KeyF44,
	tcell.KeyF45:            KeyF45,
	tcell.KeyF46:            KeyF46,
	tcell.KeyF47:            KeyF47,
	tcell.KeyF48:            KeyF48,
	tcell.KeyF49:            KeyF49,
	tcell.KeyF50:            KeyF50,
	tcell.KeyF51:            KeyF51,
	tcell.KeyF52:            KeyF52,
	tcell.KeyF53:            KeyF53,
	tcell.KeyF54:            KeyF54,
	tcell.KeyF55:            KeyF55,
	tcell.KeyF56:            KeyF56,
	tcell.KeyF57:            KeyF57,
	tcell.KeyF58:            KeyF58,
	tcell.KeyF59:            KeyF59,
	tcell.KeyF60:            KeyF60,
	tcell.KeyF61:            KeyF61,
	tcell.KeyF62:            KeyF62,
	tcell.KeyF63:            KeyF63,
	tcell.KeyF64:            KeyF64,
	tcell.KeyCtrlA:          KeyCtrlA,
	tcell.KeyCtrlB:          KeyCtrlB,
	tcell.KeyCtrlC:          KeyCtrlC,
	tcell.KeyCtrlD:          KeyCtrlD,
	tcell.KeyCtrlE:          KeyCtrlE,
	tcell.KeyCtrlF:          KeyCtrlF,
	tcell.KeyCtrlG:          KeyCtrlG,
	tcell.KeyCtrlJ:          KeyCtrlJ,
	tcell.KeyCtrlK:          KeyCtrlK,
	tcell.KeyCtrlL:          KeyCtrlL,
	tcell.KeyCtrlN:          KeyCtrlN,
	tcell.KeyCtrlO:          KeyCtrlO,
	tcell.KeyCtrlP:          KeyCtrlP,
	tcell.KeyCtrlQ:          KeyCtrlQ,
	tcell.KeyCtrlR:          KeyCtrlR,
	tcell.KeyCtrlS:          KeyCtrlS,
	tcell.KeyCtrlT:          KeyCtrlT,
	tcell.KeyCtrlU:          KeyCtrlU,
	tcell.KeyCtrlV:          KeyCtrlV,
	tcell.KeyCtrlW:          KeyCtrlW,
	tcell.KeyCtrlX:          KeyCtrlX,
	tcell.KeyCtrlY:          KeyCtrlY,
	tcell.KeyCtrlZ:          KeyCtrlZ,
	tcell.KeyCtrlSpace:      KeyCtrlSpace,
	tcell.KeyCtrlUnderscore: KeyCtrlUnderscore,
	tcell.KeyCtrlRightSq:    KeyCtrlRightSq,
	tcell.KeyCtrlBackslash:  KeyCtrlBackslash,
	tcell.KeyCtrlCarat:      KeyCtrlCarat,
}

// type Key tcell.Key
type KeyEvent struct {
	Rune rune
	Key  string
}

func ConvertEvent(event *tcell.EventKey) *KeyEvent {
	return &KeyEvent{
		Rune: event.Rune(),
		Key:  keyNames[event.Key()],
	}
}
