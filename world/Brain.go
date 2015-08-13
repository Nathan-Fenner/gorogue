package world

import "math/rand"

type Brain interface {
	Step(*Map, *Critter)
}

type Wander struct {
}

func (wander *Wander) Step(world *Map, self *Critter) {
	directions := Cardinals()
	choices := []P{}
	for _, direction := range directions {
		adjacent := self.Location.Add(direction)
		if world.MoveTo(adjacent).Type != BumpEmpty {
			continue
		}
		choices = append(choices, direction)
	}
	if rand.Intn(3) != 0 || len(choices) == 0 {
		return // Don't bother moving
	}
	choice := choices[rand.Intn(len(choices))]
	self.MoveTo(self.Location.Add(choice))
}

type Hunter struct {
}

func (hunter *Hunter) Step(world *Map, self *Critter) {
	player := world.Player()
	if self.Location.Distance(player.Location) == 1 {
		self.AttackTarget(world, player)
		return
	}
	if self.Location.Distance(player.Location) > 10 {
		(&Wander{}).Step(world, self)
		return
	}
	// Move towards the target, if possible
	directions := Cardinals()
	current := self.Location.Distance(player.Location)
	choices := []P{}
	alternates := []P{}
	for _, direction := range directions {
		adjacent := self.Location.Add(direction)
		if world.MoveTo(adjacent).Type != BumpEmpty {
			continue
		}
		if adjacent.Distance(player.Location) == current {
			alternates = append(alternates, direction)
			continue
		}
		if adjacent.Distance(player.Location) > current {
			continue
		}
		choices = append(choices, direction)
	}
	if len(choices) == 0 {
		if len(alternates) == 0 {
			(&Wander{}).Step(world, self)
			return
		}
		choices = alternates
	}
	choice := choices[rand.Intn(len(choices))]
	self.MoveTo(self.Location.Add(choice))
}
