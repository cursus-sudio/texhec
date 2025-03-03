package vo

type SizeInGame struct {
	Width, Height uint
}

func NewSizeInGame(w, h uint) SizeInGame {
	return SizeInGame{
		Width:  w,
		Height: h,
	}
}
