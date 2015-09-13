package world

import "math"
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

type ForestOptions struct {
	Enabled  bool
	Size     int
	Sparsity float64
}

type ChasmOptions struct {
	Enabled bool
}

type CavernOptions struct {
	ExtraTunnels  bool
	Watery        bool
	ForestOptions ForestOptions
	ChasmOptions  ChasmOptions
}

func GenerateCavern(options CavernOptions) *Level {
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

	chasmSpace := map[P]bool{}

	// ExtraTunnels conflicts badly with the chasms
	if options.ChasmOptions.Enabled && !options.ExtraTunnels {
		// Generate chasms
		// Pick random points and streak along
		for i := 0; i < 5; i++ {
			s := p(rand.Intn(70), rand.Intn(70))
			dx, dy := rand.Float64()-0.5, rand.Float64()-0.5
			for dx*dx+dy*dy > 1 {
				dx, dy = rand.Float64()-0.5, rand.Float64()-0.5
			}
			dm := math.Sqrt(dx*dx + dy*dy)
			dx, dy = dx/dm, dy/dm
			dtx, dty := rand.Float64()-0.5, rand.Float64()-0.5
			for dtx*dtx+dty*dty > 1 {
				dtx, dty = rand.Float64()-0.5, rand.Float64()-0.5
			}
			if dx*dtx+dy*dty < 0 {
				dtx, dty = -dtx, -dty
			}
			dtm := math.Sqrt(dtx*dtx + dty*dty)
			dtx, dty = dtx/dtm, dty/dtm
			length := rand.Float64()*35 + 15
			for t := 0.0; t < length; t += 1.0 {
				mx := dtx*t + (1-t)*dx
				my := dty*t + (1-t)*dy
				mm := math.Sqrt(mx*mx + my*my)
				mx, my = mx/mm, my/mm
				at := s.Add(p(int(mx*t), int(my*t)))
				rw := rand.Intn(5) + 2
				rh := rand.Intn(5) + 2
				for ix := 0; ix < rw; ix++ {
					for iy := 0; iy < rh; iy++ {
						tile := at.Add(p(ix-rw/2, iy-rh/2))
						chasmSpace[tile] = true
						level.Tiles[tile] = Tile{
							Solid:      false,
							Passable:   false,
							Symbol:     '.',
							Name:       "chasm",
							Foreground: termbox.ColorBlack | termbox.AttrBold,
							Background: termbox.ColorBlack,
						}
						delete(openSpace, tile)
						wholeSpace[tile] = true
					}
				}
			}
		}
	}

	waterSpace := map[P]bool{}

	if !options.ChasmOptions.Enabled {

		waterAmount := 25
		if options.Watery {
			waterAmount = 50
		}
		waterSpace = PlantWater(70, 70, waterAmount)
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
	pathCost := 1
	floorCost := 10
	waterCost := 100
	wallCost := 5000
	chasmCost := 200

	for i := range origins {
		j := rand.Intn(len(origins)-i) + i
		origins[i], origins[j] = origins[j], origins[i]
	}

	pathSpace := map[P]bool{}

	costFunc := func(n P) int {
		switch {
		case pathSpace[n]:
			return pathCost
		case chasmSpace[n]:
			return chasmCost
		case openSpace[n]:
			return floorCost
		case waterSpace[n]:
			return waterCost
		default:
			return wallCost
		}
	}

	connectPath := func(i, j int) {
		path := FindPath(origins[i], origins[j], costFunc, true)
		for _, p := range path {
			pathSpace[p] = true
			if level.Tiles[p].Passable {
				continue
			}
			if waterSpace[p] || chasmSpace[p] {
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

	for i := range origins {
		if i == 0 {
			continue
		}
		connectPath(i-1, i)
	}

	if options.ExtraTunnels {
		// New, wider paths
		pathCost = wallCost
		waterCost = wallCost * 300
		chasmCost = wallCost * 400
		for i := range origins {
			for j := range origins {
				if j >= i {
					break
				}
				if rand.Float32() < 0.2 {
					connectPath(j, i)
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

	if options.ForestOptions.Enabled {
		forestSpace := PlantWater(70, 70, 40+options.ForestOptions.Size)
		forestCells := []P{}
		for p := range forestSpace {
			if level.Tiles[p].Name == "cavern floor" || level.Tiles[p].Name == "passage floor" {
				level.Tiles[p] = Tile{
					Solid:      false,
					Passable:   true,
					Symbol:     chooseRune(",.'`"),
					Name:       "grassy soil",
					Foreground: termbox.ColorGreen,
					Background: termbox.ColorBlack,
				}
			}
			if level.Tiles[p].Name == "cavern wall" {
				level.Tiles[p] = Tile{
					Solid:      true,
					Passable:   false,
					Symbol:     '░',
					Name:       "mossy wall",
					Foreground: termbox.ColorGreen,
					Background: termbox.ColorBlack,
				}
			}
			forestCells = append(forestCells, p)
		}
		for i := range forestCells {
			j := rand.Intn(len(forestCells)-i) + i
			forestCells[i], forestCells[j] = forestCells[j], forestCells[i]
		}
		for _, at := range forestCells {
			if level.Tiles[at].Name != "grassy soil" {
				continue
			}
			obstructed := 0
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if !level.Tiles[at.Add(p(dx, dy))].Passable {
						obstructed++
					}
				}
			}
			if obstructed <= 1 && rand.Float64() > options.ForestOptions.Sparsity {
				level.Tiles[at] = Tile{
					Solid:      true,
					Passable:   false,
					Symbol:     chooseRune("♣♠"),
					Name:       "tree",
					Foreground: termbox.ColorGreen,
					Background: termbox.ColorBlack,
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
				Foreground: termbox.ColorGreen | termbox.AttrBold,
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
