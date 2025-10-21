package text

import "frontend/services/assets"

// this is required to render text
// every other component is optional and has default value
type Text struct {
	Text string
}

type TextAlign struct {
	// value between 0 and 1 where 0 means aligned to left and 1 aligned to right
	Vertical, Horizontal float32 // default is 0
}

type FontFamily struct {
	FontAsset assets.AssetID
}

type FontSize struct {
	FontSize uint
}
