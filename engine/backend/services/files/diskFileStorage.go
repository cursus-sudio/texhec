package files

import (
	"backend/services/scopes"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"shared/utils/httperrors"
	"time"

	"github.com/optimus-hft/lockset"
)

// os.PathError
// os.ErrPermission
// os.ErrNotExist
// os.ErrInvalid

type diskFile struct {
	Modified bool
	// nil stands for removed file
	// other arrays mean just content
	Content []byte
}

// diskFileStorage implements FileStorage for local disk operations
// since this is scoped it does not need to be thread safe
type diskFileStorage struct {
	baseDir string
	lockset *lockset.Set // this is thread safe because this supposed to be shared across scopes

	openedFiles            map[Path]*diskFile
	handleFailedFileChange func(file *diskFile, fullPath string)
}

func NewDiskFileStorage(
	baseDir string,
	requestEnd scopes.RequestEnd,
	lockset *lockset.Set,
) FileStorage {
	fileStorage := &diskFileStorage{
		baseDir:     baseDir,
		openedFiles: make(map[Path]*diskFile),
		lockset:     lockset,
		handleFailedFileChange: func(file *diskFile, fullPath string) {
			go func() {
				for {
					var err error
					if file.Content == nil {
						err = os.Remove(fullPath)
					} else {
						err = os.WriteFile(fullPath, file.Content, 0644)
					}
					if err == nil {
						break
					}
					log.Printf("failed to save file %s", err)
					time.Sleep(time.Second * 10)
				}
			}()
		},
	}

	requestEnd.AddCleanListener(fileStorage.CleanUp)

	return fileStorage
}

func (fileStorage *diskFileStorage) CleanUp(args scopes.RequestEndArgs) {
	for path, openedFile := range fileStorage.openedFiles {
		defer fileStorage.lockset.Unlock(path.String())
		if !openedFile.Modified {
			continue
		}
		fullPath := filepath.Join(fileStorage.baseDir, path.String())
		if openedFile.Content == nil {
			err := os.Remove(fullPath)
			if errors.Is(err, fs.ErrNotExist) {
			} else if err != nil {
				fileStorage.handleFailedFileChange(openedFile, fullPath)
			}
		} else {
			err := os.WriteFile(fullPath, openedFile.Content, 0644)
			if err != nil {
				fileStorage.handleFailedFileChange(openedFile, fullPath)
			}
		}
	}
}

// only errors which can be returned are 500 errors
func (fileStorage *diskFileStorage) LoadFileOrDefault(path Path, defaultFileArg diskFile) (*diskFile, error) {
	defaultFile := defaultFileArg
	fileStorage.lockset.Lock(path.String())
	fullPath := filepath.Join(fileStorage.baseDir, path.String())

	var file *diskFile
	content, err := os.ReadFile(fullPath)
	if err == nil {
		file = &diskFile{
			Modified: false,
			Content:  content,
		}
	} else if errors.Is(err, os.ErrNotExist) {
		file = &defaultFile
	} else if os.IsPermission(err) {
		return nil, httperrors.Err404
	} else {
		return nil, httperrors.Err503
	}

	fileStorage.openedFiles[path] = &defaultFile
	return file, nil
}

// only errors which can be returned are 500 errors
func (fileStorage *diskFileStorage) LoadFile(path Path, fileArg diskFile) (*diskFile, error) {
	file := fileArg
	_, ok := fileStorage.openedFiles[path]
	if !ok {
		fileStorage.lockset.Lock(path.String())
	}
	fileStorage.openedFiles[path] = &file
	return &file, nil
}

func (fileStorage *diskFileStorage) EnsureExists(path Path) error {
	file, ok := fileStorage.openedFiles[path]
	if !ok {
		_, err := fileStorage.LoadFileOrDefault(path, diskFile{
			Modified: true,
			Content:  []byte{},
		})

		if err != nil {
			return err
		}
	} else if file.Content == nil {
		file.Content = []byte{}
		file.Modified = true
	}
	return nil
}

func (fileStorage *diskFileStorage) Exists(path Path) (bool, error) {
	file, ok := fileStorage.openedFiles[path]
	if !ok {
		var err error
		file, err = fileStorage.LoadFileOrDefault(path, diskFile{
			Modified: false,
			Content:  nil,
		})
		if err != nil {
			return false, err
		}
	}
	return file.Content != nil, nil
}

func (fileStorage *diskFileStorage) Read(path Path) ([]byte, error) {
	file, ok := fileStorage.openedFiles[path]
	if !ok {
		var err error
		file, err = fileStorage.LoadFileOrDefault(path, diskFile{
			Modified: false,
			Content:  nil,
		})
		if err != nil {
			return nil, err
		}
	}
	if file.Content == nil {
		return nil, httperrors.Err404
	}
	return file.Content, nil
}

func (fileStorage *diskFileStorage) OverWrite(path Path, contents []byte) error {
	_, err := fileStorage.LoadFile(path, diskFile{
		Modified: true,
		Content:  contents,
	})
	return err
}

func (fileStorage *diskFileStorage) Write(path Path, contents []byte) error {
	file, ok := fileStorage.openedFiles[path]
	if !ok {
		var err error
		file, err = fileStorage.LoadFileOrDefault(path, diskFile{
			Modified: false,
			Content:  nil,
		})
		if err != nil {
			return err
		}
	}
	file.Modified = true
	file.Content = contents
	return nil
}

func (fileStorage *diskFileStorage) Delete(path Path) error {
	_, err := fileStorage.LoadFile(path, diskFile{
		Modified: true,
		Content:  nil,
	})
	return err
}
