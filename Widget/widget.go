package widget

import ui "github.com/gizak/termui/v3"

type Widget interface {
	HandleUIEvent(event ui.Event)
	ui.Drawable
	SetBorderStyle(style ui.Style)
	Refresh()
	//SetRect(x1, y1, x2, y2 int)
}
