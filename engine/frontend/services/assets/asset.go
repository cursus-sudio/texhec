package assets

type AssetID string

type Assets interface {
	Get(AssetID) (any, error)
	Release(...AssetID)
}

type assets struct {
	assetStorage AssetsStorage
	cachedAssets CachedAssets
}

func NewAssets(
	storage AssetsStorage,
	cached CachedAssets,
) Assets {
	return &assets{
		assetStorage: storage,
		cachedAssets: cached,
	}
}

func (a *assets) Get(id AssetID) (any, error) {
	{
		cached, err := a.cachedAssets.Get(id)
		if err == nil {
			return cached, nil
		}
	}
	stored, err := a.assetStorage.Get(id)
	if err != nil {
		return nil, err
	}
	cached, err := stored.Cache()
	if err != nil {
		return nil, err
	}
	if err := a.cachedAssets.Set(id, cached); err != nil {
		return nil, err
	}
	return cached, nil
}

func (a *assets) Release(ids ...AssetID) {
	for _, id := range ids {
		a.cachedAssets.Delete(id)
	}
}
