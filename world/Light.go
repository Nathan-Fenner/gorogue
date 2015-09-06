package world

import "math"

type LightGrid struct {
	from    P
	visible map[P]bool
}

func (g LightGrid) Visible(at P) bool {
	return g.visible[at.Sub(g.from)]
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

func intFloor(x float64) int {
	return int(math.Floor(x))
}

func round(x float64) int {
	return intFloor(x + 0.5)
}

func castRayOffsets(m *Level, from P, offset P, fx float64, fy float64, ox float64, oy float64) bool {
	dx := ox + float64(offset.X) - fx
	dy := oy + float64(offset.Y) - fy
	steps := 100 // for now, a really high-precision number
	px, py := float64(from.X)+fx, float64(from.Y)+fy
	l := from
	for i := 1; i <= steps; i++ {
		ax := float64(from.X) + fx + float64(i)/float64(steps)*dx
		ay := float64(from.Y) + fy + float64(i)/float64(steps)*dy
		c := p(intFloor(ax), intFloor(ay))
		if c == from.Add(offset) {
			return true
		}
		if tile, ok := m.Tiles[c]; !ok || tile.Solid {
			if tile, ok := m.Tiles[l]; !ok || tile.Solid {
				if round(ax) != round(px) || round(ay) != round(py) {
					return false
				}
			}
		}
		// This is where I'm at.
		px, py = ax, ay
		l = c
	}
	return true
}
func castRay(m *Level, from P, offset P) bool {
	os := []float64{0.05, 0.95}
	for _, fx := range os {
		for _, fy := range os {
			for _, ox := range os {
				for _, oy := range os {
					if castRayOffsets(m, from, offset, fx, fy, ox, oy) {
						return true
					}
				}
			}
		}
	}
	return false
}

func VisibleBetween(m *Level, from P, to P) bool {
	return castRay(m, from, to.Sub(from))
}

func GetVisibility(m *Level, from P) LightGrid {
	g := LightGrid{
		from:    from,
		visible: map[P]bool{},
	}
	stack := []P{p(0, 0)}
	visited := map[P]bool{p(0, 0): true}
	size := 30
	for len(stack) > 0 {
		offset := stack[0]
		stack = stack[1:]
		if offset.X > size || offset.X < -size || offset.Y > size || offset.Y < -size {
			continue
		}
		if offset != p(0, 0) && !castRay(m, from, offset) {
			continue
		}
		g.visible[offset] = true
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				n := offset.Add(p(dx, dy))
				if !visited[n] {
					visited[n] = true
					stack = append(stack, n)
				}
			}
		}
	}
	return g
}
