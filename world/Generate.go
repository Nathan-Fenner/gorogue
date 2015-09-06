package world

import "math/rand"
import "github.com/nsf/termbox-go"

func chooseRune(str string) rune {
	runes := []rune(str)
	return runes[rand.Intn(len(runes))]
}

const grassSymbols = ",;.'\"  "
const treeSymbols = "♣♠¶"

func SurfaceLevel() *Level {
	width := 200
	height := 200
	randLevel := &Level{
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
					Passable:   false,
					Solid:      false,
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
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
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
			randLevel.Tiles[p(x, y)] = Tile{
				Solid:      solid,
				Passable:   !solid,
				Symbol:     symbol,
				Name:       name,
				Foreground: fore,
				Background: termbox.ColorBlack,
			}
		}
	}
	for i := 0; i < 30; i++ {
		goblin := &Critter{
			Health: Bar{
				Value:   15,
				Maximum: 15,
			},
			Tile: Tile{
				Symbol:     'g',
				Passable:   false,
				Solid:      false,
				Foreground: termbox.ColorBlack | termbox.AttrBold,
				Name:       "goblin",
			},
			Location: P{
				X: rand.Intn(width),
				Y: rand.Intn(height),
			},
			Evasion: 3,
			Attack: Attack{
				Name: "crude spear",
				Verb: "stab",
				Damage: Dice{
					Count: 1,
					Size:  5,
					Base:  0,
				},
				Accuracy: 5,
			},
			Brain: &Hunter{},
		}
		randLevel.AddEntity(goblin)
	}
	return randLevel
}
