package world

type PassiveEffect interface {
	PassiveEffect()
	Describe() string
	PassiveDealDamage(self *Critter, target *Critter, amount int) int
	PassiveReceiveDamage(self *Critter, from *Critter, amount int) int
	PassiveHitEnemy(self *Critter, target *Critter)
	PassiveMissEnemy(self *Critter, target *Critter)
	PassiveIdle(self *Critter, target *Critter)
	PassiveDefeatEnemy(self *Critter, target *Critter)
}

type EmptyEffect struct {
}

func (e EmptyEffect) PassiveEffect() {
}
func (e EmptyEffect) PassiveDealDamage(self *Critter, target *Critter, amount int) int {
	return amount
}
func (e EmptyEffect) PassiveReceiveDamage(self *Critter, target *Critter, amount int) int {
	return amount
}
func (e EmptyEffect) PassiveHitEnemy(self *Critter, target *Critter) {
}
func (e EmptyEffect) PassiveMissEnemy(self *Critter, target *Critter) {
}
func (e EmptyEffect) PassiveIdle(self *Critter, target *Critter) {
}
func (e EmptyEffect) PassiveDefeatEnemy(self *Critter, target *Critter) {
}

// Increases damage dealt by the given amount
type BonusDamage struct {
	EmptyEffect
	Amount Dice
}

func (b BonusDamage) PassiveDealDamage(self *Critter, target *Critter, amount int) int {
	return amount + b.Amount.Roll()
}

type MomentumI interface {
}

type Momentum struct {
	Count    int
	Required int
	Effect   PassiveEffect
}
