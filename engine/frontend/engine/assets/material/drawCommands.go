package material

// DrawElementsIndirectCommand must match GL spec layout
// typedef struct {
//   GLuint count;
//   GLuint primCount;
//   GLuint firstIndex;
//   GLuint baseVertex;
//   GLuint baseInstance;
// } DrawElementsIndirectCommand;

type DrawElementsIndirectCommand struct {
	IndexCount    uint32 // Corresponds to GLuint count
	InstanceCount uint32 // Corresponds to GLuint primCount
	FirstIndex    uint32 // Corresponds to GLuint firstIndex
	FirstVertex   uint32 // Corresponds to GLuint baseVertex
	FirstInstance uint32 // Corresponds to GLuint baseInstance
}

type MeshRange struct {
	firstIndex  uint32
	indexCount  uint32
	firstVertex uint32
}

func NewDrawElementsIndirectCommand(
	meshRange MeshRange,
	instanceCount uint32,
	firstInstance uint32,
) DrawElementsIndirectCommand {
	cmd := DrawElementsIndirectCommand{
		IndexCount:    meshRange.indexCount,
		InstanceCount: instanceCount,
		FirstIndex:    meshRange.firstIndex,
		FirstVertex:   meshRange.firstVertex,
		FirstInstance: firstInstance,
	}
	return cmd
}
