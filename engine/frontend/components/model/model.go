package model

import "frontend/services/assets"

type Model struct {
	Mesh     assets.AssetID
	Material assets.AssetID
}
