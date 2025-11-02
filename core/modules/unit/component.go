package unit

type UnitPos struct{ X, Y int32 }

type UnitComponent struct {
	Pos  UnitPos
	Type uint32
}
