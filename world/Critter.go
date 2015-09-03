package world

import "fmt"
import "math/rand"
import "github.com/nsf/termbox-go"

type Critter struct {
	Location      P
	Tile          Tile
	Health        Bar
	Person        Person
	Attack        Attack
	Evasion       int
	Brain         Brain
	Drops         []Item
	DistanceField DistanceField
}

func (c *Critter) At() P {
	return c.Location
}
func (self *Critter) MoveTo(loc P) {
	self.Location = loc
	self.DistanceField = DistanceField{} // Clear the distance field
}
func (self *Critter) Appearance() Tile {
	return self.Tile
}

func (self *Critter) BasicName() string {
	if self.Appearance().Name == "$player" {
		return "you"
	}
	return fmt.Sprintf("the %s", self.Appearance().Name)
}

func (self *Critter) Distance(world *Map) DistanceField {
	if !self.DistanceField.Active {
		self.DistanceField = CreateDistanceField(world, self.Location, 40)
	}
	return self.DistanceField
}

func (self *Critter) ReceiveDamage(world *Map, amount int) {
	self.Health.Value -= amount
	if self.Health.Value <= 0 {
		self.Health.Value = 0
		AddMessage("%s %s.", self.BasicName(), ConjugatePresent("die", self.Person))
		self.BecomeCorpse(world)
	}
}

func (self *Critter) AttackTarget(world *Map, target *Critter) {
	miss := rand.Intn(target.Evasion) >= rand.Intn(self.Attack.Accuracy)
	if miss {
		// you miss the bear
		// the bear misses you
		AddMessage("%s %s %s.", self.BasicName(), ConjugatePresent("miss", self.Person), target.BasicName())
		return
	}
	damage := self.Attack.Damage.Roll()
	AddMessage("%s %s %s for %d damage.", self.BasicName(), ConjugatePresent(self.Attack.Verb, self.Person), target.BasicName(), damage)
	target.ReceiveDamage(world, damage)
}

func (self *Critter) BecomeCorpse(world *Map) {
	corpse := &Corpse{
		Location: self.Location,
		Tile: Tile{
			Symbol:     '&',
			Foreground: termbox.ColorRed,
			Solid:      false,
		},
	}
	for _, item := range self.Drops {
		item.MoveTo(self.At())
		world.PlaceItem(item)
	}
	world.AddEntity(corpse)
	world.RemoveEntity(self)
}
