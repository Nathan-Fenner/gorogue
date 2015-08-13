package world

type BumpType int

const (
	BumpEmpty = iota
	BumpSolid
	BumpAction
)

type Bump struct {
	Type   BumpType
	Tile   Tile
	Target Entity
}
