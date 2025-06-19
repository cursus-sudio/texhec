package saves

import (
	"backend/services/clock"
)

type saves struct {
	Storage SavesStorage
	Repo    SavesMetaRepo
	Clock   clock.Clock
}

func newSaves(
	storage SavesStorage,
	repo SavesMetaRepo,
	clock clock.Clock,
) Saves {
	return &saves{
		Storage: storage,
		Repo:    repo,
		Clock:   clock,
	}
}

func (saves *saves) ListSaves(query ListSavesQuery) ([]SaveMeta, error) {
	return saves.Repo.ListSaves(query)
}

func (saves *saves) SavesPages(query ListSavesQuery) (int, error) {
	return saves.Repo.SavesPages(query)
}

func (saves *saves) Load(id SaveId) error {
	return saves.Storage.LoadSave(id)
}

func (saves *saves) Save(id SaveId) error {
	meta, err := saves.Repo.GetById(id)
	if err != nil {
		return err
	}
	meta.LastModified = saves.Clock.Now()

	if err := saves.Repo.Upsert(meta); err != nil {
		return err
	}
	if err := saves.Storage.Save(id); err != nil {
		return err
	}
	return nil
}

func (saves *saves) Delete(id SaveId) error {
	if err := saves.Storage.Delete(id); err != nil {
		return err
	}
	if err := saves.Repo.Delete(id); err != nil {
		return err
	}
	return nil
}

// if save already exists it overwrites it
func (saves *saves) NewSave(meta SaveMeta) error {
	if err := saves.Repo.Upsert(meta); err != nil {
		return err
	}
	if err := saves.Storage.Save(meta.Id); err != nil {
		return err
	}
	return nil
}
