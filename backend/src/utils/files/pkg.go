package files

import (
	"backend/src/utils/services/scopecleanup"
	"os"
	"path/filepath"

	"github.com/ogiusek/ioc"
	"github.com/optimus-hft/lockset"
)

type Pkg struct {
	BaseDir string
}

func Package(
	baseDir string,
) Pkg {
	info, err := os.Stat(baseDir)
	if err != nil {
		panic(err)
	}
	if !info.IsDir() {
		panic("path is a file, not a directory: " + filepath.Clean(baseDir))
	}
	return Pkg{
		BaseDir: baseDir,
	}
}

func (pkg Pkg) Register(c ioc.Dic) {
	lockSet := lockset.New()
	ioc.RegisterScoped(c, func(c ioc.Dic) FileStorage {
		return NewDiskFileStorage(
			pkg.BaseDir,
			ioc.Get[scopecleanup.ScopeCleanUp](c),
			lockSet,
		)
	})

}
