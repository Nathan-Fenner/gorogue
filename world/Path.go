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
func (f DistanceField) Next(at P, world *Level) P {
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

func CreateDistanceField(world *Level, target P, maximum int) DistanceField {
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
				if tile, ok := world.Tiles[neighbor]; !ok || !tile.Passable {
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

func TracePath(dist map[P]int, from P, to P, onlyCardinal bool) []P {
	if from == to {
		return []P{from}
	}
	bests := []P{}
	bestValue := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			if onlyCardinal && dx != 0 && dy != 0 {
				// Paths will only use the 4 cardinal directions
				continue
			}
			n := from.Add(p(dx, dy))
			if dist[n] == 0 {
				continue
			}
			if len(bests) == 0 || dist[n] < bestValue {
				bests = []P{}
				bestValue = dist[n]
			}
			if bestValue == dist[n] {
				bests = append(bests, n)
			}
		}
	}
	parent := bests[rand.Intn(len(bests))]
	return append(TracePath(dist, parent, to, onlyCardinal), from)
}

type PathHeap struct {
	dist     []int
	points   []*[]P
	priority map[int]*[]P
}

func NewHeap() *PathHeap {
	return &PathHeap{
		dist:     []int{},
		points:   []*[]P{},
		priority: map[int]*[]P{},
	}
}
func (p *PathHeap) Empty() bool {
	return len(p.dist) == 0
}
func (p *PathHeap) swap(i, j int) {
	p.dist[i], p.dist[j] = p.dist[j], p.dist[i]
	p.points[i], p.points[j] = p.points[j], p.points[i]
}
func (p *PathHeap) retract() {
	p.dist = p.dist[:len(p.dist)-1]
	p.points = p.points[:len(p.points)-1]
}
func (p *PathHeap) bubbleDown(i int) {
	if 2*i+1 >= len(p.dist) {
		return // cannot bubble down further
	}
	if 2*i+2 == len(p.dist) {
		// only have one parent to consider
		if p.dist[i] > p.dist[2*i+1] {
			p.swap(i, 2*i+1)
			p.bubbleDown(2*i + 1)
		}
		return
	}
	next := 2*i + 1
	if p.dist[2*i+1] > p.dist[2*i+2] {
		next = 2*i + 2
	}
	if p.dist[i] > p.dist[next] {
		p.swap(i, next)
		p.bubbleDown(next)
	}
}
func (p *PathHeap) bubbleUp(i int) {
	if i == 0 {
		return
	}
	if p.dist[(i-1)/2] > p.dist[i] {
		p.swap((i-1)/2, i)
		p.bubbleUp((i - 1) / 2)
	}
}
func (p *PathHeap) Best() P {
	bestList := p.points[0]
	i := rand.Intn(len(*bestList))
	best := (*bestList)[i]
	(*bestList)[i], (*bestList)[len(*bestList)-1] = (*bestList)[len(*bestList)-1], (*bestList)[i]
	*bestList = (*bestList)[:len(*bestList)-1]
	if len(*bestList) == 0 {
		p.swap(0, len(p.dist)-1)
		p.retract()
		p.bubbleDown(0)
	}
	return best
}
func (p *PathHeap) Insert(at P, cost int) {
	if p.priority[cost] == nil {
		list := []P{}
		p.priority[cost] = &list
	}
	*p.priority[cost] = append(*p.priority[cost], at)
	if len(*p.priority[cost]) == 1 {
		p.dist = append(p.dist, cost)
		p.points = append(p.points, p.priority[cost])
		p.bubbleUp(len(p.dist) - 1)
	}
}

func FindPath(from P, to P, costFunc func(P) int, onlyCardinal bool) []P {
	dist := map[P]int{from: 1}

	heap := NewHeap()
	heap.Insert(from, 1)

	inHeap := map[P]bool{from: true}

	for {
		s := heap.Best()
		if s == to {
			return TracePath(dist, to, from, onlyCardinal)
		}
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if onlyCardinal && dx != 0 && dy != 0 {
					// Only use the 4 cardinal directions.
					continue
				}
				n := s.Add(p(dx, dy))
				cost := costFunc(n) + dist[s]
				if dist[n] == 0 || dist[n] > cost {
					dist[n] = cost
					if !inHeap[n] {
						heap.Insert(n, costFunc(n))
						inHeap[n] = true
					}
				}
			}
		}
	}
}
