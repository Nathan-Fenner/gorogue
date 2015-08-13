package world

import "github.com/nsf/termbox-go"

type Map struct {
	Width    int
	Height   int
	Tiles    map[P]Tile
	Entities []Entity
	Remove   map[Entity]bool
}

func (m *Map) EntitiesAt(loc P) []Entity {
	result := []Entity{}
	for _, entity := range m.Entities {
		if !m.Remove[entity] && entity.At() == loc {
			result = append(result, entity)
		}
	}
	return result
}

func (m *Map) MoveTo(loc P) Bump {
	if loc.X < 0 || loc.Y < 0 || loc.X >= m.Width || loc.Y >= m.Height {
		return Bump{
			Type: BumpSolid,
			Tile: Tile{
				Symbol:     '#',
				Name:       "edge of the world",
				Foreground: termbox.ColorMagenta,
				Background: termbox.ColorBlack,
			},
		}
	}
	tile, ok := m.Tiles[loc]
	if !ok || tile.Solid {
		return Bump{
			Type: BumpSolid,
			Tile: tile,
		}
	}
	entitiesAt := m.EntitiesAt(loc)
	for _, entity := range entitiesAt {
		if entity.Appearance().Solid {
			return Bump{
				Type:   BumpAction,
				Target: entity,
				Tile:   m.Tiles[loc],
			}
		}
	}
	return Bump{
		Type: BumpEmpty,
	}
}

func (m *Map) Player() *Critter {
	for _, entity := range m.Entities {
		if m.Remove[entity] {
			continue
		}
		if critter, ok := entity.(*Critter); ok && critter.Tile.Name == "$player" {
			return critter
		}
	}
	panic("World has no player")
}

func (m *Map) AddEntity(entity Entity) {
	m.Entities = append(m.Entities, entity)
}

func (m *Map) RemoveEntity(entity Entity) {
	m.Remove[entity] = true
}
func (m *Map) CommitRemovals() {
	newEntities := []Entity{}
	for _, entity := range m.Entities {
		if !m.Remove[entity] {
			newEntities = append(newEntities, entity)
		}
	}
	m.Entities = newEntities
	m.Remove = map[Entity]bool{}
}