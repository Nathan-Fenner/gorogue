package world

type Entity interface {
	At() P
	MoveTo(P)
	Appearance() Tile
	BasicName() string
}
