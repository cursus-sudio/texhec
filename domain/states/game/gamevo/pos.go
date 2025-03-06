package gamevo

// SERVICE
type PosErrors interface {
}

type Pos struct {
	x int
	y int
}

func NewPos(x int, y int) Pos {
	return Pos{
		x: x,
		y: y,
	}
}

func (pos *Pos) X() int {
	return pos.x
}

func (pos *Pos) Y() int {
	return pos.y
}

func (pos *Pos) Move(dir Direction) (Pos, error) {
	xDif := 0
	yDif := 0
	switch dir {
	case North:
		yDif += 1
	case East:
		xDif += 1
	case West:
		xDif -= 1
	case South:
		yDif -= 1
	default:
		return Pos{}, nil
	}
	return Pos{
		x: pos.x + xDif,
		y: pos.y + yDif,
	}, nil
}
