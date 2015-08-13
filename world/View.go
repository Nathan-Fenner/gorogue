package world

import (
	"github.com/nsf/termbox-go"
)

type View interface {
	SetCell(x int, y int, symbol rune, fg, bg termbox.Attribute)
}

type ViewStandard struct {
}

func (view ViewStandard) SetCell(x int, y int, symbol rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x, y, symbol, fg, bg)
}

type ViewOffset struct {
	Child View
	X     int
	Y     int
}

func (view ViewOffset) SetCell(x int, y int, symbol rune, fg, bg termbox.Attribute) {
	view.Child.SetCell(x+view.X, y+view.Y, symbol, fg, bg)
}

type ViewClip struct {
	Child  View
	Width  int
	Height int
}

func (view ViewClip) SetCell(x int, y int, symbol rune, fg, bg termbox.Attribute) {
	if x < 0 || y < 0 || x >= view.Width || y >= view.Height {
		return
	}
	view.Child.SetCell(x, y, symbol, fg, bg)
}
func (view ViewClip) Clear(symbol rune, fg, bg termbox.Attribute) {
	for x := 0; x < view.Width; x++ {
		for y := 0; y < view.Height; y++ {
			view.SetCell(x, y, symbol, fg, bg)
		}
	}
}
