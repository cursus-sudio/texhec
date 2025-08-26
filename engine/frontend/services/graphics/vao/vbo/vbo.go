package vbo

type VBO interface {
	ID() uint32
	Len() int
	Configure()
	Release()
}

type VBOSetter[Vertex any] interface {
	VBO
	SetVertices(vertices []Vertex)
}
