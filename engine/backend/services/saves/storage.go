package saves

import (
	"backend/services/files"
)

type savesStorage struct {
	Codec       StateCodec
	FileStorage files.FileStorage
}

func newSavesStorage(
	codec StateCodec,
	fileStorage files.FileStorage,
) SavesStorage {
	return &savesStorage{
		Codec:       codec,
		FileStorage: fileStorage,
	}
}

func (storage *savesStorage) LoadSave(id SaveId) error {
	bytes, err := storage.FileStorage.Read(files.NewPath(id.String()))
	if err != nil {
		return err
	}
	if err := storage.Codec.Load(bytes); err != nil {
		return err
	}
	return nil
}

func (storage *savesStorage) Save(id SaveId) error {
	data, err := storage.Codec.Serialize()
	if err != nil {
		return err
	}
	if err := storage.FileStorage.OverWrite(files.NewPath(id.String()), data); err != nil {
		return err
	}
	return nil
}

func (storage *savesStorage) Delete(id SaveId) error {
	err := storage.FileStorage.Delete(files.NewPath(id.String()))
	return err
}
