package render

// import "frontend/services/graphics/vao/ebo"
//
// type MeshRange struct {
// 	FirstIndex  uint32
// 	IndexCount  uint32
// 	FirstVertex uint32
// }
//
// type PackedMesh[Vertex any] struct {
// 	Vertices []Vertex
// 	Indices  []ebo.Index
// 	Ranges   []MeshRange
// }
//
// func Pack[Vertex any](meshes ...MeshAsset[Vertex]) PackedMesh[Vertex] {
// 	p := PackedMesh[Vertex]{}
// 	for _, m := range meshes {
// 		var firstVertex = uint32(len(p.Vertices))
// 		var firstIndex = uint32(len(p.Indices))
//
// 		p.Ranges = append(p.Ranges, MeshRange{
// 			FirstIndex:  firstIndex,
// 			IndexCount:  uint32(len(m.Indices())),
// 			FirstVertex: firstVertex,
// 		})
//
// 		p.Vertices = append(p.Vertices, m.Vertices()...)
// 		for _, i := range m.Indices() {
// 			p.Indices = append(p.Indices, i)
// 		}
// 	}
//
// 	return p
// }
