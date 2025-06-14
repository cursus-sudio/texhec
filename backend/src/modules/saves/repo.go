package saves

type savableRepositories struct {
	repositories map[RepoId]SavableRepo
	sealed       bool
}

func newSavableRepositories() SavableRepositories {
	return &savableRepositories{
		repositories: map[RepoId]SavableRepo{},
		sealed:       false,
	}
}

func (savableRepositories *savableRepositories) AddRepo(id RepoId, repo SavableRepo) error {
	if savableRepositories.sealed {
		return ErrSealedSavableRepositories
	}
	savableRepositories.repositories[id] = repo
	return nil
}

func (savableRepositories *savableRepositories) GetRepositories() map[RepoId]SavableRepo {
	return savableRepositories.repositories
}

func (savableRepositories *savableRepositories) Seal() {
	savableRepositories.sealed = true
}
