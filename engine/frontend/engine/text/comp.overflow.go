package text

// type Overflow struct {
// 	Visible bool
// }

const (
	BreakNone uint8 = iota
	BreakWord
	BreakAny
)

type Break struct {
	Break uint8
}
