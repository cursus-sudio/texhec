package game

import (
	"domain/vo"

	"github.com/ogiusek/ioc"
)

// SERVICE
type DefaultGameOptions struct {
	StartBudget vo.Budget
	MapSize     vo.SizeInGame
}

type GameOptions struct {
	StartBudget vo.Budget
	MapSize     vo.SizeInGame
	Seed        vo.Seed
	ManualSeed  bool
}

func NewPrepareOptions(c ioc.Dic) GameOptions {
	defaults := ioc.Get[DefaultGameOptions](c)
	return GameOptions{
		StartBudget: defaults.StartBudget,
		MapSize:     defaults.MapSize,
		Seed:        vo.NewSeed(c),
		ManualSeed:  false,
	}
}

func (options *GameOptions) ChangeSeed(seed vo.Seed) {
	options.Seed = seed
	options.ManualSeed = true
}

func (options *GameOptions) RandomizeSeed(c ioc.Dic) {
	options.Seed = vo.NewSeed(c)
	options.ManualSeed = false
}
