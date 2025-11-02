package saves

import (
	"backend/services/files"
	"backend/services/scopes"
	"shared/services/clock"
	"shared/services/codec"
	"shared/services/db"
	"shared/services/ecs"
	"shared/services/uuid"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterTransient(b, func(c ioc.Dic) Saves {
		return newSaves(
			ioc.Get[SavesStorage](c),
			ioc.Get[SavesMetaRepo](c),
			ioc.Get[clock.Clock](c),
		)
	})
	ioc.RegisterDependency[Saves, SavesStorage](b)
	ioc.RegisterDependency[Saves, SavesMetaRepo](b)
	ioc.RegisterDependency[Saves, clock.Clock](b)

	ioc.RegisterSingleton(b, func(c ioc.Dic) SaveMetaFactory {
		return newSaveMetaFactory(
			ioc.Get[clock.Clock](c),
			ioc.Get[uuid.Factory](c),
		)
	})
	ioc.RegisterDependency[SaveMetaFactory, clock.Clock](b)
	ioc.RegisterDependency[SaveMetaFactory, uuid.Factory](b)

	ioc.RegisterSingleton(b, func(c ioc.Dic) ListSavesQueryBuilder { return newQueryBuilder() })

	ioc.RegisterScoped(b, scopes.Request, func(c ioc.Dic) SavesMetaRepo {
		return newSavesMetaRepo(
			ioc.Get[db.Tx](c),
			ioc.Get[clock.DateFormat](c),
		)
	})
	ioc.RegisterDependency[SavesMetaRepo, db.Tx](b)
	ioc.RegisterDependency[SavesMetaRepo, clock.DateFormat](b)

	ioc.RegisterScoped(b, scopes.Request, func(c ioc.Dic) SavesStorage {
		return newSavesStorage(
			ioc.Get[StateCodec](c),
			ioc.Get[files.FileStorage](c),
		)
	})
	ioc.RegisterDependency[SavesStorage, StateCodec](b)
	ioc.RegisterDependency[SavesStorage, files.FileStorage](b)

	ioc.RegisterSingleton(b, func(c ioc.Dic) WorldStateCodecBuilder { return NewWorldStateCodecBuilder() })

	ioc.RegisterSingleton(b, func(c ioc.Dic) StateCodecRWMutex { return newStateCodecRWMutex() })
	ioc.RegisterTransient(b, func(c ioc.Dic) StateCodec {
		return newStateCodec(
			ioc.Get[StateCodecRWMutex](c).RWMutex(),
			ioc.Get[ecs.World](c),
			ioc.Get[codec.Codec](c),
			ioc.Get[WorldStateCodecBuilder](c).arrays,
		)
	})
	ioc.RegisterDependency[StateCodec, SavableRepositories](b)
	ioc.RegisterDependency[StateCodec, StateCodecRWMutex](b)
	ioc.RegisterDependency[StateCodec, WorldStateCodecBuilder](b)

	ioc.RegisterSingleton(b, func(c ioc.Dic) SavableRepoBuilder { return newSavableRepoBuilder() })

	ioc.RegisterSingleton(b, func(c ioc.Dic) SavableRepositories { return ioc.Get[SavableRepoBuilder](c).Build() })
	ioc.RegisterDependency[SavableRepositories, SavableRepoBuilder](b)
}
