package vo

type SizeInBattle struct {
	Width, Height uint
}

func NewSizeInBattle(w, h uint) SizeInBattle {
	return SizeInBattle{
		Width:  w,
		Height: h,
	}
}
