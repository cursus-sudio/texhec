package inputs

import "github.com/ogiusek/ioc"

type Pkg struct {
}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(c ioc.Dic) {
}
