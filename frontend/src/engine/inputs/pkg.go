package inputs

import "github.com/ogiusek/ioc/v2"

type Pkg struct {
}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(c ioc.Builder) {
}
