package files

import (
	"backend/services/scopes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ogiusek/ioc/v2"
	"github.com/optimus-hft/lockset"
)

type Pkg struct {
	BaseDir string
}

func Package(
	baseDir string,
) Pkg {
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("error creating directories %s", err.Error()))
	}
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

func (pkg Pkg) Register(b ioc.Builder) {
	lockSet := lockset.New()
	ioc.RegisterScoped(b, scopes.Request, func(c ioc.Dic) FileStorage {
		return NewDiskFileStorage(
			pkg.BaseDir,
			ioc.Get[scopes.RequestEnd](c),
			lockSet,
		)
	})

}
