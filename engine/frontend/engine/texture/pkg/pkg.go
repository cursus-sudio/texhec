package texturepkg

import "github.com/ogiusek/ioc/v2"

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
}
