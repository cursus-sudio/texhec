package groups

type Group int

const (
	Groupless Group = 0
)

type Groups struct {
	Mask uint32 // this can be swapped to uint64 etc (remember to swap all uint32 occurencies)
}

func EmptyGroups() Groups {
	return Groups{}
}

func DefaultGroups() Groups {
	return Groups{
		Mask: 0b1,
	}
}

func (g1 Groups) Ptr() *Groups { return &g1 }
func (g1 *Groups) Val() Groups { return *g1 }

func (g1 *Groups) Enable(g Group) *Groups {
	var mask uint32 = 0b1 << g
	g1.Mask = g1.Mask | mask
	return g1
}

func (g1 *Groups) Enabled(g Group) bool {
	var mask uint32 = 0b1 << g
	return g1.Mask&mask != 0
}

func (g1 *Groups) Disable(g Group) *Groups {
	var mask uint32 = ^(0b1 << g)
	g1.Mask = g1.Mask & mask
	return g1
}

func (g1 *Groups) GetSharedWith(g2 Groups) Groups {
	return Groups{Mask: g1.Mask & g2.Mask}
}

func (g1 *Groups) SharesAnyGroup(g2 Groups) bool {
	return g1.GetSharedWith(g2).Mask != 0
}
