package vo

type TurnsInBattle struct {
	Turns uint
}

func NewTurnsInBattle(turns uint) TurnsInBattle {
	return TurnsInBattle{Turns: turns}
}
