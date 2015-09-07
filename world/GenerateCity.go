package world

import "math/rand"
import "github.com/nsf/termbox-go"

func splatRectangle(at P, width int, height int, onto map[P]bool) {
	for ax := at.X; ax < at.X+width; ax++ {
		for ay := at.Y; ay < at.Y+height; ay++ {
			onto[p(ax, ay)] = true
		}
	}
}

func CarveSpace(width, height, density int) map[P]bool {
	r := map[P]bool{}
	for i := 0; i < width*height*density/20000; i++ {
		spot := p(rand.Intn(width), rand.Intn(height))
		w, h := rand.Intn(4)+3, rand.Intn(4)+3
		for j := 0; j < rand.Intn(4)+2; j++ {
			splatRectangle(spot.Sub(p(w/2, h/2)).Add(p(rand.Intn(7)-3, rand.Intn(7)-3)), w, h, r)
		}
	}
	return r
}

func PlantWater(width, height, density int) map[P]bool {
	r := map[P]bool{}
	sources := []P{}
	for i := 0; i < width*height*density/20000; i++ {
		sources = append(sources, p(rand.Intn(width), rand.Intn(height)))
	}

	padding := 30

	for x := -padding; x < width+padding; x++ {
		for y := -padding; y < height+padding; y++ {
			s := 0.0
			for _, source := range sources {
				d := float64(p(x, y).Distance(source)) + 1
				s += 1 / d / d
			}
			if s > 0.04 {
				r[p(x, y)] = true
			}
		}
	}
	return r
}

func GenerateCity() *Level {
	level := &Level{
		Tiles:    map[P]Tile{},
		Entities: []Entity{},
		Remove:   map[Entity]bool{},
	}
	wholeSpace := map[P]bool{}
	openSpace := CarveSpace(70, 70, 100)
	for p := range openSpace {
		level.Tiles[p] = Tile{
			Solid:      false,
			Passable:   true,
			Symbol:     chooseRune(",.'`"),
			Name:       "cavern floor",
			Foreground: termbox.ColorWhite,
			Background: termbox.ColorBlack,
		}
		wholeSpace[p] = true
	}
	waterSpace := PlantWater(70, 70, 25)
	for p := range waterSpace {
		if openSpace[p] {
			continue
		}
		level.Tiles[p] = Tile{
			Solid:      false,
			Passable:   false,
			Symbol:     '≈',
			Name:       "water",
			Foreground: termbox.ColorBlue | termbox.AttrBold,
			Background: termbox.ColorBlack,
		}
		wholeSpace[p] = true
	}

	for at := range wholeSpace {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				n := at.Add(p(dx, dy))
				if !wholeSpace[n] {
					level.Tiles[n] = Tile{
						Solid:      true,
						Passable:   false,
						Symbol:     '░',
						Name:       "cavern wall",
						Foreground: termbox.ColorWhite,
						Background: termbox.ColorBlack,
					}
				}
			}
		}
	}

	// Now, we connect everything
	cells := make([]P, 0, len(openSpace))
	for p := range openSpace {
		cells = append(cells, p)
	}
	for i := range cells {
		j := rand.Intn(len(cells)-i) + i
		cells[i], cells[j] = cells[j], cells[i]
	}

	marked := map[P]bool{}
	origins := []P{}
	for _, at := range cells {
		if marked[at] {
			continue
		}
		marked[at] = true
		origins = append(origins, at)
		stack := []P{at}
		for len(stack) > 0 {
			s := stack[0]
			stack = stack[1:]
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					n := s.Add(p(dx, dy))
					if openSpace[n] && !marked[n] {
						marked[n] = true
						stack = append(stack, n)
					}
				}
			}
		}
	}

	// Now, we have a bunch of origins.
	// We'll pathfind between them in order to make the best possible network.
	floorCost := 10
	waterCost := 100
	wallCost := 5000

	for i := range origins {
		j := rand.Intn(len(origins)-i) + i
		origins[i], origins[j] = origins[j], origins[i]
	}

	pathSpace := map[P]bool{}

	costFunc := func(n P) int {
		switch {
		case pathSpace[n]:
			return 1
		case openSpace[n]:
			return floorCost
		case waterSpace[n]:
			return waterCost
		default:
			return wallCost
		}
	}

	for i := range origins {
		if i == 0 {
			continue
		}
		path := FindPath(origins[i-1], origins[i], costFunc, true)
		for _, p := range path {
			pathSpace[p] = true
			if level.Tiles[p].Passable {
				continue
			}
			if waterSpace[p] {
				level.Tiles[p] = Tile{
					Solid:      false,
					Passable:   true,
					Symbol:     '%',
					Name:       "wooden bridge",
					Foreground: termbox.ColorYellow,
					Background: termbox.ColorBlack,
				}
			} else {
				level.Tiles[p] = Tile{
					Solid:      false,
					Passable:   true,
					Symbol:     chooseRune(",.'`"),
					Name:       "passage floor",
					Foreground: termbox.ColorWhite,
					Background: termbox.ColorBlack,
				}
			}
		}
	}

	for at := range pathSpace {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				n := at.Add(p(dx, dy))
				if !wholeSpace[n] && !pathSpace[n] {
					level.Tiles[n] = Tile{
						Solid:      true,
						Passable:   false,
						Symbol:     '▓',
						Name:       "cavern wall",
						Foreground: termbox.ColorWhite,
						Background: termbox.ColorBlack,
					}
				}
			}
		}
	}

	level.AddEntity(&Critter{
		Health: Bar{
			Value:   200,
			Maximum: 200,
		},
		Tile: Tile{
			Symbol:     '@',
			Solid:      false,
			Passable:   false,
			Foreground: termbox.ColorRed | termbox.AttrBold,
			Name:       "$player",
		},
		Person:   Second,
		Location: origins[0],
		Evasion:  10,
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

	for i := 0; i < 10+rand.Intn(10); i++ {
		at := cells[i]
		level.AddEntity(&Critter{
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
	return level
}
