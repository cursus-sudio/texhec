package saves

import "shared/utils/httperrors"

type SavableRepoBuilder interface {
	// can return Err409
	AddRepo(RepoId, SavableRepo) error
	Build() SavableRepositories
}

type savableRepoBuilder struct {
	repositories map[RepoId]SavableRepo
}

func newSavableRepoBuilder() SavableRepoBuilder {
	return &savableRepoBuilder{
		repositories: map[RepoId]SavableRepo{},
	}
}

func (savableRepoBuilder *savableRepoBuilder) AddRepo(id RepoId, repo SavableRepo) error {
	_, ok := savableRepoBuilder.repositories[id]
	if ok {
		return httperrors.Err409
	}
	savableRepoBuilder.repositories[id] = repo
	return nil
}

func (savableRepoBuilder *savableRepoBuilder) Build() SavableRepositories {
	return &savableRepositories{
		repositories: savableRepoBuilder.repositories,
	}
}
