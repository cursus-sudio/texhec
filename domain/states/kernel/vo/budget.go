package vo

import "github.com/ogiusek/ioc"

// SERVICE
type BudgetErrors interface {
	MissingMoney(m1, m2 uint) error
	MissingManpower(m1, m2 uint) error
	MissingTexHec(t1, t2 uint) error
}

// budget

type Budget struct {
	money    uint
	manpower uint
	texHec   uint
}

func NewBudget(money uint, manpower uint, texHec uint) Budget {
	return Budget{
		money:    money,
		manpower: manpower,
		texHec:   texHec,
	}
}

func EmptyBudget() Budget {
	return Budget{
		money:    0,
		manpower: 0,
		texHec:   0,
	}
}

// return (sum Budget, overflow Budget)
func (max Budget) MaxSum(b1 Budget, b2 Budget) (Budget, Budget) {
	sum := EmptyBudget()
	overflow := EmptyBudget()
	moneySum := b1.money + b2.money
	if moneySum > max.money {
		sum.money = max.money
		overflow.money = moneySum - max.money
	}
	manpowerSum := b1.manpower + b2.manpower
	if manpowerSum > max.manpower {
		sum.manpower = max.manpower
		overflow.manpower = manpowerSum - max.manpower
	}
	texHecSum := b1.texHec + b2.texHec
	if texHecSum > max.texHec {
		sum.texHec = max.texHec
		overflow.texHec = texHecSum - max.texHec
	}
	return sum, overflow
}

func (b1 *Budget) Add(c ioc.Dic, b2 Budget) Budget {
	return Budget{
		money:    b1.money + b2.money,
		manpower: b1.manpower + b2.manpower,
		texHec:   b1.texHec + b2.texHec,
	}
}

func (b1 *Budget) Subtract(c ioc.Dic, b2 Budget) (Budget, []error) {
	budegetErrors := ioc.Get[BudgetErrors](c)
	var errs []error
	if b1.money < b2.money {
		errs = append(errs, budegetErrors.MissingMoney(b1.money, b2.money))
	}
	if b1.manpower < b2.manpower {
		errs = append(errs, budegetErrors.MissingManpower(b1.manpower, b2.manpower))
	}
	if b1.texHec < b2.texHec {
		errs = append(errs, budegetErrors.MissingTexHec(b1.texHec, b2.texHec))
	}

	if len(errs) != 0 {
		return Budget{}, errs
	}

	return Budget{
		money:    b1.money - b2.money,
		manpower: b1.manpower - b2.manpower,
		texHec:   b1.texHec - b2.texHec,
	}, nil
}
