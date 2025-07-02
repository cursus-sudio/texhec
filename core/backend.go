package main

import (
	"backend"
	backendapi "backend/services/api"
	backendtcp "backend/services/api/tcp"
	"backend/services/db"
	"backend/services/files"
	backendscopes "backend/services/scopes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"shared"
	"shared/services/api"
	"shared/services/api/netconnection"
	"shared/services/clock"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/null"
)

func backendDic(
	netconnectionPkg netconnection.Pkg,
	clockPkg clock.Pkg,
) ioc.Dic {
	engineDir, err := os.Getwd()
	if err != nil {
		panic(errors.Join(errors.New("current wordking direcotry"), err))
	}
	// parent of both /backend and /frontend directory
	engineDir = filepath.Dir(engineDir)
	userStorage := filepath.Join(engineDir, "user_storage")

	var backendPkg backend.Pkg = backend.Package(
		shared.Package(
			api.Package(
				netconnectionPkg,
				func(c ioc.Dic) ioc.Dic { return c.Scope(backendscopes.Request) },
			),
			clockPkg,
			logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		),
		backendapi.Package(
			backendtcp.Package(
				"0.0.0.0",
				"8080",
				"tcp",
			),
		),
		db.Package(
			fmt.Sprintf("%s/db.sql", userStorage),
			null.New(fmt.Sprintf("%s/engine/backend/services/db/migrations", engineDir)),
		),
		files.Package(fmt.Sprintf("%s/files", userStorage)),
		[]ioc.Pkg{
			exBackendModPkg{},
			ServerPackage(),
		},
	)

	b := ioc.NewBuilder()
	backendPkg.Register(b)

	return b.Build()
}
