package world

import "math/rand"

type Dice struct {
	Count int
	Size  int
	Base  int
}

func (d Dice) Roll() int {
	value := d.Base
	for i := 0; i < d.Count; i++ {
		value += rand.Intn(d.Size) + 1
	}
	return value
}
