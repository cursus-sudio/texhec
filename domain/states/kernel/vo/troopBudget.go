package vo

import (
	"github.com/ogiusek/ioc"
)

// SERVICE
type TroopBudgetErrors interface {
	PayExceededBudget(budget TroopBudget, pay TroopBudget) error
}

// battle budget

type TroopBudget struct {
	budget uint
}

func (troop *TroopBudget) Valid(c ioc.Dic) []error {
	// TODO
	return nil
}

func (b *TroopBudget) Multiply(m *TroopBudgetMultiplier) TroopBudget {
	var res TroopBudget = *b
	res.budget = uint(m.budget * float32(b.budget))
	return res
}

func MinTroopBudget(t1 *TroopBudget, t2 *TroopBudget) TroopBudget {
	return TroopBudget{
		budget: min(t1.budget, t2.budget),
	}
}

func (budget *TroopBudget) Pay(c ioc.Dic, pay *TroopBudget) (TroopBudget, error) {
	if budget.budget < pay.budget {
		errors := ioc.Get[TroopBudgetErrors](c)
		return TroopBudget{}, errors.PayExceededBudget(*budget, *pay)
	}
	return TroopBudget{
		budget: budget.budget - pay.budget,
	}, nil
}

// battle budget multiplier

type TroopBudgetMultiplier struct {
	budget float32
}

func (m *TroopBudgetMultiplier) Multiply(b *TroopBudget) TroopBudget {
	return b.Multiply(m)
}
