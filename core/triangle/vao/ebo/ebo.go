package ebo

import "github.com/go-gl/gl/v4.5-core/gl"

type Index uint32

type EBO struct {
	ID  uint32
	Len int
}

func NewEBO() EBO {
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	return EBO{
		ID:  ebo,
		Len: 0,
	}
}

func (ebo *EBO) Release() {
	gl.DeleteBuffers(1, &ebo.ID)
}

func (ebo *EBO) SetIndices(indices []Index) {
	indiciesLen := len(indices)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.ID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indiciesLen*4, gl.Ptr(indices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	ebo.Len = indiciesLen
}
