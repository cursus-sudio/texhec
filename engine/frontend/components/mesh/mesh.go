package mesh

import "frontend/components/transform"

type Mesh struct {
	ID   string
	Size transform.Size
}

func NewMesh(id string, size transform.Size) Mesh {
	return Mesh{
		ID:   id,
		Size: size,
	}
}
