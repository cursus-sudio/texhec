package text

// type Overflow struct {
// 	Visible bool
// }

const (
	BreakNone uint8 = iota
	BreakWord
	BreakAny
)

type BreakComponent struct {
	Break uint8
}
