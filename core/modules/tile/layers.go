package tile

type Layer uint8

const (
	GroundLayer Layer = iota
	BuildingLayer
	UnitLayer
)
