package world

import "math/rand"
import "github.com/nsf/termbox-go"

func chooseRune(str string) rune {
	runes := []rune(str)
	return runes[rand.Intn(len(runes))]
}

const grassSymbols = ",;.'\"  "
const treeSymbols = "♣♠¶"

func MakeMap() *Map {
	randMap := &Map{
		Width:  70,
		Height: 40,
		Tiles:  make(map[P]Tile),
		Remove: map[Entity]bool{},
		Entities: []Entity{
			&Critter{
				Health: Bar{
					Value:   200,
					Maximum: 200,
				},
				Tile: Tile{
					Symbol:     '@',
					Solid:      true,
					Foreground: termbox.ColorBlue | termbox.AttrBold,
					Name:       "$player",
				},
				Person: Second,
				Location: P{
					X: 4,
					Y: 3,
				},
				Evasion: 10,
				Attack: Attack{
					Name: "holy sword",
					Verb: "slash",
					Damage: Dice{
						Count: 2,
						Size:  3,
						Base:  0,
					},
					Accuracy: 10,
				},
			},
		},
	}
	for y := 0; y < randMap.Height; y++ {
		for x := 0; x < randMap.Width; x++ {
			solid := false
			name := "grass"
			fore := termbox.ColorGreen
			if rand.Intn(2) == 0 {
				fore |= termbox.AttrBold
			}
			symbol := chooseRune(grassSymbols)
			if rand.Intn(10) == 0 {
				solid = true
				name = "tree"
				symbol = chooseRune(treeSymbols)
			}
			randMap.Tiles[p(x, y)] = Tile{
				Solid:      solid,
				Symbol:     symbol,
				Name:       name,
				Foreground: fore,
				Background: termbox.ColorBlack,
			}
		}
	}
	for i := 0; i < 30; i++ {
		bear := &Critter{
			Health: Bar{
				Value:   35,
				Maximum: 35,
			},
			Tile: Tile{
				Symbol:     'B',
				Solid:      true,
				Foreground: termbox.ColorBlack | termbox.AttrBold,
				Name:       "bear",
			},
			Location: P{
				X: rand.Intn(randMap.Width),
				Y: rand.Intn(randMap.Height),
			},
			Evasion: 3,
			Attack: Attack{
				Name: "bear claws",
				Verb: "claw",
				Damage: Dice{
					Count: 1,
					Size:  20,
					Base:  0,
				},
				Accuracy: 5,
			},
			Brain: &Hunter{},
		}
		randMap.AddEntity(bear)
	}
	return randMap
}
