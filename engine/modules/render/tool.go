package render

import "engine/services/ecs"

type Tool interface {
	Transaction() Transaction
	Error() error
}

type Transaction interface {
	GetObject(ecs.EntityID) Object
	Transactions() []ecs.AnyComponentsArrayTransaction
	Flush() error
}

type Object interface {
	Color() ecs.EntityComponent[ColorComponent]
	Mesh() ecs.EntityComponent[MeshComponent]
	Texture() ecs.EntityComponent[TextureComponent]
}
