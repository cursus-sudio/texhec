package text

import (
	"frontend/services/assets"

	"github.com/go-gl/mathgl/mgl32"
)

// this is required to render text
// every other component is optional and has default value
type TextComponent struct {
	Text string
}

type TextAlignComponent struct {
	// value between 0 and 1 where 0 means aligned to left and 1 aligned to right
	Vertical, Horizontal float32 // default is 0
}

type TextColorComponent struct {
	Color mgl32.Vec4
}

type FontFamilyComponent struct {
	FontFamily assets.AssetID
}

type FontSizeComponent struct {
	FontSize uint
}
