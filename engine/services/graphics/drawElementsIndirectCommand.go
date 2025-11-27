package graphics

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

func NewDrawElementsIndirectCommand(
	indexCount uint32,
	instanceCount uint32,
	firstIndex uint32,
	firstVertex uint32,
	firstInstance uint32,
) DrawElementsIndirectCommand {
	cmd := DrawElementsIndirectCommand{
		IndexCount:    indexCount,
		InstanceCount: instanceCount,
		FirstIndex:    firstIndex,
		FirstVertex:   firstVertex,
		FirstInstance: firstInstance,
	}
	return cmd
}
