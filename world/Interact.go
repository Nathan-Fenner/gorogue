package world

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

func drawWorld(world *Map) {
	player := world.Player()
	viewSize := 15
	view := ViewOffset{ViewStandard{}, 2, 2}
	view = ViewOffset{ViewClip{view, viewSize*2 + 1, viewSize*2 + 1}, viewSize - player.Location.X, viewSize - player.Location.Y}

	for y := player.Location.Y - viewSize; y <= player.Location.Y+viewSize; y++ {
		for x := player.Location.X - viewSize; x <= player.Location.X+viewSize; x++ {
			if tile, ok := world.Tiles[p(x, y)]; ok {
				view.SetCell(x, y, tile.Symbol, tile.Foreground, tile.Background)
			} else {
				view.SetCell(x, y, '.', termbox.ColorWhite, termbox.ColorBlack)
			}
		}
	}
	for _, entity := range world.Entities {
		if world.Remove[entity] {
			continue
		}
		loc := entity.At()
		tile := world.Tiles[loc]
		appearance := entity.Appearance()
		if !appearance.Solid {
			view.SetCell(loc.X, loc.Y, appearance.Symbol, appearance.Foreground, tile.Background)
		}
	}
	for _, entity := range world.Entities {
		if world.Remove[entity] {
			continue
		}
		loc := entity.At()
		tile := world.Tiles[loc]
		appearance := entity.Appearance()
		if appearance.Solid {
			view.SetCell(loc.X, loc.Y, appearance.Symbol, appearance.Foreground, tile.Background)
		}
	}
}

func debug(s string) {
	chars := []rune(s)
	if len(chars) > 80 {
		chars = chars[0:80]
	}
	for i := 0; i < 80; i++ {
		termbox.SetCell(i, 26, ' ', termbox.ColorWhite, termbox.ColorBlack)
	}
	for i, char := range chars {
		termbox.SetCell(i, 26, char, termbox.ColorWhite, termbox.ColorBlack)
	}
	termbox.Flush()
}

var TurnCount = 0
var Messages = []Announcement{}

type Announcement struct {
	Turn    int
	Message string
}

func AddMessage(message string, options ...interface{}) {
	Messages = append(Messages, Announcement{
		Turn:    TurnCount,
		Message: fmt.Sprintf(message, options...),
	})
}

func DrawText(x int, y int, text string) {
	chars := []rune(text)
	forestack := []termbox.Attribute{termbox.ColorWhite}
	backstack := []termbox.Attribute{termbox.ColorBlack}
	peak := 0
	mode := "print"
	key := []rune{}
	value := []rune{}
	colorMap := map[string]termbox.Attribute{
		"white":   termbox.ColorWhite,
		"black":   termbox.ColorBlack,
		"red":     termbox.ColorRed,
		"yellow":  termbox.ColorYellow,
		"green":   termbox.ColorGreen,
		"cyan":    termbox.ColorCyan,
		"blue":    termbox.ColorBlue,
		"magenta": termbox.ColorMagenta,
		"WHITE":   termbox.ColorWhite | termbox.AttrBold,
		"BLACK":   termbox.ColorBlack | termbox.AttrBold,
		"RED":     termbox.ColorRed | termbox.AttrBold,
		"YELLOW":  termbox.ColorYellow | termbox.AttrBold,
		"GREEN":   termbox.ColorGreen | termbox.AttrBold,
		"CYAN":    termbox.ColorCyan | termbox.AttrBold,
		"BLUE":    termbox.ColorBlue | termbox.AttrBold,
		"MAGENTA": termbox.ColorMagenta | termbox.AttrBold,
	}
	dx := 0
	for _, char := range chars {
		switch mode {
		case "format-key":
			if char == ':' {
				mode = "format-value"
				continue
			}
			key = append(key, char)
		case "format-value":
			if char == ';' || char == '|' {
				keyString := string(key)
				switch keyString {
				case "f":
					forestack[peak] = colorMap[string(value)]
				case "b":
					backstack[peak] = colorMap[string(value)]
				default:
					panic(fmt.Sprintf("unknown format key `%s`", key))
				}
				if char == ';' {
					mode = "format-key"
				} else {
					mode = "print"
				}

				key = nil
				value = nil
				continue
			}
			value = append(value, char)
		case "print":
			switch char {
			case '{':
				forestack = append(forestack, forestack[peak])
				backstack = append(backstack, backstack[peak])
				peak++
				mode = "format-key"
				key = nil
				value = nil
			case '}':
				forestack = forestack[:peak]
				backstack = backstack[:peak]
				peak--
			default:
				termbox.SetCell(x+dx, y, char, forestack[peak], backstack[peak])
				dx++
			}
		}
	}
}

func DisplayMessages() {
	shown := 10
	position := 31
	for x := 0; x < 90; x++ {
		for y := position; y <= position+shown; y++ {
			termbox.SetCell(x, y, ' ', termbox.ColorWhite, termbox.ColorBlack)
		}
	}
	lastFew := len(Messages) - shown
	if lastFew < 0 {
		lastFew = 0
	}
	for i := lastFew; i < len(Messages); i++ {
		if Messages[i].Turn == TurnCount {
			DrawText(1, position+i-lastFew, Messages[i].Message)
		} else {
			DrawText(1, position+i-lastFew, Messages[i].Message)
		}
	}
}

func MovePlayer(world *Map, dx int, dy int) {
	player := world.Player()
	np := player.At().Add(p(dx, dy))
	bump := world.MoveTo(np)
	switch bump.Type {
	case BumpEmpty:
		player.MoveTo(np)
		AdvanceWorld(world)
	case BumpSolid:
		AddMessage("You're blocked by the %s.", bump.Tile.Name)
	case BumpAction:
		switch target := bump.Target.(type) {
		case *Critter:
			player.AttackTarget(world, target)
			AdvanceWorld(world)
		default:
			AddMessage("You bump into {f:WHITE|%s}.", target.BasicName())
		}
	}
}

func AdvanceWorld(world *Map) {
	for _, entity := range world.Entities {
		if world.Remove[entity] {
			continue
		}
		if critter, ok := entity.(*Critter); ok && critter.Brain != nil && critter.Health.Value > 0 {
			critter.Brain.Step(world, critter)
		}
	}
}

func Play() {
	rand.Seed(time.Now().Unix())
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	world := MakeMap()
	defer termbox.Close()
	for {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyArrowLeft:
				MovePlayer(world, -1, 0)
			case termbox.KeyArrowRight:
				MovePlayer(world, 1, 0)
			case termbox.KeyArrowUp:
				MovePlayer(world, 0, -1)
			case termbox.KeyArrowDown:
				MovePlayer(world, 0, 1)
			case termbox.KeyEsc:
				return
			default:
				switch event.Ch {
				case '5':
					// Wait
					AdvanceWorld(world)
				case '0', '1', '2', '3', '4', '6', '7', '8', '9':
					num := int(event.Ch - 48)
					dx := ((num - 1) % 3) - 1
					dy := 1 - ((num - 1) / 3)
					MovePlayer(world, dx, dy)
				}
			}

		}

		drawWorld(world)
		DisplayMessages()
		DrawText(0, 0, "Forest - Depth 0")
		DrawText(34, 0, fmt.Sprintf("Player: %s", world.Player().Health.Render(40)))

		err := termbox.Flush()
		if err != nil {
			panic(err)
		}
	}
}
