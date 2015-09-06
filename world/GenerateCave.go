package world

import "math/rand"

import "github.com/nsf/termbox-go"

func goblinDrops() []Item {
	items := []Item{
		&Weapon{
			Tile: Tile{
				Solid:      false,
				Passable:   true,
				Symbol:     '*',
				Name:       "gold pile",
				Foreground: termbox.ColorYellow | termbox.AttrBold,
				Background: termbox.ColorBlack,
			},
		},
		&Weapon{
			Tile: Tile{
				Solid:      false,
				Passable:   true,
				Symbol:     '↑',
				Name:       "wooden spear",
				Foreground: termbox.ColorYellow | termbox.AttrBold,
				Background: termbox.ColorBlack,
			},
		},
		&Weapon{
			Tile: Tile{
				Solid:      false,
				Passable:   true,
				Symbol:     ']',
				Name:       "leather armor",
				Foreground: termbox.ColorYellow | termbox.AttrBold,
				Background: termbox.ColorBlack,
			},
		},
		&Weapon{
			Tile: Tile{
				Solid:      false,
				Passable:   true,
				Symbol:     '§',
				Name:       "goblin spellbook",
				Foreground: termbox.ColorYellow | termbox.AttrBold,
				Background: termbox.ColorBlack,
			},
		},
	}
	selected := []Item{}
	for _, item := range items {
		if rand.Intn(3) == 1 {
			selected = append(selected, item)
		}
	}
	if len(selected) == 0 {
		selected = append(selected, items[rand.Intn(len(items))])
	}
	return selected
}

func CaveLevel() *Level {
	rooms := 15
	result := &Level{
		Tiles:    make(map[P]Tile),
		Remove:   map[Entity]bool{},
		Entities: []Entity{},
	}
	locations := []P{p(0, 0)}
	StampRoom(result, -1, -1, 3, 3)
	rooms--
	for rooms > 0 {
		base := locations[rand.Intn(len(locations))]
		offset := p(rand.Intn(37)-18, rand.Intn(37)-18)
		size := p(rand.Intn(4)+3, rand.Intn(4)+3)
		corner := base.Add(offset).Sub(p(size.X/2, size.Y/2))
		ok := StampRoom(result, corner.X, corner.Y, size.X, size.Y)
		if ok {
			rooms--
			locations = append(locations, base.Add(offset))
		}
	}

	result.AddEntity(&Critter{
		Health: Bar{
			Value:   200,
			Maximum: 200,
		},
		Tile: Tile{
			Symbol:     '@',
			Solid:      false,
			Passable:   false,
			Foreground: termbox.ColorBlue | termbox.AttrBold,
			Name:       "$player",
		},
		Person: Second,
		Location: P{
			X: 0,
			Y: 0,
		},
		Evasion: 10,
		Attack: Attack{
			Name: "steel dagger",
			Verb: "slash",
			Damage: Dice{
				Count: 2,
				Size:  3,
				Base:  0,
			},
			Accuracy: 10,
		},
	})

	for i, at := range locations {
		if i == 0 {
			continue
		}
		target := locations[rand.Intn(i)]
		StampHall(result, at, target)
	}

	openTiles := []P{}
	for at, tile := range result.Tiles {
		if tile.Passable {
			openTiles = append(openTiles, at)
		}
	}

	for i := range openTiles {
		j := rand.Intn(len(openTiles)-i) + i
		openTiles[i], openTiles[j] = openTiles[j], openTiles[i]
	}

	for i := 0; i < 10+rand.Intn(10); i++ {
		at := openTiles[i]
		result.AddEntity(&Critter{
			Health: Bar{
				Value:   20,
				Maximum: 20,
			},
			Tile: Tile{
				Symbol:     'g',
				Solid:      false,
				Passable:   false,
				Foreground: termbox.ColorGreen,
				Name:       "goblin",
			},
			Person:   Third,
			Location: at,
			Evasion:  10,
			Attack: Attack{
				Name: "stone spear",
				Verb: "stab",
				Damage: Dice{
					Count: 1,
					Size:  4,
					Base:  0,
				},
				Accuracy: 5,
			},
			Brain:             &Hunter{},
			Inventory:         goblinDrops(),
			InventoryCapacity: 3,
		})
	}

	return result
}

func StampHallHorizontal(m *Level, from P, to P) P {
	at := from
	for at.X != to.X {
		StampWalledHall(m, at)
		if at.X < to.X {
			at.X++
		} else {
			at.X--
		}
	}
	return at
}
func StampHallVertical(m *Level, from P, to P) P {
	at := from
	for at.Y != to.Y {
		StampWalledHall(m, at)
		if at.Y < to.Y {
			at.Y++
		} else {
			at.Y--
		}
	}
	return at
}

func StampHall(m *Level, from P, to P) {
	if rand.Intn(2) == 0 {
		corner := StampHallHorizontal(m, from, to)
		StampHallVertical(m, corner, to)
	} else {
		corner := StampHallVertical(m, from, to)
		StampHallHorizontal(m, corner, to)
	}
}

func StampRoom(m *Level, x int, y int, w int, h int) bool {
	for nx := x - 1; nx <= x+w; nx++ {
		for ny := y - 1; ny <= y+h; ny++ {
			if _, ok := m.Tiles[p(nx, ny)]; !ok {
				continue
			}
			expectedName := ""
			if nx >= x && nx < x+w && ny >= y && ny < y+h {
				expectedName = "stone tile"
			} else {
				expectedName = "stone wall"
			}
			if m.Tiles[p(nx, ny)].Name != expectedName {
				return false
			}
		}
	}
	for nx := x - 1; nx <= x+w; nx++ {
		for ny := y - 1; ny <= y+h; ny++ {
			if nx >= x && nx < x+w && ny >= y && ny < y+h {
				StampFloor(m, p(nx, ny))
			} else {
				StampWall(m, p(nx, ny))
			}
		}
	}
	return true
}

func StampFloor(m *Level, at P) {
	m.Tiles[at] = Tile{
		Solid:      false,
		Passable:   true,
		Symbol:     chooseRune(",.'`"),
		Name:       "stone tile",
		Foreground: termbox.ColorWhite,
		Background: termbox.ColorBlack,
	}
}
func StampWall(m *Level, at P) {
	if m.Tiles[at].Name == "stone tile" {
		return
	}
	m.Tiles[at] = Tile{
		Solid:      true,
		Passable:   false,
		Symbol:     '▒',
		Name:       "stone wall",
		Foreground: termbox.ColorWhite,
		Background: termbox.ColorBlack,
	}
}

func StampWalledHall(m *Level, at P) {
	StampFloor(m, at)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx != 0 || dy != 0 {
				StampWall(m, at.Add(p(dx, dy)))
			}
		}
	}
}
