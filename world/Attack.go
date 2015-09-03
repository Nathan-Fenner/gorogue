package world

type Attack struct {
	Name     string
	Verb     string
	Damage   Dice
	Accuracy int
	Effects  []PassiveEffect
}
