package world

import "fmt"

type Corpse struct {
	Location P
	Tile     Tile
}

func (self *Corpse) At() P {
	return self.Location
}

func (self *Corpse) MoveTo(loc P) {
	self.Location = loc
}

func (self *Corpse) Appearance() Tile {
	return self.Tile
}

func (self *Corpse) BasicName() string {
	if self.Appearance().Name == "$player" {
		return "your corpse"
	}
	return fmt.Sprintf("%s corpse", MakeIndefinite(self.Appearance().Name))
}
