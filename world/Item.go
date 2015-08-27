package world

type Item interface {
	Entity
	Item()
}

type Weapon struct {
	Location P
	Attack   Attack
	Tile     Tile
}

func (w *Weapon) At() P {
	return w.Location
}
func (w *Weapon) MoveTo(p P) {
	w.Location = p
}
func (w *Weapon) Appearance() Tile {
	return w.Tile
}
func (w *Weapon) BasicName() string {
	return w.Tile.Name
}
func (w *Weapon) Item() {
}
