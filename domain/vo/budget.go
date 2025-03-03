package vo

// budget

type Budget struct {
	Money    uint
	Manpower uint
	TexHec   uint
}

func (b Budget) Multiply(mutiplayer BudgetMultiplier) Budget {
	b.Money = uint(mutiplayer.Money * float32(b.Money))
	b.Manpower = uint(mutiplayer.Manpower * float32(b.Manpower))
	b.TexHec = uint(mutiplayer.TexHec * float32(b.TexHec))
	return b
}

// budget multiplier

type BudgetMultiplier struct {
	Money    float32
	Manpower float32
	TexHec   float32
}

func (m BudgetMultiplier) Multiply(budget Budget) Budget {
	return budget.Multiply(m)
}
