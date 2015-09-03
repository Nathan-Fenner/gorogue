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

	field := player.Distance(world)

	self.MoveTo(field.Next(self.Location, world))
}
