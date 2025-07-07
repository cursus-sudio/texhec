package example

import (
	"backend/services/saves"
	"reflect"

	"github.com/ogiusek/ioc/v2"
)

type BackendPkg struct{}

func BackendPackage() BackendPkg {
	return BackendPkg{}
}

func (BackendPkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) *intRepo {
		return NewIntRepository(
			0,
			false,
			ioc.Get[saves.StateCodecRWMutex](c).RWMutex().RLocker(),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) IntRepo { return ioc.Get[*intRepo](c) })
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s saves.SavableRepoBuilder) saves.SavableRepoBuilder {
		repoId := reflect.TypeFor[IntRepo]().String()
		s.AddRepo(saves.RepoId(repoId), ioc.Get[*intRepo](c))
		return s
	})
}
