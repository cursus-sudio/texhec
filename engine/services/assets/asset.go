package assets

// its just a general asset type for assets which do not use game engine at all and are just golang objects
type GoAsset interface {
	// CachableAsset
	CachedAsset
}

type goAsset struct {
	asset GoAsset
}

func NewGoAsset(asset GoAsset) GoAsset {
	return &goAsset{asset: asset}
}

func (a *goAsset) Release() {}
