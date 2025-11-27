package ebo

import "github.com/go-gl/gl/v4.5-core/gl"

type Index uint32

type EBO interface {
	ID() uint32
	Len() int
	Configure()
	Release()
	SetIndices(indices []Index)
}

type ebo struct {
	id  uint32
	len int
}

func NewEBO() EBO {
	var id uint32
	gl.GenBuffers(1, &id)
	return &ebo{
		id:  id,
		len: 0,
	}
}

func (ebo *ebo) ID() uint32 { return ebo.id }
func (ebo *ebo) Len() int   { return ebo.len }

func (ebo *ebo) Configure() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.id)
}

func (ebo *ebo) Release() {
	gl.DeleteBuffers(1, &ebo.id)
}

func (ebo *ebo) SetIndices(indices []Index) {
	indicesLen := len(indices)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.id)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indicesLen*4, gl.Ptr(indices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	ebo.len = indicesLen
}
