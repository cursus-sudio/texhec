package vo

type TurnsInGame struct {
	Turns uint
}

func NewTurnsInGame(turns uint) TurnsInGame {
	return TurnsInGame{Turns: turns}
}
