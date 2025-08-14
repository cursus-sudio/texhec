package grid

import (
	"errors"
	"frontend/services/assets"
)

var (
	ErrCellsShareCoords error = errors.New("cells share coordinates")
	ErrExpectedOneGrid  error = errors.New("there should be a single grid for a cell type")
)

type Pos struct{ X, Y int }
type TextureCoord struct{ X, Y int }

type Cell interface {
	Pos() Pos
	TextureCoord() TextureCoord
}

type cell struct {
	pos          Pos
	textureCoord TextureCoord
}

func NewCell(pos Pos, textureCoord TextureCoord) Cell {
	return cell{
		pos:          pos,
		textureCoord: textureCoord,
	}
}

func (cell cell) Pos() Pos                   { return cell.pos }
func (cell cell) TextureCoord() TextureCoord { return cell.textureCoord }

type Grid[GridCell Cell] struct {
	MinPos, MaxPos Pos
	Texture        assets.AssetID
}
