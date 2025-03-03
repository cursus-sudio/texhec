package vo

import (
	"github.com/ogiusek/ioc"
)

// battle budget

type TroopBudget struct { // example troop for reference
	//                   // budget
	UniversalBudget uint // 100
	MoveBudget      uint // 5
	ShootBudget     uint // 0
	// troop skills:
	// move square cost: 5
	// shoot: 25
}

func (troop *TroopBudget) Valid(c ioc.Dic) []error {
	// TODO
	return nil
}

func (b TroopBudget) Multiply(m TroopBudgetMultiplier) TroopBudget {
	b.UniversalBudget = uint(m.UniversalBudget * float32(b.UniversalBudget))
	b.MoveBudget = uint(m.MoveBudget * float32(b.MoveBudget))
	b.ShootBudget = uint(m.ShootBudget * float32(b.ShootBudget))
	return b
}

func MinTroopBudget(t1 TroopBudget, t2 TroopBudget) TroopBudget {
	return TroopBudget{
		UniversalBudget: min(t1.UniversalBudget, t2.UniversalBudget),
		MoveBudget:      min(t1.MoveBudget, t2.MoveBudget),
		ShootBudget:     min(t1.ShootBudget, t2.ShootBudget),
	}
}

// battle budget multiplier

type TroopBudgetMultiplier struct {
	UniversalBudget float32
	MoveBudget      float32
	ShootBudget     float32
}

func (m TroopBudgetMultiplier) Multiply(b TroopBudget) TroopBudget {
	return b.Multiply(m)
}
