package main

import (
	backendapi "backend/services/api"
	backendtcp "backend/services/api/tcp"
	"backend/services/clients"
	"backend/services/db"
	"backend/services/files"
	"backend/services/saves"
	"backend/services/scopes"
	backendscopes "backend/services/scopes"
	"core/ping"
	"core/tacticalmap"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"shared/services/api"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/null"
)

func backendDic(
	sharedPkg SharedPkg,
) ioc.Dic {
	engineDir, err := os.Getwd()
	if err != nil {
		panic(errors.Join(errors.New("current wordking direcotry"), err))
	}
	// parent of both /backend and /frontend directory
	engineDir = filepath.Dir(engineDir)
	userStorage := filepath.Join(engineDir, "user_storage")

	pkgs := []ioc.Pkg{
		sharedPkg,
		api.Package(
			func(c ioc.Dic) ioc.Dic { return c.Scope(backendscopes.Request) },
		),
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		backendtcp.Package(
			"0.0.0.0",
			"8080",
			"tcp",
		),
		backendapi.Package(),
		clients.Package(),
		db.Package(
			fmt.Sprintf("%s/db.sql", userStorage),
			null.New(fmt.Sprintf("%s/engine/backend/services/db/migrations", engineDir)),
		),
		files.Package(fmt.Sprintf("%s/files", userStorage)),
		saves.Package(),
		scopes.Package(),

		// mods
		exBackendModPkg{},
		ping.BackendPackage(),
		tacticalmap.BackendPackage(),
	}

	b := ioc.NewBuilder()
	for _, pkg := range pkgs {
		pkg.Register(b)
	}
	// backendPkg.Register(b)

	return b.Build()
}
