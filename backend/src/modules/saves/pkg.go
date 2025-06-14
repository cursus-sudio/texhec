package saves

import (
	"backend/src/utils/clock"
	"backend/src/utils/db"
	"backend/src/utils/files"
	"backend/src/utils/uuid"

	"github.com/ogiusek/ioc"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(c ioc.Dic) {
	ioc.RegisterTransient(c, func(c ioc.Dic) Saves {
		return newSaves(
			ioc.Get[SavesStorage](c),
			ioc.Get[SavesMetaRepo](c),
			ioc.Get[clock.Clock](c),
		)
	})
	ioc.RegisterSingleton(c, func(c ioc.Dic) SaveMetaFactory {
		return newSaveMetaFactory(
			ioc.Get[clock.Clock](c),
			ioc.Get[uuid.Factory](c),
		)
	})
	ioc.RegisterScoped(c, func(c ioc.Dic) SavesMetaRepo {
		return newSavesMetaRepo(
			ioc.Get[db.Tx](c),
			ioc.Get[clock.DateFormat](c),
		)
	})
	ioc.RegisterScoped(c, func(c ioc.Dic) SavesStorage {
		return newSavesStorage(
			ioc.Get[StateCodec](c),
			ioc.Get[files.FileStorage](c),
		)
	})
	ioc.RegisterSingleton(c, func(c ioc.Dic) StateCodecRWMutex { return newStateCodecRWMutex() })
	ioc.RegisterSingleton(c, func(c ioc.Dic) StateCodec {
		savableRepositories := ioc.Get[SavableRepositories](c)
		savableRepositories.Seal()
		return newStateCodec(
			savableRepositories.GetRepositories(),
			ioc.Get[StateCodecRWMutex](c).RWMutex(),
		)
	})
	ioc.RegisterSingleton(c, func(c ioc.Dic) SavableRepositories { return newSavableRepositories() })
}
