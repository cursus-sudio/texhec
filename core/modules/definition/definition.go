package definition

type DefinitionID uint32

type DefinitionComponent struct {
	ID DefinitionID
}

func NewDefinition(id DefinitionID) DefinitionComponent {
	return DefinitionComponent{id}
}

type DefinitionLinkComponent struct {
	DefinitionID DefinitionID
}

func NewLink(id DefinitionID) DefinitionLinkComponent {
	return DefinitionLinkComponent{id}
}
