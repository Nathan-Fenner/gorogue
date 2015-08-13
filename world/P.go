package world

type P struct {
	X int
	Y int
}

func (loc P) Add(d P) P {
	return P{
		X: loc.X + d.X,
		Y: loc.Y + d.Y,
	}
}
func (loc P) Sub(other P) P {
	return P{
		X: loc.X - other.X,
		Y: loc.Y - other.Y,
	}
}

func Cardinals() []P {
	return []P{
		P{1, 0},
		P{1, 1},
		P{0, 1},
		P{-1, 1},
		P{-1, 0},
		P{-1, -1},
		P{0, -1},
		P{1, -1},
	}
}

func (loc P) Length() int {
	dx := loc.X
	if dx < 0 {
		dx = -dx
	}
	dy := loc.Y
	if dy < 0 {
		dy = -dy
	}
	if dx > dy {
		return dx
	}
	return dy
}

func (loc P) Distance(other P) int {
	return loc.Sub(other).Length()
}

func p(x int, y int) P {
	return P{
		X: x,
		Y: y,
	}
}
