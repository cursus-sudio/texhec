package groups

type Group uint

const (
	Groupless Group = 0
)

type GroupsComponent struct {
	Mask uint32 // this can be swapped to uint64 etc (remember to swap all uint32 occurencies)
}

func EmptyGroups() GroupsComponent {
	return GroupsComponent{}
}

func DefaultGroups() GroupsComponent {
	return GroupsComponent{
		Mask: 0b1,
	}
}

func (g1 GroupsComponent) Ptr() *GroupsComponent { return &g1 }
func (g1 *GroupsComponent) Val() GroupsComponent { return *g1 }

func (g1 *GroupsComponent) Enable(g Group) *GroupsComponent {
	var mask uint32 = 0b1 << g
	g1.Mask = g1.Mask | mask
	return g1
}

func (g1 *GroupsComponent) Enabled(g Group) bool {
	var mask uint32 = 0b1 << g
	return g1.Mask&mask != 0
}

func (g1 *GroupsComponent) Disable(g Group) *GroupsComponent {
	var mask uint32 = ^(0b1 << g)
	g1.Mask = g1.Mask & mask
	return g1
}

func (g1 *GroupsComponent) GetSharedWith(g2 GroupsComponent) GroupsComponent {
	return GroupsComponent{Mask: g1.Mask & g2.Mask}
}

func (g1 *GroupsComponent) SharesAnyGroup(g2 GroupsComponent) bool {
	return g1.GetSharedWith(g2).Mask != 0
}
