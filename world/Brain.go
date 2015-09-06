package world

import "math/rand"

type Brain interface {
	Step(*Level, *Critter)
}

type Wander struct {
}

func (wander *Wander) Step(world *Level, self *Critter) {
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
	TargetPosition P
}

func (hunter *Hunter) Step(world *Level, self *Critter) {
	player := world.Player()
	if self.Location.Distance(player.Location) == 1 {
		self.AttackTarget(world, player)
		return
	}
	playerSighted := VisibleBetween(world, self.Location, player.Location)
	playerField := player.Distance(world)
	if playerSighted && playerField.Distance(self.Location) <= 15 {
		hunter.TargetPosition = player.Location
	}

	if self.Location == hunter.TargetPosition {
		(&Wander{}).Step(world, self)
		hunter.TargetPosition = self.Location // To ensure that wandering continues
		return
	}

	field := CreateDistanceField(world, hunter.TargetPosition, 15)
	self.MoveTo(field.Next(self.Location, world))
}
