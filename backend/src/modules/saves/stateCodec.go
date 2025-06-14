package saves

import (
	"encoding/json"
	"fmt"
	"sync"
)

type repoStateCodec struct {
	repoWMutex   sync.Locker
	repositories map[RepoId]SavableRepo
}

func newStateCodec(
	repositories map[RepoId]SavableRepo,
	repoWMutex sync.Locker,
) StateCodec {
	return &repoStateCodec{
		repoWMutex:   repoWMutex,
		repositories: repositories,
	}
}

func (repoStateCodec *repoStateCodec) Serialize() SaveData {
	repoStateCodec.repoWMutex.Lock()
	defer repoStateCodec.repoWMutex.Unlock()
	serializable := make(map[RepoId]RepoSnapshot, len(repoStateCodec.repositories))
	for repoId, repo := range repoStateCodec.repositories {
		serializable[repoId] = repo.TakeSnapshot()
	}
	bytes, _ := json.Marshal(serializable)
	return NewSaveData(bytes)
}

func (repoStateCodec *repoStateCodec) HasChanges() bool {
	for _, repo := range repoStateCodec.repositories {
		if repo.HasChanges() {
			return true
		}
	}
	return false
}

func (repoStateCodec *repoStateCodec) Load(data SaveData) error {
	repoStateCodec.repoWMutex.Lock()
	defer repoStateCodec.repoWMutex.Unlock()

	// get repositories snapshots
	snapshots := make(map[RepoId]RepoSnapshot, len(repoStateCodec.repositories))
	if err := json.Unmarshal(data, &snapshots); err != nil {
		return ErrInvalidSaveFormat
	}

	// verify snapshots
	for key, snapshot := range snapshots {
		if valid := repoStateCodec.repositories[key].IsValidSnapshot(snapshot); !valid {
			return ErrInvalidSaveFormat
		}
	}

	// apply snapshots
	for key, snapshot := range snapshots {
		err := repoStateCodec.repositories[key].LoadSnapshot(snapshot)

		// if verified snapshot is invalid then something is wrong.
		// therefor panic.
		if err != nil {
			panic(fmt.Sprintf("repository rejected snapshot which was earlier marked as valid by it\nrepo returned error %s", err.Error()))
		}
	}

	return nil
}
