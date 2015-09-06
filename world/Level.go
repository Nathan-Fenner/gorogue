package world

import "math/rand"

type Level struct {
	Tiles    map[P]Tile
	Entities []Entity
	Remove   map[Entity]bool
}

func (m *Level) EntitiesAt(loc P) []Entity {
	result := []Entity{}
	for _, entity := range m.Entities {
		if !m.Remove[entity] && entity.At() == loc {
			result = append(result, entity)
		}
	}
	return result
}

func (m *Level) MoveTo(loc P) Bump {
	tile, ok := m.Tiles[loc]
	if !ok || !tile.Passable {
		return Bump{
			Type: BumpSolid,
			Tile: tile,
		}
	}
	entitiesAt := m.EntitiesAt(loc)
	for _, entity := range entitiesAt {
		if !entity.Appearance().Passable {
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

func (m *Level) Player() *Critter {
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

func (m *Level) AddEntity(entity Entity) {
	m.Entities = append(m.Entities, entity)
}

func (m *Level) RemoveEntity(entity Entity) {
	m.Remove[entity] = true
}
func (m *Level) CommitRemovals() {
	newEntities := []Entity{}
	for _, entity := range m.Entities {
		if !m.Remove[entity] {
			newEntities = append(newEntities, entity)
		}
	}
	m.Entities = newEntities
	m.Remove = map[Entity]bool{}
}

func (m *Level) ItemAt(loc P) Item {
	entities := m.EntitiesAt(loc)
	for _, entity := range entities {
		if item, ok := entity.(Item); ok {
			return item
		}
	}
	return nil
}

func (m *Level) PlaceItem(item Item) {
	at := item.At()
	for i := 0; i < 1000; i++ {
		if m.ItemAt(at) == nil {
			item.MoveTo(at)
			m.AddEntity(item)
			return
		}
		dir := []P{p(1, 0), p(0, 1), p(-1, 0), p(0, -1)}[rand.Intn(4)]
		if m.MoveTo(at.Add(dir)).Tile.Passable {
			at = at.Add(dir)
		}
	}
	// sorry, the item was lost!
}

// Interaction

func (world *Level) Advance() {
	for _, entity := range world.Entities {
		if world.Remove[entity] {
			continue
		}
		if critter, ok := entity.(*Critter); ok && critter.Brain != nil && critter.Health.Value > 0 {
			critter.Brain.Step(world, critter)
		}
	}
}

func (world *Level) MovePlayer(dx int, dy int) {
	player := world.Player()
	np := player.At().Add(p(dx, dy))
	bump := world.MoveTo(np)
	switch bump.Type {
	case BumpEmpty:
		player.MoveTo(np)
		entities := world.EntitiesAt(player.Location)
		for _, entity := range entities {
			if item, ok := entity.(Item); ok {
				AddMessage("You're standing on a {f:yellow|%s}.", item.BasicName())
			}
		}
		world.Advance()
	case BumpSolid:
		AddMessage("You're blocked by the %s.", bump.Tile.Name)
	case BumpAction:
		switch target := bump.Target.(type) {
		case *Critter:
			player.AttackTarget(world, target)
			world.Advance()
		default:
			AddMessage("You bump into {f:blue|%s}.", target.BasicName())
		}
	}
}

func (world *Level) TryPickup() {
	entities := world.EntitiesAt(world.Player().Location)
	foundItem := Item(nil)
	for _, entity := range entities {
		if item, ok := entity.(Item); ok {
			foundItem = item
			break
		}
	}
	if foundItem == nil {
		AddMessage("There's nothing to pick up where you're standing.")
		return
	}

	if len(world.Player().Inventory) >= world.Player().InventoryCapacity {
		AddMessage("You don't have enough room for a {f:yellow|%s}.", foundItem.BasicName())
		return
	}

	AddMessage("You pick up a {f:yellow|%s}.", foundItem.BasicName())
	world.RemoveEntity(foundItem)
	world.Player().Inventory = append(world.Player().Inventory, foundItem)

	world.Advance()
}
