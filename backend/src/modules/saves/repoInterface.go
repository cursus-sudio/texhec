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

var (
	ErrSealedSavableRepositories error = errors.New("already sealed savable repo")
)

type SavableRepositories interface {
	// can return ErrSealedSavableRepositories
	AddRepo(RepoId, SavableRepo) error
	GetRepositories() map[RepoId]SavableRepo
	Seal()
}
