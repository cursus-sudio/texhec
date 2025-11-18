package tilerenderer

import "core/modules/tile"

type PosComponent struct{ X, Y, Z int32 }

func NewPos(pos tile.PosComponent) PosComponent { return PosComponent{pos.X, pos.Y, int32(pos.Layer)} }
