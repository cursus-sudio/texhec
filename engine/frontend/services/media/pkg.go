package media

import (
	"frontend/services/media/inputs"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	// TODO
	if true {
		return
	}
	inputs.Package().Register(b)
}
