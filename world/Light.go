package world

type LightGrid struct {
	from    P
	visible map[P]bool
	active  map[P]bool
	lower   map[P]P
	upper   map[P]P
}

func (g LightGrid) Visible(at P) bool {
	return g.visible[at.Sub(g.from)]
}

func naturalTarget(from int, to int) float64 {
	if from == to {
		return 0.5
	}
	if from < to {
		return 1
	}
	return 0
}

func sign(i int) int {
	if i < 0 {
		return -1
	}
	if i > 0 {
		return 1
	}
	return 0
}

func spiral(n int) []P {
	r := []P{p(0, 0)}
	for i := 1; i <= n; i++ {
		for x := -i + 1; x < i; x++ {
			r = append(r, p(i, x))
			r = append(r, p(-i, x))
			r = append(r, p(x, i))
			r = append(r, p(x, -i))
		}
		r = append(r, p(i, i))
		r = append(r, p(-i, i))
		r = append(r, p(-i, -i))
		r = append(r, p(i, -i))
	}
	return r
}
func primary(at P) P {
	x, sx := at.X, 1
	y, sy := at.Y, 1
	if x < 0 {
		sx = -1
	}
	if y < 0 {
		sy = -1
	}
	if x == 0 {
		return p(0, sy)
	}
	if y == 0 {
		return p(sx, 0)
	}
	if x*sx == y*sy {
		return p(0, 0) // it's a pure angle
	}
	if x*sx > y*sy {
		return p(sx, 0)
	}
	return p(0, sy)
}
func secondary(at P) P {
	if at.X == 0 || at.Y == 0 {
		return p(0, 0) // it's a pure cardinal
	}
	sx := 1
	if at.X == 0 {
		sx = 0
	}
	if at.X < 0 {
		sx = -1
	}
	sy := 1
	if at.Y == 0 {
		sy = 0
	}
	if at.Y < 0 {
		sy = -1
	}
	return p(sx, sy)
}

func inOrder(low P, high P) bool {
	// low.Y / low.X <= high.Y / high.X
	// or equivalently,
	// low.Y * high.X <= high.Y * low.X
	// (provided all quantities are positive)
	return low.Y*high.X <= high.Y*low.X
}

func lower(a P, b P) P {
	if inOrder(a, b) {
		return a
	}
	return b
}
func higher(a P, b P) P {
	if inOrder(a, b) {
		return b
	}
	return a
}

func inBounds(low P, value P, high P) bool {
	return inOrder(low, value) && inOrder(value, high)
}

func angleCoord(at P) P {
	if at.X < 0 {
		at.X *= -1
	}
	if at.Y < 0 {
		at.Y *= -1
	}
	if at.X == at.Y {
		return p(0, at.X)
	}
	if at.X < at.Y {
		return p(at.Y-at.X, at.X)
	} else {
		return p(at.X-at.Y, at.Y)
	}
}

func reorientPrinciple(point P) P {
	if point.X < 0 {
		point.X *= -1
	}
	if point.Y < 0 {
		point.Y *= -1
	}
	if point.X < point.Y {
		point.X, point.Y = point.Y, point.X
	}
	return point
}

func GetVisibility(m *Map, from P) LightGrid {
	g := LightGrid{
		from:    from,
		visible: map[P]bool{},
		active:  map[P]bool{},
		lower:   map[P]P{},
		upper:   map[P]P{},
	}
	points := spiral(30)
	zero := p(0, 0)
	for i, point := range points {
		if i == 0 {
			// The central point is always active and visible
			g.visible[zero] = true
			g.active[zero] = true
			g.lower[zero] = p(2, 0)
			g.upper[zero] = p(2, 2)
			continue
		}
		// Step backwards in the primary direction:

		forgiveness := 3

		if parent := point.Sub(primary(point)); g.active[parent] {
			oriented := reorientPrinciple(point)
			// The point is active: check whether I lie in its ratios
			if inBounds(g.lower[parent], oriented, g.upper[parent]) {
				// Update visible region
				g.lower[point] = higher(g.lower[parent], oriented.Add(p(forgiveness, 0)))
				g.upper[point] = lower(g.upper[parent], oriented.Add(p(forgiveness, forgiveness)))
				g.visible[point] = true
			}
		}

		if parent := point.Sub(secondary(point)); g.active[parent] {
			oriented := reorientPrinciple(point)
			// The point is active: check whether I lie in its ratios
			if inBounds(g.lower[parent], oriented, g.upper[parent]) {
				// Update visible region, with caveat if already visible from primary:
				// We expand the sight instead of restricting.
				lowerBound := higher(g.lower[parent], oriented.Add(p(forgiveness, 0)))
				upperBound := lower(g.upper[parent], oriented.Add(p(forgiveness, forgiveness)))
				if g.visible[point] {
					g.lower[point] = lower(g.lower[point], lowerBound)
					g.upper[point] = lower(g.upper[point], upperBound)
				} else {
					g.lower[point] = lowerBound
					g.upper[point] = upperBound
				}
				g.visible[point] = true
			}
		}

		if g.visible[point] && !m.Tiles[point.Add(from)].Solid {
			g.active[point] = true
		}
	}
	return g
}
