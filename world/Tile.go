package world

import "github.com/nsf/termbox-go"

type Tile struct {
	Solid      bool
	Passable   bool
	Symbol     rune
	Name       string
	Foreground termbox.Attribute
	Background termbox.Attribute
}
