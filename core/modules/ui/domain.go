package ui

// changes ui screen using activate below.
// attached elements are dependent from data in domain specific events
type SettingsEvent struct{}
type SelectedTileEvent struct{ X, Y uint32 }
