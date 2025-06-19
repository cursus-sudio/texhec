package saves

import "errors"

var (
	ErrInvalidRepoSnapshot error = errors.New("invalid repo snapshot")
)

type RepoId string
type RepoSnapshot []byte

func (spanshot *RepoSnapshot) Bytes() []byte {
	return []byte(*spanshot)
}

func NewRepoSnapshot(bytes []byte) RepoSnapshot {
	return RepoSnapshot(bytes)
}

type SavableRepo interface {
	IsValidSnapshot(RepoSnapshot) bool

	TakeSnapshot() RepoSnapshot

	// replaces current save with loaded
	// can return:
	// - ErrInvalidSaveFormat
	LoadSnapshot(RepoSnapshot) error

	// returns was modified since last snapshot loading or saving
	HasChanges() bool
}

type SavableRepositories interface {
	GetRepositories() map[RepoId]SavableRepo
}

// imp

type savableRepositories struct {
	repositories map[RepoId]SavableRepo
}

func newSavableRepositories() SavableRepositories {
	return &savableRepositories{
		repositories: map[RepoId]SavableRepo{},
	}
}

func (savableRepositories *savableRepositories) GetRepositories() map[RepoId]SavableRepo {
	return savableRepositories.repositories
}
