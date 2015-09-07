package world

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"unicode"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

var memory = map[P]bool{}

var inventory = []Item{}
var inventoryCapacity = 20

func drawWorld(world *Level, visibility LightGrid) {
	player := world.Player()
	viewSize := 15
	view := ViewOffset{ViewStandard{}, 2, 2}
	view = ViewOffset{ViewClip{view, viewSize*2 + 1, viewSize*2 + 1}, viewSize - player.Location.X, viewSize - player.Location.Y}

	for y := player.Location.Y - viewSize; y <= player.Location.Y+viewSize; y++ {
		for x := player.Location.X - viewSize; x <= player.Location.X+viewSize; x++ {
			if tile, ok := world.Tiles[p(x, y)]; ok && visibility.Visible(p(x, y)) {
				memory[p(x, y)] = true
				view.SetCell(x, y, tile.Symbol, tile.Foreground, tile.Background)
			} else {
				if memory[p(x, y)] {
					view.SetCell(x, y, tile.Symbol, termbox.ColorBlue, termbox.ColorBlack)
				} else {
					view.SetCell(x, y, ' ', termbox.ColorBlue, termbox.ColorBlack)
				}

			}
		}
	}
	// Draw ground entities:
	for _, entity := range world.Entities {
		if world.Remove[entity] {
			continue
		}
		if !visibility.Visible(entity.At()) {
			continue
		}
		loc := entity.At()
		tile := world.Tiles[loc]
		appearance := entity.Appearance()
		if appearance.Passable {
			view.SetCell(loc.X, loc.Y, appearance.Symbol, appearance.Foreground, tile.Background)
		}
	}
	// Draw items:
	for _, entity := range world.Entities {
		if world.Remove[entity] {
			continue
		}
		if !visibility.Visible(entity.At()) {
			continue
		}
		loc := entity.At()
		tile := world.Tiles[loc]
		appearance := entity.Appearance()
		if _, ok := entity.(Item); ok {
			view.SetCell(loc.X, loc.Y, appearance.Symbol, appearance.Foreground, tile.Background)
		}
	}
	// Draw solid entities:
	for _, entity := range world.Entities {
		if world.Remove[entity] {
			continue
		}
		if !visibility.Visible(entity.At()) {
			continue
		}
		loc := entity.At()
		tile := world.Tiles[loc]
		appearance := entity.Appearance()
		if !appearance.Passable {
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

func uppercaseFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func AddMessage(message string, options ...interface{}) {
	message = uppercaseFirst(fmt.Sprintf(message, options...))
	if len(message) == 0 {
		panic("AddMessage() given empty message")
	}
	Messages = append(Messages, Announcement{
		Turn:    TurnCount,
		Message: message,
	})
}

func DrawText(x int, y int, text string) {
	chars := []rune(text)
	forestack := []string{"grey"}
	backstack := []string{"black"}
	modestack := []string{"normal"}
	peak := 0
	mode := "print"
	key := []rune{}
	value := []rune{}
	colorLevel := map[string]termbox.Attribute{
		"white":   termbox.ColorWhite,
		"black":   termbox.ColorBlack,
		"red":     termbox.ColorRed,
		"yellow":  termbox.ColorYellow,
		"green":   termbox.ColorGreen,
		"cyan":    termbox.ColorCyan,
		"blue":    termbox.ColorBlue,
		"magenta": termbox.ColorMagenta,
		"grey":    termbox.ColorWhite,
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
					forestack[peak] = string(value)
				case "b":
					backstack[peak] = string(value)
				case "m":
					modestack[peak] = string(value)
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
				modestack = append(modestack, modestack[peak])
				peak++
				mode = "format-key"
				key = nil
				value = nil
			case '}':
				forestack = forestack[:peak]
				backstack = backstack[:peak]
				modestack = modestack[:peak]
				peak--
			default:
				foreColor := colorLevel[forestack[peak]]
				backColor := colorLevel[backstack[peak]]
				if modestack[peak] == "bold" {
					foreColor = foreColor | termbox.AttrBold
				}
				if forestack[peak] == "grey" && modestack[peak] == "normal" {
					foreColor = termbox.ColorBlack | termbox.AttrBold
				}
				termbox.SetCell(x+dx, y, char, foreColor, backColor)
				dx++
			}
		}
	}
}

func DisplayMessages() {
	shown := 10
	position := 35
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
			DrawText(1, position+i-lastFew, fmt.Sprintf("{m:bold|%s}", Messages[i].Message))
		} else {
			DrawText(1, position+i-lastFew, fmt.Sprintf("%s", Messages[i].Message))
		}
	}
}

type Sidebar struct {
	Visible Region
	Status  [][]*Critter
}

func (sidebar *Sidebar) AddCritter(c *Critter) {
	if !sidebar.Visible.Inside(c.Location) {
		return
	}
	y := c.Location.Y - sidebar.Visible.Low.Y
	sidebar.Status[y] = append(sidebar.Status[y], c)
}

type CritterOrder []*Critter

func (a CritterOrder) Len() int {
	return len(a)
}
func (a CritterOrder) Less(i, j int) bool {
	return a[i].Location.X < a[j].Location.X
}
func (a CritterOrder) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (sidebar *Sidebar) Display() {
	for x := 33; x < 33+200; x++ {
		for y := 2; y <= 2+31; y++ {
			termbox.SetCell(x, y, ' ', termbox.ColorWhite, termbox.ColorBlack)
		}
	}
	for y, critters := range sidebar.Status {
		if len(critters) == 0 {
			continue
		}
		sort.Sort(CritterOrder(critters))
		for x, critter := range critters {
			DrawText(33+x*10, 2+y, critter.Health.Render(9))
		}
	}
}

func Play() {
	rand.Seed(time.Now().Unix())

	world := GenerateCity()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	termbox.Flush()

	defer termbox.Close()
	for {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyArrowLeft:
				world.MovePlayer(-1, 0)
			case termbox.KeyArrowRight:
				world.MovePlayer(1, 0)
			case termbox.KeyArrowUp:
				world.MovePlayer(0, -1)
			case termbox.KeyArrowDown:
				world.MovePlayer(0, 1)
			case termbox.KeyEsc:
				return
			default:
				switch event.Ch {
				case 'g':
					world.TryPickup()
				case '5':
					// Wait
					world.Advance()
				case '0', '1', '2', '3', '4', '6', '7', '8', '9':
					num := int(event.Ch - 48)
					dx := ((num - 1) % 3) - 1
					dy := 1 - ((num - 1) / 3)
					world.MovePlayer(dx, dy)
				}
			}

		}

		world.CommitRemovals()

		visibility := GetVisibility(world, world.Player().Location)

		drawWorld(world, visibility)
		DisplayMessages()
		TurnCount++
		DrawText(0, 0, "Forest - Depth 0")
		DrawText(34, 0, fmt.Sprintf("Player: %s", world.Player().Health.Render(40)))

		sidebar := &Sidebar{
			Visible: Region{
				Low:  world.Player().Location.Sub(p(15, 15)),
				High: world.Player().Location.Add(p(15, 15)),
			},
			Status: make([][]*Critter, 31),
		}

		for _, entity := range world.Entities {
			if critter, ok := entity.(*Critter); ok && visibility.Visible(critter.Location) {
				sidebar.AddCritter(critter)
			}
		}

		sidebar.Display()

		err := termbox.Flush()
		if err != nil {
			panic(err)
		}
	}
}
