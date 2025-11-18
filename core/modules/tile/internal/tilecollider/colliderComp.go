package tilecollider

import "core/modules/tile"

type ColliderComponent struct {
	LayersBitmask uint8
}

func NewCollider() ColliderComponent {
	return ColliderComponent{0}
}

func (comp ColliderComponent) Layers() []tile.Layer {
	layers := make([]tile.Layer, 0, 8)
	for i := uint(0); i < 8; i++ {
		currentLayerBit := tile.Layer(1 << i)

		if comp.Has(currentLayerBit) {
			layers = append(layers, tile.Layer(currentLayerBit))
		}
	}
	return layers
}

func (comp ColliderComponent) Ptr() *ColliderComponent { return &comp }
func (comp *ColliderComponent) Val() ColliderComponent { return *comp }

func layerMask(layer tile.Layer) uint8 {
	return uint8(1) << layer
}

func (comp *ColliderComponent) Has(layer tile.Layer) bool {
	mask := layerMask(layer)
	return comp.LayersBitmask&mask != 0
}

func (comp *ColliderComponent) Add(layer tile.Layer) *ColliderComponent {
	mask := layerMask(layer)
	comp.LayersBitmask |= mask
	return comp
}

func (comp *ColliderComponent) Remove(layer tile.Layer) *ColliderComponent {
	mask := layerMask(layer)
	comp.LayersBitmask &= ^mask
	return comp
}
