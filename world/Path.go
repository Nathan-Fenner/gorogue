package world

import "math/rand"

type DistanceField struct {
	Active       bool
	Target       P
	Maximum      int
	DistanceFrom map[P]int
}

func (f DistanceField) Distance(at P) int {
	if value, ok := f.DistanceFrom[at]; ok {
		return value
	}
	return f.Maximum + 1
}
func (f DistanceField) Next(at P, world *Map) P {
	bestNeighbors := []P{at}
	bestValue := f.Distance(at)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			neighbor := p(dx, dy).Add(at)
			if hit := world.MoveTo(neighbor); hit.Type != BumpEmpty {
				continue
			}
			if f.Distance(neighbor) < bestValue {
				bestNeighbors = []P{}
				bestValue = f.Distance(neighbor)
			}
			if f.Distance(neighbor) == bestValue {
				bestNeighbors = append(bestNeighbors, neighbor)
			}
		}
	}
	index := rand.Intn(len(bestNeighbors))
	return bestNeighbors[index]
}

func CreateDistanceField(world *Map, target P, maximum int) DistanceField {
	stack := []P{target}
	field := DistanceField{
		Active:       true,
		Target:       target,
		Maximum:      maximum,
		DistanceFrom: map[P]int{target: 0},
	}
	for len(stack) > 0 {
		at := stack[0]
		stack = stack[1:]
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				neighbor := at.Add(p(dx, dy))
				if tile, ok := world.Tiles[neighbor]; !ok || tile.Solid {
					continue
				}
				if _, ok := field.DistanceFrom[neighbor]; ok {
					continue
				}
				field.DistanceFrom[neighbor] = field.DistanceFrom[at] + 1
				if field.DistanceFrom[neighbor] == maximum {
					continue
				}
				stack = append(stack, neighbor)
			}
		}
	}
	return field
}
