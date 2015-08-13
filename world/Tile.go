package world

import "github.com/nsf/termbox-go"

type Tile struct {
	Solid      bool
	Symbol     rune
	Name       string
	Foreground termbox.Attribute
	Background termbox.Attribute
}
